package main

import (
	"os"
	"strconv"
	"time"

	"github.com/JovanniChen/SteamDB/Steam"
	"github.com/JovanniChen/SteamDB/Steam/Dao"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
)

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
	sessionPath := "temp/session_" + strconv.Itoa(accountIndex) + ".json"
	needLogin := true

	if info, err := os.Stat(sessionPath); err == nil {
		if time.Since(info.ModTime()) > 4*time.Hour {
			if rmErr := os.Remove(sessionPath); rmErr != nil {
				Logger.Error("删除过期session文件失败:", rmErr)
				return
			}
			Logger.Info("session文件已超过4小时，已删除并重新登录")
		} else {
			needLogin = false
			Logger.Info("session文件未超过4小时，跳过登录")
		}
	} else if !os.IsNotExist(err) {
		Logger.Error("检查session文件失败:", err)
		return
	}

	if !needLogin {
		return
	}

	account := getAccount(accountIndex)

	client, err := Steam.NewClient(config)
	if err != nil {
		Logger.Error(err)
		return
	}

	maFile, err := os.ReadFile("mafiles/" + account.Username + ".maFile")
	if err != nil {
		Logger.Info("没有发现maFile文件")
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

	Logger.Info("登录成功")
	Logger.Info(userInfo)

	// 提取访问令牌
	accessToken, err := client.GetAccessToken()
	if err != nil {
		accessToken = ""
	}

	steamOffset := client.GetSteamOffset()

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
		AccountIndex:  accountIndex,
		Username:      account.GetUsername(),
		SteamID:       client.GetSteamID(),
		Nickname:      client.GetNickname(),
		CountryCode:   client.GetCountryCode(),
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		LoginCookies:  loginCookies,
		LoginTime:     time.Now(),
		SteamOffset:   steamOffset,
		SteamLanguage: client.GetLanguage(),
	}
	session.Save(accountIndex)
}

func TestCheckAccountAvailable(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.CheckAccountAvailable(strconv.FormatUint(client.GetSteamID(), 10)))
}

func TestSetLanguage(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.SetLanguage("english"))
}

func TestGetCountryCode(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetCountryCode())
}

func TestLogoutAll(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.LogoutAll())
}

func TestSetPrivacy(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.SetPrivacy())
}
