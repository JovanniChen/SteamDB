// Dao包 - 数据访问对象，负责Steam平台的HTTP请求和数据交互
// 提供HTTP客户端、Cookie管理、请求重试等核心功能
package Dao

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	u "net/url"
	"sync"
	"time"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
)

type globalConfig struct {
	daos sync.Map
}

var globalDaos *globalConfig // 全局DAO对象池

// Dao 数据访问对象结构体
// 封装了HTTP客户端和用户凭据，提供Steam API交互功能
type Dao struct {
	httpCli         *http.Client // HTTP客户端，用于发送网络请求
	credentials     *Credentials // 用户凭据信息，包含登录状态和认证信息
	requestCallback func()       // HTTP请求成功后的回调函数，用于外部监控请求
	proxy           string
}

func reInit() {
	globalDaos = &globalConfig{
		daos: sync.Map{},
	}
}

func init() {
	reInit()
}

// Request 创建包含认证信息的HTTP请求
// 自动添加登录必需的Cookie信息，包括steamLoginSecure、sessionid等
// 参数：
//
//	method - HTTP方法（GET、POST等）
//	url - 请求URL
//	body - 请求体数据
//
// 返回值：配置好认证信息的HTTP请求对象
func (d *Dao) Request(method, url string, body io.Reader) (*http.Request, error) {
	// 创建基础HTTP请求
	req, err := d.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 解析URL以获取主机名，用于设置Cookie的域
	ur, _ := u.Parse(url)

	// 获取该域名对应的登录Cookie信息
	ck := d.GetCookiesString(url)
	if ck == nil {
		return nil, Errors.Error("Cookie not exist")
	}

	// 添加Steam登录安全Cookie，用于身份验证
	req.AddCookie(&http.Cookie{
		Name:   "steamLoginSecure",
		Value:  ck.SteamLoginSecure,
		Domain: ur.Host,
		Path:   "/",
	})

	// 添加会话ID Cookie，维持会话状态
	req.AddCookie(&http.Cookie{
		Name:   "sessionid",
		Value:  ck.SessionId,
		Domain: ur.Host,
		Path:   "/",
	})

	// 如果设置了语言偏好，添加语言Cookie
	if d.credentials.Language != "" {
		req.AddCookie(&http.Cookie{
			Name:   "Steam_Language",
			Value:  d.credentials.Language,
			Domain: ur.Host,
			Path:   "/",
		})
	}
	return req, nil
}

// NewRequest 创建基础HTTP请求
// 设置通用的请求头信息，包括User-Agent和Content-Type
// 参数：
//
//	method - HTTP请求方法
//	url - 目标URL
//	body - 请求体内容
//
// 返回值：配置好基础头信息的HTTP请求对象
func (d *Dao) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 设置浏览器用户代理，模拟真实浏览器请求以避免被反爬虫检测
	req.Header.Set("User-Agent", Constants.UserAgent)
	// req.Header.Set("accept-language", "zh-CN,zh;q=0.9")

	// 对于POST请求，设置表单数据的Content-Type
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
}

// Do 执行HTTP请求
// 使用内部HTTP客户端发送请求
// 参数：req - 要执行的HTTP请求
// 返回值：HTTP响应和可能的错误
func (d *Dao) Do(req *http.Request) (*http.Response, error) {
	return d.httpCli.Do(req)
}

// RetryRequest 带重试机制的HTTP请求
// 在请求失败时自动重试，提高请求成功率
// 参数：
//
//	tries - 最大重试次数
//	request - 要执行的HTTP请求
//
// 返回值：成功的HTTP响应或最终的错误
func (d *Dao) RetryRequest(tries int, request *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// 循环重试指定次数
	for try := 0; try < tries; try++ {
		resp, err = d.Do(request)

		// 如果网络请求失败，等待1秒后重试
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// 如果HTTP状态码不是200，关闭响应体并等待2秒后重试
		// if resp.StatusCode != http.StatusOK {
		// 	resp.Body.Close()
		// 	time.Sleep(2 * time.Second)
		// 	continue
		// }

		// 请求成功，触发回调（如果设置了）
		if d.requestCallback != nil {
			d.requestCallback()
		}

		// 返回响应
		return resp, err
	}

	// 所有重试都失败，返回错误信息
	return resp, Errors.ErrNetwork
}

// SetRequestCallback 设置HTTP请求成功回调
// 用于外部监控每次成功的HTTP请求，通常用于统计请求次数
// 参数：callback - 回调函数，每次HTTP请求成功后调用
func (d *Dao) SetRequestCallback(callback func()) {
	d.requestCallback = callback
}

// CheckLogin 检查指定URL的登录状态
// 通过检查该域名是否存在登录Cookie来判断登录状态
// 参数：ul - 要检查的URL地址
// 返回值：true表示已登录，false表示未登录
func (d *Dao) CheckLogin(ul string) bool {
	// 解析URL获取主机名
	ur, _ := u.Parse(ul)

	// 检查该主机是否存在登录Cookie信息
	if _, ok := d.credentials.LoginCookies[ur.Host]; ok {
		return true
	}
	return false
}

func (d *Dao) GetProxy() string {
	return d.proxy
}

// SetProxy 切换代理设置
// 允许在不改变登录状态的情况下切换网络代理
// 参数：newProxy - 新的代理地址，空字符串表示不使用代理
func (d *Dao) SetProxy(newProxy string) {
	// 获取或创建对应代理的 transport
	transport := getOrCreateTransport(newProxy)

	// 更新代理配置和 HTTP 客户端的 transport
	d.proxy = newProxy
	d.httpCli.Transport = transport
}

// GetCookiesString 获取指定URL对应的登录Cookie信息
// 根据URL的主机名查找对应的Cookie数据
// 参数：ul - URL地址
// 返回值：登录Cookie信息或nil(如果不存在)
func (d *Dao) GetCookiesString(ul string) *LoginCookie {
	// 解析URL获取主机名
	ur, _ := u.Parse(ul)

	// 查找并返回该主机对应的Cookie信息
	if cookeStr, ok := d.credentials.LoginCookies[ur.Host]; ok {
		return cookeStr
	}
	return nil
}

// getOrCreateTransport 获取或创建指定代理的HTTP Transport
// 从全局缓存中查找，如果不存在则创建新的Transport并缓存
// 参数：proxy - 代理服务器地址，空字符串表示不使用代理
// 返回值：配置好的HTTP Transport对象
func getOrCreateTransport(proxy string) *http.Transport {
	// 尝试从缓存中加载已有的 transport
	if t, ok := globalDaos.daos.Load(proxy); ok {
		return t.(*http.Transport)
	}

	// 缓存中不存在，创建新的 transport
	// 代理函数配置，根据传入的proxy参数决定是否使用代理
	proxyFn := func(_ *http.Request) (*u.URL, error) {
		if proxy == "" {
			return nil, nil // 不使用代理
		}
		// 解析代理地址并返回URL对象
		ul, err := u.Parse("http://" + proxy)
		return ul, err
	}

	transport := &http.Transport{
		Proxy:        proxyFn,                                                                // 代理配置
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper), // 禁用HTTP/2
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,  // 连接超时时间
			KeepAlive: 90 * time.Second, // 保持连接时间
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second, // TLS握手超时时间

		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过TLS证书验证（生产环境建议移除）
		},
		DisableCompression:  true, // 禁用压缩
		DisableKeepAlives:   true, // 禁用长连接
		MaxIdleConns:        100,  // 最大空闲连接数
		MaxIdleConnsPerHost: 10,   // 每个主机最大空闲连接数
		MaxConnsPerHost:     20,   // 每个主机最大连接数
	}

	// 存入缓存供后续复用
	globalDaos.daos.Store(proxy, transport)
	return transport
}

// New 创建新的Dao实例
// 初始化HTTP客户端和相关配置，支持代理设置
// 参数：proxy - 代理服务器地址，空字符串表示不使用代理
// 返回值：配置完成的Dao实例
func New(proxy string) *Dao {
	// 获取或创建对应代理的 transport
	transport := getOrCreateTransport(proxy)

	// 创建Cookie存储对象，用于自动管理HTTP Cookie
	jar, _ := cookiejar.New(nil)
	return &Dao{
		proxy: proxy,
		httpCli: &http.Client{
			Jar:       jar,                 // 设置Cookie存储
			Transport: transport,            // 使用缓存的transport
			Timeout:   10 * time.Second,     // 整体请求超时时间
		},
		credentials: &Credentials{}, // 初始化空的用户凭据
	}
}
