package main

import "github.com/JovanniChen/SteamDB/Steam/Logger"

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
