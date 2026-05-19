package main

import (
	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
)

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

func TestGetUserInfo(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	userInfo, err := client.GetUserInfo()
	if err != nil {
		Logger.Error(err)
		return
	}

	countryCode := userInfo.CountryCode

	countryInfo := Constants.Countries[countryCode]
	currencySymbol := countryInfo.CurrencySymbol
	Logger.Infof("当前账号余额：%s%.2f ", currencySymbol, float64(userInfo.Balance)/100)
	Logger.Infof("当前账号待处理余额：%s%d ", currencySymbol, userInfo.WaitBalance)
}
