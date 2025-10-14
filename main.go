// Steam数据库操作主程序
// 本程序用于连接Steam平台，进行用户登录、获取令牌代码和添加反应等操作
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/JovanniChen/SteamDB/Steam"
	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Dao"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
)

var accounts = []Account{
	{Username: "za0ww9ml4xl2", Password: "HLHxGyRMm6Zi", SharedSecret: "F54xOr9Tpyd5fAxgKx+RHR7vHik="},   // [0] [xv6753] [46]
	{Username: "zytmnd2097", Password: "awtekBcEkXz9", SharedSecret: "vNVDHuqBle/rnsG7EQW2xQUqlME="},     // [1] [4wzwg]  [45]
	{Username: "zwrvsq6897", Password: "5uoIBclSSBI8", SharedSecret: "kUcQLn0pJutKt9oeh8yRDG7t+o8="},     // [2] [wqrmhz] [44]
	{Username: "zuzuaw8238", Password: "uYj035ynLA5N", SharedSecret: "yKuRsv/OmI584XxMt2LUWWbCM+Y="},     // [3] [kxweoq] [40]
	{Username: "yrknu899", Password: "FyoR1QV8brUd", SharedSecret: "q8JcjcE5jc65C7YntMrME8HJ9sY="},       // [4] [3zgmh7] [47]
	{Username: "mbkle379", Password: "CFs91IvocA39", SharedSecret: "sIF2wljQzxzya9xVO/VtEs1pUwc="},       // [5] [x5x3g8] [48]
	{Username: "ugsxh51037", Password: "z0dAC0nic9Ec", SharedSecret: "KXdQ/El9khZe6K3HIxLS7IwrDi4="},     // [6] [x5x3g8] [49]
	{Username: "cv71oebl0wvj6z", Password: "uolMwmIPT8Uo", SharedSecret: "ireKAD4ZX7HfC45M23iKiYiobqU="}, // [67 [x5x3g8] [49]

}

// var config *Steam.Config = Steam.NewConfig("your_username:your_password@54.215.254.6:8080")

var config *Steam.Config = Steam.NewConfig("")

// main 主函数，程序入口点
// 执行Steam平台相关操作的演示流程
func main() {
	// TestGetTokenCode(6)
	TestLogin(5)
	// TestGetSummary(7)
	// TestGetInventory(7)
	// TestGetMyListings(3)
	// TestPutList(5)
	TestPutList2(5)
	// TestGetConfirmations(3)
	// TestRemoveMyListings(4)
	// TestBuyListing(4)
	// TestGetBalance(7)
	// TestGetWaitBalance(3)
	// TestGetInventoryAndPutList(4)
	// TestCreateOrder(3)
}

func TestGetTokenCode(accountIndex int) {
	client, err := Steam.NewClient(config)
	if err != nil {
		Logger.Error(err)
		return
	}
	code, _ := client.GetTokenCode(getAccount(accountIndex).GetSharedSecret())
	Logger.Info(code)
}

func TestLogin(accountIndex int) {
	account := getAccount(accountIndex)

	client, err := Steam.NewClient(config)
	if err != nil {
		Logger.Error(err)
		return
	}

	maFile, err := os.ReadFile(account.Username + ".maFile")
	if err != nil {
		return
	}

	userInfo, err := client.Login(&Steam.LoginCredentials{
		Username:     account.GetUsername(),
		Password:     account.GetPassword(),
		SharedSecret: account.GetSharedSecret(),
		MaFile:       string(maFile),
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
	Logger.Info("GetSummary -> ", summary)
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

	items, err := client.GetInventory(Constants.Dota2, Constants.Catetory)
	if err != nil {
		Logger.Error(err)
		return
	}

	// 随机
	randomIndex := rand.Intn(len(items))
	randomItem := items[randomIndex]

	data, err := os.ReadFile(client.GetUsername() + ".maFile")
	if err != nil {
		Logger.Error(err)
		return
	}

	if err := client.PutList(randomItem.AssetID, 2, 10, string(data)); err != nil {
		Logger.Error(err)
		return
	}

	client.GetMyListings()
}

func TestPutList2(accountIndex int) {
	// 32849705541
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	data, err := os.ReadFile(client.GetUsername() + ".maFile")
	if err != nil {
		Logger.Error(err)
		return
	}

	if err := client.PutList("32849705541", 2, 10, string(data)); err != nil {
		Logger.Error(err)
		return
	}
}

func TestGetInventoryAndPutList(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	items, err := client.GetInventory(Constants.Dota2, Constants.Catetory)
	if err != nil {
		Logger.Error("获取库存失败，错误：", err)
		return
	}

	if len(items) == 0 {
		Logger.Warn("该用户没有可用库存")
		return
	}

	if err := client.PutList(items[0].AssetID, 0.1, 1, ""); err != nil {
		Logger.Error("上架失败失败，错误：", err)
	}

	Logger.Info("上架成功!")
}

func TestBuyListing(accountIndex int) {
	// client, err := loadFromSession(accountIndex)
	// if err != nil {
	// 	Logger.Error(err)
	// 	return
	// }
	// Logger.Info(client.BuyListing("640197019947684853", "Skirt of the Mage Slayer"))
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
	data, err := os.ReadFile(client.GetUsername() + ".maFile")
	if err != nil {
		return
	}

	maFileContent := string(data)

	Logger.Info(client.GetConfirmations(maFileContent))
}

func TestGetBalance(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetBalance())
}

func TestGetWaitBalance(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetWaitBalance())
}

func TestCreateOrder(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	data, err := os.ReadFile(client.GetUsername() + ".maFile")
	if err != nil {
		return
	}

	maFileContent := string(data)

	Logger.Info(client.CreateOrder("Heat of the Sixth Hell - Back", 0.05, 5, maFileContent))
}

func loadFromSession(accountIndex int) (*Steam.Client, error) {
	session := &SteamSession{}
	session.Load(accountIndex)
	client, err := Steam.NewClient(config)
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
