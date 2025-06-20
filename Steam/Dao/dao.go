package Dao

import (
	"crypto/tls"
	"example.com/m/v2/Steam"
	"example.com/m/v2/Steam/Errors"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	u "net/url"
	"time"
)

type Dao struct {
	httpCli     *http.Client
	credentials *Credentials
}

func (d *Dao) Request(method, url string, body io.Reader) (*http.Request, error) {
	req, err := d.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	ur, _ := u.Parse(url)
	ck := d.GetCookiesString(url)
	if ck == nil {
		return nil, Errors.Error("Cookie not exist")
	}
	req.AddCookie(&http.Cookie{
		Name:   "steamLoginSecure",
		Value:  ck.SteamLoginSecure,
		Domain: ur.Host,
		Path:   "/",
	})
	req.AddCookie(&http.Cookie{
		Name:   "sessionid",
		Value:  ck.SessionId,
		Domain: ur.Host,
		Path:   "/",
	})
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

func (d *Dao) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", Steam.UseAgent)
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
}

func (d *Dao) Do(req *http.Request) (*http.Response, error) {
	return d.httpCli.Do(req)
}

func (d *Dao) RetryRequest(tries int, request *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	for try := 0; try < tries; try++ {
		resp, err = d.Do(request)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			time.Sleep(2 * time.Second)
			continue
		}
		return resp, err
	}
	return resp, Errors.Error("多次请求后还是失败-->" + request.URL.String())
}

func (d *Dao) CheckLogin(ul string) bool {
	ur, _ := u.Parse(ul)
	if _, ok := d.credentials.LoginCookies[ur.Host]; ok {
		return true
	}
	return false
}

func (d *Dao) GetCookiesString(ul string) *LoginCookie {
	ur, _ := u.Parse(ul)
	if cookeStr, ok := d.credentials.LoginCookies[ur.Host]; ok {
		return cookeStr
	}
	return nil
}

func New(proxy string) *Dao {
	proxyFn := func(_ *http.Request) (*u.URL, error) {
		if proxy == "" {
			return nil, nil
		}
		ul, err := u.Parse("http://" + proxy) //根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
		return ul, err
	}
	// 创建一个 Cookie 存储对象
	jar, _ := cookiejar.New(nil)

	return &Dao{
		httpCli: &http.Client{
			//CheckRedirect: func(req *http.Request, via []*http.Request) error {
			//	// 获取上一次请求的响应，从中提取 Cookie
			//	if len(via) > 0 {
			//		resp := via[len(via)-1].Response
			//		if resp != nil {
			//			for _, cookie := range resp.Cookies() {
			//				// 将 Cookie 添加到下一次请求中
			//				req.AddCookie(cookie)
			//			}
			//		}
			//	}
			//	return nil
			//},
			Jar: jar,
			Transport: &http.Transport{
				Proxy:        proxyFn,
				TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
				Dial: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 30 * time.Second,

				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableCompression:  true,
				DisableKeepAlives:   true,
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				MaxConnsPerHost:     20,
			},
			Timeout: 10 * time.Second,
		},
		credentials: &Credentials{},
	}
}
