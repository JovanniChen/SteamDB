package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Dao"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
)

func TestIsAccountBanned(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	banned := client.IsAccountBanned()

	Logger.Info("账号是否红信:", banned)
}

func TestGetSteamRate(accountIndex int) {
	start, count := 0, 10
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	resp, err := client.GetMarketListings(440, "Specialized Killstreak L'Etranger Kit Fabricator", start, count, "CN", "schinese", 23)
	if err != nil {
		Logger.Error(err)
		return
	}

	if resp == nil {
		Logger.Info("未获取到市场列表")
		return
	}

	total := resp.TotalCount
	pages := total / count
	if total%count == 0 {
		pages--
	}
	start = pages * count

	type ratio struct {
		code  string
		ratio string
	}

	ratios := make([]ratio, 0, len(Constants.Countries))

	for _, v := range Constants.Countries {
		countryCode := v.Code
		countryName := v.Name
		steamCurrencyID := v.SteamCurrencyID
		resp, err = client.GetMarketListings(440, "Specialized Killstreak L'Etranger Kit Fabricator", start, count, countryCode, "schinese", steamCurrencyID)
		if err != nil {
			Logger.Error(err)
			return
		}

		if resp == nil {
			Logger.Infof("[%s]未获取到市场列表", countryName)
			continue
		}

		Logger.Infof("起始位置: %d, 每页大小: %d, 总数量: %d", resp.Start, resp.PageSize, resp.TotalCount)
		for _, item := range resp.Items {
			if item.AssetID == "16958482695" {
				// Logger.Infof("AssetID: %s, ListingID: %s, Price: %d, Fee: %d, ConvertedPrice: %d, ConvertedFee: %d, ConvertedSteamFee: %d, ConvertedPublisherFee: %d, ConvertedPricePerUnit: %d",
				// 	item.AssetID, item.ListingID, item.Price, item.Fee, item.ConvertedPrice, item.ConvertedFee, item.ConvertedSteamFee, item.ConvertedPublisherFee, item.ConvertedPricePerUnit)
				ratios = append(ratios, ratio{code: countryName, ratio: fmt.Sprintf("%.08f", float64(item.Price)/float64(item.ConvertedPrice))})
			}
		}
	}

	for _, r := range ratios {
		Logger.Infof("%s/人民币汇率: %s", r.code, r.ratio)
	}

}

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

	resp, err := client.GetMyListings()
	if err != nil {
		Logger.Error("获取上架列表失败: ", err)
		return
	}

	Logger.Infof("已上架物品 (%d个) -> %+v", len(resp.Listings), resp.Listings)
	Logger.Infof("待确认上架物品 (%d个) -> %+v", len(resp.ListingsToConfirm), resp.ListingsToConfirm)
}

func TestGetMarketListings(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	response, err := client.GetMarketListings(440, "Specialized%20Killstreak%20Baby%20Face%27s%20Blaster%20Kit%20Fabricator", 0, 10, "CN", "schinese", 23)
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.Info("获取市场列表成功,数量:", len(response.Items))
	Logger.Infof("Start: %d, PageSize: %d, TotalCount: %d", response.Start, response.PageSize, response.TotalCount)
	for _, listing := range response.Items {
		Logger.Infof("AssetID: %s, ListingID: %s, ConvertedSteamFee: %d, ConvertedPublisherFee: %d, ConvertedPricePerUnit: %d",
			listing.AssetID, listing.ListingID, listing.ConvertedSteamFee, listing.ConvertedPublisherFee, listing.ConvertedPricePerUnit)
	}
}

func TestGetInventory(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	items, err := client.GetInventory(Constants.TeamFortress2, Constants.TeamFortress2Catetory, Dao.WithTradable(1), Dao.WithMarketable(1), Dao.WithCommodity(0))
	if err != nil {
		Logger.Error("获取库存失败: ", err)
		return
	}

	Logger.Infof("库存物品 (%d 个)", len(items))

	for _, item := range items {
		Logger.Infof("物品ID: %s, 名称: %s, Icon: %s,市场名称: %s, 市场 Hash名称: %s,价格: %f, 货币: %d, 是否可玩家交易: %t, 是否可在市场交易: %t, 是否标准化商品: %t",
			item.AssetID, item.Name, item.Icon, item.MarketName, item.MarketHashName, item.Price, item.Currency, item.Tradable, item.Marketable, item.Commodity)
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
		Logger.Infof("物品ID: %s, 名称: %s, 描述: %s, 市场名称: %s, 市场Hash名称: %s, 英文市场名称: %s, 英文市场Hash名称: %s, 接受者: %s, 是否可交易: %t, 是否可在市场交易: %t, 是否标准化商品: %t", item.AssetID, item.Name, item.DescriptionName, item.MarketName, item.MarketHashName, item.EnglishMarketName, item.EnglishMarketHashName, item.ReceiverName, item.Tradable, item.Marketable, item.Commodity)
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

	items, err := client.GetInventory(Constants.TeamFortress2, Constants.TeamFortress2Catetory, Dao.WithTradable(1), Dao.WithMarketable(1), nil)
	if err != nil {
		Logger.Error(err)
		return
	}

	// 随机选择一个饰品
	randomIndex := rand.Intn(len(items))

	response, err := client.PutList(Constants.TeamFortress2, Constants.TeamFortress2Catetory, items[randomIndex].AssetID, 10000, 23, string(data))
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.Infof("饰品上架成功: %+v", response)
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

	if err := client.BuyListing(440, "23", "527629225053396147", "", 30, 16, maFileContent); err != nil {
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

func TestGetPartnerInventory(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	items, err := client.GetPartnerInventory("https://steamcommunity.com/tradeoffer/new/?partner=739009475&token=k0hOOOS4", 440, 2)
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.Infof("伙伴库存物品 (%d 个)", len(items))
	for _, item := range items {
		Logger.Infof("[物品ID]: %s, [名称]: %s, [市场名称]: %s", item.ID, item.MarketName, item.MarketHashName)
	}
}

func TestSendGift(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	items, err := client.GetInventory(Constants.TeamFortress2, Constants.TeamFortress2Catetory, Dao.WithTradable(1), Dao.WithMarketable(1), Dao.WithCommodity(0))
	if err != nil {
		Logger.Error(err)
		return
	}

	if len(items) == 0 {
		Logger.Error("没有库存物品")
		return
	}

	data, err := os.ReadFile("mafiles/" + client.GetUsername() + ".maFile")
	if err != nil {
		return
	}

	maFileContent := string(data)

	Logger.Info(client.SendGift("https://steamcommunity.com/tradeoffer/new/?partner=1603469310&token=orFqPFWu", items[0].AssetID, maFileContent))
}
