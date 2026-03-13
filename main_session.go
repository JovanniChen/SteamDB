package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/JovanniChen/SteamDB/Steam"
	"github.com/JovanniChen/SteamDB/Steam/Dao"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
)

func loadFromSession(accountIndex int) (*Steam.Client, error) {
	session := &SteamSession{}
	session.Load(accountIndex)
	client, err := Steam.NewClient(config)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	client.SetLoginInfo(session.Username, session.SteamID, session.Nickname, session.CountryCode, session.AccessToken, session.RefreshToken, session.LoginCookies, session.SteamOffset, session.SteamLanguage)
	return client, nil
}

// SteamSession Steam会话信息
type SteamSession struct {
	AccountIndex  int                         `json:"account_index"`
	Username      string                      `json:"username"`
	SteamID       uint64                      `json:"steam_id"`
	Nickname      string                      `json:"nickname"`
	CountryCode   string                      `json:"country_code"`
	AccessToken   string                      `json:"access_token"`
	RefreshToken  string                      `json:"refresh_token"`
	LoginCookies  map[string]*Dao.LoginCookie `json:"login_cookies"`
	LoginTime     time.Time                   `json:"login_time"`
	SteamOffset   int64                       `json:"steam_offset"`
	SteamLanguage string                      `json:"steam_language"`
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
