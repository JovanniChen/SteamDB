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
	{Username: "za0ww9ml4xl2", Password: "HLHxGyRMm6Zi", SharedSecret: "F54xOr9Tpyd5fAxgKx+RHR7vHik="}, // [0] [xv6753] [46]
	{Username: "zytmnd2097", Password: "awtekBcEkXz9", SharedSecret: "vNVDHuqBle/rnsG7EQW2xQUqlME="},   // [1] [4wzwg]  [45]
	{Username: "zwrvsq6897", Password: "5uoIBclSSBI8", SharedSecret: "kUcQLn0pJutKt9oeh8yRDG7t+o8="},   // [2] [wqrmhz] [44]
	{Username: "zuzuaw8238", Password: "uYj035ynLA5N", SharedSecret: "yKuRsv/OmI584XxMt2LUWWbCM+Y="},   // [3] [kxweoq] [40]
	{Username: "yrknu899", Password: "FyoR1QV8brUd", SharedSecret: "q8JcjcE5jc65C7YntMrME8HJ9sY="},     // [4] [3zgmh7] [47]
	{Username: "mbkle379", Password: "CFs91IvocA39", SharedSecret: "sIF2wljQzxzya9xVO/VtEs1pUwc="},     // [5] [x5x3g8] [48]
	{Username: "lvqxpe8572", Password: "5gfweOafGM3S", SharedSecret: "QLWiEAN8ebHLkGtt7HHtuZyMwDg="},   // [6] [x5x3g8] [49]
	{Username: "naotqp7801", Password: "ja9C5LZelku0", SharedSecret: "g+kIH7JuL98R5O00j87379CkFus="},   // [7] [x5x3g8] [49]
	{Username: "iatfqv6444", Password: "NOJsp0b1aqbj", SharedSecret: "wCdOSNrhPjXrJEpg3FX643+fseQ="},   // [8] [x5x3g8] [49]
	{Username: "uwxhfw8800", Password: "ybuYe33Qg2Dr", SharedSecret: "ViUburoMwWe88QJfL5f0KPPoY68="},   // [9] [x5x3g8] [49]
	{Username: "xqkea03549", Password: "wuwQJ5WFdZp1", SharedSecret: "59z0KMWJFdgfWrSgYYADD/LBPyU="},   // [10] [6ck2bcax] [53]
	{Username: "ffotd74229", Password: "oP4M4CMHAftX", SharedSecret: "IDhBX3NM+8fZCti4C3d6oFhXI6E="},   // [11] [x5x3g8] [54]
	{Username: "ttmsq72777", Password: "yoRD7x6LQvgu", SharedSecret: "5boHTiGFhQoszGcpFDLB7H7thng="},   // [12] [x5x3g8] [52]
	{Username: "ddndd12412", Password: "New0KJYVv16", SharedSecret: "VoSY5VrnD+CJooEVrlADofTGTok="},    // [13] [x5x3g8] [51]
}

// var config *Steam.Config = Steam.NewConfig("your_username:your_password@13.52.178.34:8080")

var config *Steam.Config = Steam.NewConfig("")

// main 主函数，程序入口点
// 执行Steam平台相关操作的演示流程
func main() {
	//TestGetTokenCode(13)
	// TestLogin(12)
	// TestGetSummary(7)
	// TestGetInventory(10)
	//TestGetMyListings(13)

	TestPutList(10)
	// TestBuyListing(13)

	// TestPutList2(5)
	// TestGetConfirmations(12)
	// TestRemoveMyListings(4)
	// TestGetBalance(0)
	// TestGetWaitBalance(3)
	// TestGetInventoryAndPutList(4)
	// TestCreateOrder(10)
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
	fmt.Println("开始测试登录")
	account := getAccount(accountIndex)

	client, err := Steam.NewClient(config)
	if err != nil {
		Logger.Error(err)
		return
	}

	maFile, err := os.ReadFile("mafiles/" + account.Username + ".maFile")
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
	Logger.Info(client.GetInventory(Constants.TeamFortress2, Constants.Catetory))
}

func TestPutList(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	items, err := client.GetInventory(Constants.TeamFortress2, Constants.Catetory)
	if err != nil {
		Logger.Error(err)
		return
	}

	if len(items) == 0 {
		Logger.Error("无可用库存")
		return
	}

	// 随机
	randomIndex := rand.Intn(len(items))
	randomItem := items[randomIndex]

	data, err := os.ReadFile("mafiles/" + client.GetUsername() + ".maFile")
	if err != nil {
		Logger.Error(err)
		return
	}

	listingIds, err := client.PutList(randomItem.AssetID, 0.14, 23, string(data))
	if err != nil {
		Logger.Error(err)
		return
	}

	for _, listingId := range listingIds {
		Logger.Debug("上架成功：", listingId)
	}

}

func TestBuyListing(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	data, err := os.ReadFile("mafiles/" + client.GetUsername() + ".maFile")
	if err != nil {
		return
	}

	maFileContent := string(data)
	//563637814507189205
	//563637814507189205

	//563637814507222895
	Logger.Info(client.BuyListing("625562309703857739", "", 0.16, 0.14, maFileContent))
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
	data, err := os.ReadFile("mafiles/" + client.GetUsername() + ".maFile")
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

	data, err := os.ReadFile("mafiles/" + client.GetUsername() + ".maFile")
	if err != nil {
		return
	}

	maFileContent := string(data)

	Logger.Info(client.CreateOrder("Giftapult", 0.12, 10, maFileContent))
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
	os.WriteFile(fmt.Sprintf("temp/session_%d.json", accountIndex), json, 0644)
}

func (s *SteamSession) Load(accountIndex int) {
	data, _ := os.ReadFile(fmt.Sprintf("temp/session_%d.json", accountIndex))
	json.Unmarshal(data, s)
}

func (s *SteamSession) IsExist(accountIndex int) bool {
	_, err := os.ReadFile(fmt.Sprintf("temp/session_%d.json", accountIndex))
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
