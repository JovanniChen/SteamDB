// Steam数据库操作主程序
// 本程序用于连接Steam平台，进行用户登录、获取令牌代码和添加反应等操作
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/JovanniChen/SteamDB/Steam"
	"github.com/JovanniChen/SteamDB/Steam/Dao"
)

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

var accounts = []Account{
	{Username: "za0ww9ml4xl2", Password: "HLHxGyRMm6Zi", SharedSecret: "F54xOr9Tpyd5fAxgKx+RHR7vHik="}, // xv6753
	{Username: "zytmnd2097", Password: "awtekBcEkXz9", SharedSecret: "vNVDHuqBle/rnsG7EQW2xQUqlME="},   // 4wzwg
	{Username: "zwrvsq6897", Password: "5uoIBclSSBI8", SharedSecret: "kUcQLn0pJutKt9oeh8yRDG7t+o8="},   // wqrmhz
	{Username: "zuzuaw8238", Password: "uYj035ynLA5N", SharedSecret: "yKuRsv/OmI584XxMt2LUWWbCM+Y="},   // kxweoq
	{Username: "yrknu899", Password: "FyoR1QV8brUd", SharedSecret: "q8JcjcE5jc65C7YntMrME8HJ9sY="},     // 3zgmh7
}

func getAccount(index int) *Account {
	return &accounts[index]
}

// main 主函数，程序入口点
// 执行Steam平台相关操作的演示流程
func main() {
	// 以下是一些其他功能的示例调用（已注释）
	// result, err := d.GetReacionts(76561198313222178, 3)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(len(result.Reactionids))
	// fmt.Println(result) // 获取反应

	// fmt.Println(d.GetReactionConfig()) // 获取反应配置
	// fmt.Println(d.GetSteamIDByFriendLink("https://steamcommunity.com/user/chbn-qbdd/KCDRCPRT/"))

	// 为指定用户添加反应
	// 参数：用户SteamID，反应类型，反应ID
	// fmt.Println(d.AddReaction(76561198313222178, 3, 23))
	// fmt.Println(d.CheckLogin("https://steamcommunity.com/"))

	// fmt.Println(d.GetMyListings()) // 获取上架列表
	// fmt.Println(d.GetInventory())                   // 获取库存
	// fmt.Println(d.PutList("32169089104", 100, 1)) // 上架物品

	// TestLogin(4)
	// TestGetSummary(4)
	// TestGetMyListings(4)
	// TestGetInventory(4)
	TestPutList(4)
	TestGetConfirmations(4)
}

func TestGetTokenCode(accountIndex int) {
	client, err := Steam.NewClient(Steam.NewConfig(""))
	if err != nil {
		fmt.Println(err)
		return
	}
	code, _ := client.GetTokenCode(getAccount(accountIndex).GetSharedSecret())
	fmt.Println(code)
}

func TestLogin(accountIndex int) {
	account := getAccount(accountIndex)

	client, err := Steam.NewClient(Steam.NewConfig(""))
	if err != nil {
		fmt.Println(err)
		return
	}

	userInfo, err := client.Login(&Steam.LoginCredentials{
		Username:     account.GetUsername(),
		Password:     account.GetPassword(),
		SharedSecret: account.GetSharedSecret(),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(userInfo)

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

	client.GetMyListings()
}

func TestGetSummary(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		fmt.Println(err)
		return
	}
	summary, err := client.GetPointsSummary(client.GetSteamID())
	fmt.Println(summary)
}

func TestGetMyListings(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(client.GetMyListings())
}

func TestGetInventory(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(client.GetInventory())
}

func TestPutList(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(client.PutList("32169075192", 100, 1))
}

func TestGetConfirmations(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(client.GetConfirmations())
}

func loadFromSession(accountIndex int) (*Steam.Client, error) {
	session := &SteamSession{}
	session.Load(accountIndex)
	client, err := Steam.NewClient(Steam.NewConfig(""))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	client.SetLoginInfo(session.Username, session.SteamID, session.Nickname, session.CountryCode, session.AccessToken, session.RefreshToken, session.LoginCookies)
	return client, nil
}
