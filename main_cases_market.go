package main

import (
	"os"
	"time"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
)

func TestUnsendGift(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	items, err := client.GetSteamGift(Constants.Steam, Constants.SteamGiftCategory)
	if err != nil {
		Logger.Error(err)
		return
	}

	for _, item := range items {
		time.Sleep(1 * time.Second)
		if err := client.UnsendGift(item.AssetID); err != nil {
			Logger.Errorf("撤回赠送礼物失败[%s]: %v", item.AssetID, err)
			continue
		}
		Logger.Info("撤回赠送礼物成功: ", item.AssetID)
	}
}

func TestUnsendAllGift(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.UnsendAllGift())
}

func TestGetMyListings(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	client.GetMyListings()
	// Logger.Infof("已上架物品 (%d 个) -> %+v\n", len(activeListings), activeListings)
}

func TestGetMarketListings(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	response, err := client.GetMarketListings(440, "Giftapult", 500, 100, "CN", "schinese", 23)
	if err != nil {
		Logger.Error(err)
		return
	}

	if !response.Success {
		Logger.Error("获取市场列表失败")
		return
	}

	Logger.Info("获取市场列表成功")
	for _, listing := range response.ListingInfo {
		Logger.Infof("ListingID: %s, ConvertedSteamFee: %d, ConvertedPublisherFee: %d, ConvertedPricePerUnit: %d", listing.ListingID, listing.ConvertedSteamFee, listing.ConvertedPublisherFee, listing.ConvertedPricePerUnit)
	}
}

func TestGetInventory(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	items, err := client.GetInventory(Constants.Steam, Constants.SteamCategory)
	if err != nil {
		Logger.Error("获取库存失败: ", err)
	}

	for _, item := range items {
		Logger.Infof("物品ID: %s, 名称: %s, 市场名称: %s, 价格: %f, 货币: %d, 是否可交易: %t, 是否可在市场交易: %t", item.AssetID, item.Name, item.MarketName, item.Price, item.Currency, item.Tradable, item.Marketable)
	}
}

func TestGetSteamGift(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	items, err := client.GetSteamGift(Constants.Steam, Constants.SteamGiftCategory)
	if err != nil {
		Logger.Error("获取库存失败: ", err)
	}

	for _, item := range items {
		Logger.Infof("物品ID: %s, 名称: %s, 市场名称: %s, 价格: %f, 货币: %d, 是否可交易: %t, 是否可在市场交易: %t", item.AssetID, item.Name, item.MarketName, item.Price, item.Currency, item.Tradable, item.Marketable)
	}
}

func TestPutList(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	data, err := os.ReadFile("mafiles/" + client.GetUsername() + ".maFile")
	if err != nil {
		Logger.Error(err)
		return
	}

	items, err := client.GetInventory(Constants.TeamFortress2, Constants.TeamFortress2Catetory)
	if err != nil {
		Logger.Error(err)
		return
	}

	response, err := client.PutList(Constants.TeamFortress2, Constants.TeamFortress2Catetory, items[0].AssetID, 1, 23, string(data))
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.Info("上架成功")
	Logger.Infof("上架物品: %+v", response)
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

	if err := client.BuyListing(440, "631208307256406320", "", 21, 7, maFileContent); err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info("购买成功")
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

	Logger.Info(client.CreateOrder("Giftapult", 0.12, 15, maFileContent))
}
