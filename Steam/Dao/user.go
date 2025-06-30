package Dao

import (
	"bytes"
	"encoding/json"
	"example.com/m/v2/Steam"
	"example.com/m/v2/Steam/Errors"
	"example.com/m/v2/Steam/Param"
	"github.com/antchfx/htmlquery"
	"io"
	"strconv"
	"strings"
)

type UserInfo struct {
	Balance     int
	WaitBalance int
	Point       int
	PersonName  string
	CountryCode string
	Language    string
}

func (d *Dao) userInfo(body io.ReadCloser) (*UserInfo, error) {
	// 解析html
	doc, err := htmlquery.Parse(body)
	if err != nil {
		return nil, err
	}
	info := &UserInfo{}
	for _, name := range htmlquery.Find(doc, `//div[@class="accountData price"]/a/text()`) {
		info.Balance, _ = strconv.Atoi(strings.TrimSpace(name.Data))
	}
	for _, name := range htmlquery.Find(doc, `//*[@id="account_pulldown"]/text()`) {
		info.PersonName = strings.TrimSpace(name.Data)
	}
	for _, name := range htmlquery.Find(doc, `//a[@id="header_wallet_balance"]/text()`) {
		info.WaitBalance, _ = strconv.Atoi(strings.TrimSpace(name.Data))
	}
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

// getUserInfo 获取用户明细
func (d *Dao) getUserInfo() (*UserInfo, error) {
	accountUrl := Steam.Account
	if d.CheckLogin(accountUrl) {
		req, err := d.Request("GET", accountUrl, nil)
		if err != nil {
			return nil, err
		}
		resp, err := d.RetryRequest(Steam.Tries, req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, Errors.ResponseError(resp.StatusCode)
		}
		return d.userInfo(resp.Body)
	}
	return nil, nil
}

// SetCookiesLanguage 设置语言
func (d *Dao) SetCookiesLanguage(language string) {
	d.credentials.Language = language
}

// UserInfo 获取用户详细信息
func (d *Dao) UserInfo() error {
	info, err := d.getUserInfo()
	if err != nil {
		return err
	}
	d.credentials.CountryCode = info.CountryCode
	d.credentials.Language = info.Language
	d.credentials.Nickname = info.PersonName
	d.SetCookiesLanguage(info.Language)
	return nil
}

// GetUserCookies 获取所有域名登录后的cookies
func (d *Dao) GetUserCookies() ([]byte, error) {
	return json.Marshal(d.credentials.LoginCookies)
}

// SetLanguage 设置语言为英语
func (d *Dao) SetLanguage(language string) error {
	languageUrl := Steam.Language
	if d.CheckLogin(languageUrl) {
		lg := d.GetCookiesString(languageUrl)
		params := Param.Params{}
		params.SetString("language", language)
		params.SetString("sessionid", lg.SessionId)
		req, err := d.Request("POST", languageUrl, strings.NewReader(params.Encode()))
		if err != nil {
			return err
		}
		resp, err := d.RetryRequest(Steam.Tries, req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return Errors.ResponseError(resp.StatusCode)
		}
		buf := new(bytes.Buffer)
		if _, err = buf.ReadFrom(resp.Body); err != nil {
			return err
		}
		if buf.String() != "true" {
			return Errors.Error("语言设置失败")
		}
	}
	return nil
}
