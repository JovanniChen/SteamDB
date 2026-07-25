package main

import (
	"fmt"

	"github.com/JovanniChen/SteamDB/Steam"
	"github.com/JovanniChen/SteamDB/Steam/Dao"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
	"github.com/JovanniChen/SteamDB/Steam/Model"
)

func TestGetGameUpdateInofs(gameID int) {
	client, err := Steam.NewClient(config)
	if err != nil {
		Logger.Error(err)
		return
	}

	// 使用简化方法，直接获取提取的更新事件
	updateEvents, totalFound, needUpdate, err := client.GetGameUpdateEvents(gameID, 1)
	if err != nil {
		Logger.Error(err)
		return
	}

	// 简洁输出
	fmt.Printf("游戏ID: %d | 找到: %d条 | 提取: %d条 | 需要更新: %v\n",
		gameID, totalFound, len(updateEvents), needUpdate)

	if len(updateEvents) == 0 {
		fmt.Println("  ⚠️  未找到任何更新事件")
		return
	}

	// 只显示最新的一条事件
	if len(updateEvents) > 0 {
		event := updateEvents[0]
		fmt.Printf("  最新事件: %s (ID: %s)(Time: %d)\n", event.EventName, event.UniqueID, event.StartTime)
	}
}

func TestGetSummary(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	summary, err := client.GetPointsSummary(client.GetSteamID())
	Logger.Info("GetSummary -> ", summary)

	// 提取登录Cookies
	loginCookies := make(map[string]*Dao.LoginCookie)
	if cookies := client.GetLoginCookies(); cookies != nil {
		loginCookies = cookies
	}

	if loginCookies["checkout.steampowered.com"] != nil {
		fmt.Println(loginCookies["checkout.steampowered.com"])
	}
	fmt.Println(loginCookies["checkout.steampowered.com"])

	if loginCookies["steamcommunity.com"] != nil {
		fmt.Println(loginCookies["steamcommunity.com"])
	}

	if loginCookies["store.steampowered.com"] != nil {
		fmt.Println(loginCookies["store.steampowered.com"])
	}
}

func TestGetProductByAppUrl(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	// Logger.Info(client.GetProductByAppUrl("https://store.steampowered.com/app/2479810"))
	Logger.Info(client.GetProductByAppUrl("https://store.steampowered.com/app/1119730"))

}

func TestGetPackageDetails(accountIndex int) {
	client, err := Steam.NewClient(config)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetPackageDetails(3513600))
}

func TestGetStorePurchaseHistory(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	result, err := client.GetStorePurchaseHistory()
	if err != nil {
		Logger.Error(err)
		return
	}

	fmt.Printf("消费记录总数: %d\n", result.TotalRecords)
	fmt.Printf("最新连续未退款礼物购买记录数量: %d\n", len(result.LatestUnrefundedGiftPurchases))
	for _, record := range result.LatestUnrefundedGiftPurchases {
		printStorePurchaseHistoryRecord(record)
	}
}

func printStorePurchaseHistoryRecord(record Model.StorePurchaseHistoryRecord) {
	fmt.Printf("---------- 消费记录 %d ----------\n", record.Index)
	fmt.Printf("交易ID: %s\n", record.TransactionID)
	fmt.Printf("日期: %s\n", record.Date)
	fmt.Printf("物品: %s\n", record.Item)
	if record.Receiver != "" {
		fmt.Printf("收礼人: %s\n", record.Receiver)
	}
	fmt.Printf("类型: %s\n", record.TransactionType)
	fmt.Printf("支付方式: %s\n", record.Payment)
	fmt.Printf("基础价格: %s\n", record.BasePrice)
	fmt.Printf("税额: %s\n", record.Tax)
	fmt.Printf("运费: %s\n", record.Shipping)
	fmt.Printf("总计: %s\n", record.Total)
	fmt.Printf("钱包变更: %s\n", record.WalletChange)
	fmt.Printf("钱包余额: %s\n", record.WalletBalance)
	fmt.Printf("已退款: %t\n", record.Refunded)
}
