// user.go - Steam用户信息管理功能
// 提供用户信息获取、语言设置、Cookie管理等用户相关操作
package Dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Errors"
	"github.com/JovanniChen/SteamDB/Steam/Param"
	"github.com/JovanniChen/SteamDB/Steam/Utils"
	"github.com/antchfx/htmlquery"
)

// UserInfo 用户信息结构体
// 包含Steam用户的基本信息和账户状态
type UserInfo struct {
	Balance     int    // 钱包余额(分)
	WaitBalance int    // 待处理余额
	Point       int    // Steam积分
	PersonName  string // 用户昵称
	CountryCode string // 国家代码
	Language    string // 语言设置
}

// userInfo 解析HTML页面获取用户信息
// 从Steam账户页面的HTML中提取用户相关信息
// 参数：body - HTTP响应体(HTML格式)
// 返回值：用户信息结构体和可能的错误
func (d *Dao) userInfo(body io.ReadCloser) (*UserInfo, error) {
	// 解析HTML文档
	doc, err := htmlquery.Parse(body)
	if err != nil {
		return nil, err
	}

	info := &UserInfo{}

	// 提取钱包余额信息
	//div[@class="accountData price"]/a/text()
	for _, name := range htmlquery.Find(doc, `//a[@id="header_wallet_balance"]/text()`) {
		// fmt.Println("info.Balance.name ->", name.Data)
		// info.Balance, _ = strconv.Atoi(strings.TrimSpace(name.Data))
		info.Balance = Utils.WalletConvert(name.Data)
	}

	// 提取用户昵称
	for _, name := range htmlquery.Find(doc, `//*[@id="account_pulldown"]/text()`) {
		info.PersonName = strings.TrimSpace(name.Data)
	}

	// 提取待处理余额
	//a[@id="header_wallet_balance"]/text()
	for _, name := range htmlquery.Find(doc, `//a[@id="header_wallet_balance"]/span/text()`) {
		// fmt.Println("info.WaitBalance.name ->", name.Data)
		// info.WaitBalance, _ = strconv.Atoi(strings.TrimSpace(name.Data))
		info.WaitBalance = Utils.WalletConvert(name.Data)
	}

	// 提取应用配置信息(语言、国家等)
	for _, name := range htmlquery.Find(doc, `//div[@id='application_config']`) {
		src := htmlquery.SelectAttr(name, "data-config")
		m := make(map[string]interface{})
		err = json.Unmarshal([]byte(src), &m)
		if err != nil {
			return nil, err
		}
		info.Language = m["LANGUAGE"].(string)
		info.CountryCode = m["COUNTRY"].(string)
	}

	return info, nil
}

// getUserInfo 获取用户详细信息
// 通过访问Steam账户页面获取用户的详细信息
// 返回值：用户信息结构体和可能的错误
func (d *Dao) getUserInfo() (*UserInfo, error) {
	accountUrl := Constants.Account

	// 检查是否已登录该域名
	if d.CheckLogin(accountUrl) {
		// 创建认证请求
		req, err := d.Request("GET", accountUrl, nil)
		if err != nil {
			return nil, err
		}

		// 发送请求获取响应
		resp, err := d.RetryRequest(Constants.Tries, req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// 检查响应状态
		if resp.StatusCode != 200 {
			return nil, Errors.ResponseError(resp.StatusCode)
		}

		// 解析HTML获取用户信息
		return d.userInfo(resp.Body)
	}

	return nil, nil
}

func (d *Dao) GetSteamIDByFriendLink(friendLink string) (uint64, error) {
	// 创建认证请求
	req, err := d.Request("GET", friendLink, nil)
	if err != nil {
		return 0, err
	}

	// 发送请求获取响应
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != 200 {
		return 0, Errors.ResponseError(resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	re := regexp.MustCompile(`"steamid":"(\d+)"`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return 0, fmt.Errorf("SteamID not found in the page")
	}

	steamID, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return steamID, nil
}

// SetCookiesLanguage 设置Cookie中的语言偏好
// 更新内部凭据中的语言设置
// 参数：language - 语言代码(如"schinese", "english"等)
func (d *Dao) SetCookiesLanguage(language string) {
	if d.credentials != nil {
		d.credentials.Language = language
	}
}

// UserInfo 更新用户详细信息到内部凭据
// 获取并更新用户的国家代码、语言、昵称等信息
// 返回值：操作成功返回nil，失败返回错误
func (d *Dao) UserInfo() error {
	// 检查credentials是否已初始化
	if d.credentials == nil {
		return errors.New("credentials未初始化，请先执行登录操作")
	}

	// 获取用户信息
	info, err := d.getUserInfo()
	if err != nil {
		return err
	}

	// 更新内部凭据
	d.credentials.CountryCode = info.CountryCode
	d.credentials.Language = info.Language
	d.credentials.Nickname = info.PersonName
	d.SetCookiesLanguage(info.Language)

	return nil
}

// GetUserCookies 获取所有域名的登录Cookie
// 返回JSON格式的Cookie数据，用于持久化存储或传输
// 返回值：JSON格式的Cookie字节数组和可能的错误
func (d *Dao) GetUserCookies() ([]byte, error) {
	return json.Marshal(d.credentials.LoginCookies)
}

// SetLanguage 设置Steam界面语言
// 向Steam发送语言设置请求，更改用户界面语言
// 参数：language - 语言代码(如"schinese"表示简体中文)
// 返回值：设置成功返回nil，失败返回错误
func (d *Dao) SetLanguage(language string) error {
	if !d.CheckLogin(Constants.Language) {
		return fmt.Errorf("未登录")
	}

	// 构建POST请求参数
	params := Param.Params{}
	params.SetString("language", language)
	params.SetString("sessionid", d.GetCookiesString(Constants.Language).SessionId) // 需要会话ID验证

	// 创建POST请求
	req, err := d.Request("POST", Constants.Language, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	// 发送请求
	resp, err := d.RetryRequest(Constants.Tries, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != 200 {
		return fmt.Errorf("语言设置失败: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(body), "true") {
		return fmt.Errorf("语言设置失败: %s", string(body))
	}
	d.SetCookiesLanguage(language)
	return nil
}
