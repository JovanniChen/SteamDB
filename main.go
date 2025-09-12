// Steam数据库操作主程序
// 本程序用于连接Steam平台，进行用户登录、获取令牌代码和添加反应等操作
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/JovanniChen/SteamDB/Steam"
	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Dao"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
)

var accounts = []Account{
	{Username: "za0ww9ml4xl2", Password: "HLHxGyRMm6Zi", SharedSecret: "F54xOr9Tpyd5fAxgKx+RHR7vHik="}, // [0] [xv6753] [46]
	{Username: "zytmnd2097", Password: "awtekBcEkXz9", SharedSecret: "vNVDHuqBle/rnsG7EQW2xQUqlME="},   // [1] [4wzwg]  [45]
	{Username: "zwrvsq6897", Password: "5uoIBclSSBI8", SharedSecret: "kUcQLn0pJutKt9oeh8yRDG7t+o8="},   // [2] [wqrmhz] [44]
	{Username: "zuzuaw8238", Password: "uYj035ynLA5N", SharedSecret: "yKuRsv/OmI584XxMt2LUWWbCM+Y="},   // [3] [kxweoq] [40]
	{Username: "yrknu899", Password: "FyoR1QV8brUd", SharedSecret: "q8JcjcE5jc65C7YntMrME8HJ9sY="},     // [4] [3zgmh7] [47]
	{Username: "mbkle379", Password: "CFs91IvocA39", SharedSecret: "sIF2wljQzxzya9xVO/VtEs1pUwc="},     // [5] [x5x3g8] [48]
}

// main 主函数，程序入口点
// 执行Steam平台相关操作的演示流程
func main() {
	// TestGetTokenCode(2)
	// TestLogin(4)
	// TestGetSummary(4)
	// TestGetInventory(4)

	TestGetMyListings(4)
	// TestPutList(4)
	// TestGetConfirmations(4)
	// TestRemoveMyListings(4)
	// TestBuyListing(2)
}

func TestGetTokenCode(accountIndex int) {
	client, err := Steam.NewClient(Steam.NewConfig(""))
	if err != nil {
		Logger.Error(err)
		return
	}
	code, _ := client.GetTokenCode(getAccount(accountIndex).GetSharedSecret())
	Logger.Info(code)
}

func TestLogin(accountIndex int) {
	account := getAccount(accountIndex)

	client, err := Steam.NewClient(Steam.NewConfig(""))
	if err != nil {
		Logger.Error(err)
		return
	}

	userInfo, err := client.Login(&Steam.LoginCredentials{
		Username:     account.GetUsername(),
		Password:     account.GetPassword(),
		SharedSecret: account.GetSharedSecret(),
	})
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(userInfo)

	// 提取访问令牌
	accessToken, err := client.GetAccessToken()
	if err != nil {
		accessToken = ""
	}

	// 提取刷新令牌
	refreshToken := ""
	if rt := client.GetRefreshToken(); rt != "" {
		refreshToken = rt
	}

	// 提取登录Cookies
	loginCookies := make(map[string]*Dao.LoginCookie)
	if cookies := client.GetLoginCookies(); cookies != nil {
		loginCookies = cookies
	}
	session := &SteamSession{
		AccountIndex: accountIndex,
		Username:     account.GetUsername(),
		SteamID:      client.GetSteamID(),
		Nickname:     client.GetNickname(),
		CountryCode:  client.GetCountryCode(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		LoginCookies: loginCookies,
		LoginTime:    time.Now(),
	}
	session.Save(accountIndex)
}

func TestGetSummary(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	summary, err := client.GetPointsSummary(client.GetSteamID())
	Logger.Info(summary)
}

func TestGetMyListings(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetMyListings())
}

func TestGetInventory(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetInventory(Constants.Dota2, Constants.Catetory))
}

func TestPutList(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	// 32169268728
	Logger.Info(client.PutList("32169283736", 0.1, 1))
}

func TestBuyListing(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.BuyListing("654831914925835322", "Skirt of the Mage Slayer"))
}

func TestRemoveMyListings(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.RemoveMyListings("654831914925591572"))
}

func TestGetConfirmations(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetConfirmations())
}

func loadFromSession(accountIndex int) (*Steam.Client, error) {
	session := &SteamSession{}
	session.Load(accountIndex)
	client, err := Steam.NewClient(Steam.NewConfig(""))
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	client.SetLoginInfo(session.Username, session.SteamID, session.Nickname, session.CountryCode, session.AccessToken, session.RefreshToken, session.LoginCookies)
	return client, nil
}

// SteamSession Steam会话信息
type SteamSession struct {
	AccountIndex int                         `json:"account_index"`
	Username     string                      `json:"username"`
	SteamID      uint64                      `json:"steam_id"`
	Nickname     string                      `json:"nickname"`
	CountryCode  string                      `json:"country_code"`
	AccessToken  string                      `json:"access_token"`
	RefreshToken string                      `json:"refresh_token"`
	LoginCookies map[string]*Dao.LoginCookie `json:"login_cookies"`
	LoginTime    time.Time                   `json:"login_time"`
}

func (s *SteamSession) Save(accountIndex int) {
	json, _ := json.Marshal(s)
	os.WriteFile(fmt.Sprintf("session_%d.json", accountIndex), json, 0644)
}

func (s *SteamSession) Load(accountIndex int) {
	data, _ := os.ReadFile(fmt.Sprintf("session_%d.json", accountIndex))
	json.Unmarshal(data, s)
}

func (s *SteamSession) IsExist(accountIndex int) bool {
	_, err := os.ReadFile(fmt.Sprintf("session_%d.json", accountIndex))
	return err == nil
}

type Account struct {
	Username     string // Steam用户名
	Password     string // Steam密码
	SharedSecret string // Steam Guard共享密钥(base64编码)
}

func (a *Account) GetUsername() string {
	return a.Username
}
func (a *Account) GetPassword() string {
	return a.Password
}
func (a *Account) GetSharedSecret() string {
	return a.SharedSecret
}

func getAccount(index int) *Account {
	return &accounts[index]
}
