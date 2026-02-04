// Steam数据库操作主程序
// 本程序用于连接Steam平台，进行用户登录、获取令牌代码和添加反应等操作
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/JovanniChen/SteamDB/Steam"
	"github.com/JovanniChen/SteamDB/Steam/Constants"
	"github.com/JovanniChen/SteamDB/Steam/Dao"
	"github.com/JovanniChen/SteamDB/Steam/Logger"
	"github.com/JovanniChen/SteamDB/Steam/Model"
)

//  {Username: "za0ww9ml4xl2", Password: "HLHxGyRMm6Zi", SharedSecret: "F54xOr9Tpyd5fAxgKx+RHR7vHik="}, // [0] [xv6753] [46]
// 	{Username: "zytmnd2097", Password: "awtekBcEkXz9", SharedSecret: "vNVDHuqBle/rnsG7EQW2xQUqlME="},   // [1] [4wzwg]  [45]
// 	{Username: "zwrvsq6897", Password: "5uoIBclSSBI8", SharedSecret: "kUcQLn0pJutKt9oeh8yRDG7t+o8="},   // [2] [wqrmhz] [44]
// 	{Username: "zuzuaw8238", Password: "uYj035ynLA5N", SharedSecret: "yKuRsv/OmI584XxMt2LUWWbCM+Y="},   // [3] [kxweoq] [40]
// 	{Username: "yrknu899", Password: "FyoR1QV8brUd", SharedSecret: "q8JcjcE5jc65C7YntMrME8HJ9sY="},     // [4] [3zgmh7] [47]
// 	{Username: "mbkle379", Password: "CFs91IvocA39", SharedSecret: "sIF2wljQzxzya9xVO/VtEs1pUwc="},     // [5] [x5x3g8] [48]
// 	{Username: "lvqxpe8572", Password: "5gfweOafGM3S", SharedSecret: "QLWiEAN8ebHLkGtt7HHtuZyMwDg="},   // [6] [x5x3g8] [49]
// 	{Username: "naotqp7801", Password: "ja9C5LZelku0", SharedSecret: "g+kIH7JuL98R5O00j87379CkFus="},   // [7] [x5x3g8] [49]
// 	{Username: "iatfqv6444", Password: "NOJsp0b1aqbj", SharedSecret: "wCdOSNrhPjXrJEpg3FX643+fseQ="},   // [8] [x5x3g8] [49]
// 	{Username: "uwxhfw8800", Password: "ybuYe33Qg2Dr", SharedSecret: "ViUburoMwWe88QJfL5f0KPPoY68="},   // [9] [x5x3g8] [49]
// 	{Username: "xqkea03549", Password: "wuwQJ5WFdZp1", SharedSecret: "59z0KMWJFdgfWrSgYYADD/LBPyU="},   // [10] [6ck2bcax] [53]

var accounts = []Account{
	{Username: "zszvlv6362", Password: "ejuj7Rnof1BB", SharedSecret: "mQI147JxRz78GWjDdQEBoL7aaBc="},   // [0] [45]
	{Username: "gbmqnl7210", Password: "i80sMCigz1rw", SharedSecret: "Uinb4sxNpcP8KQBcYgdAZ2eiJDg="},   // [1] [46]
	{Username: "zdckla1506", Password: "d3c9InY7Epwi", SharedSecret: "rQJ4b42FyGsvGcp6XYx+SEYylyo="},   // [2] [47]
	{Username: "uvtjrm4501", Password: "u9NIlsVugLH5", SharedSecret: "y77Jk5v4rxrck/149zDMB+b3s/U="},   // [3] [48]
	{Username: "ddndd12412", Password: "New0KJYVv16", SharedSecret: "VoSY5VrnD+CJooEVrlADofTGTok="},    // [4] [51]
	{Username: "ttmsq72777", Password: "yoRD7x6LQvgu", SharedSecret: "5boHTiGFhQoszGcpFDLB7H7thng="},   // [5] [52]
	{Username: "xqkea03549", Password: "wuwQJ5WFdZp1", SharedSecret: "59z0KMWJFdgfWrSgYYADD/LBPyU="},   // [6] [53]
	{Username: "ffotd74229", Password: "oP4M4CMHAftX", SharedSecret: "IDhBX3NM+8fZCti4C3d6oFhXI6E="},   // [7] [54]
	{Username: "j47hoord6j", Password: "NewRP7IhC9Z", SharedSecret: "Gwgztog4anK0soQp4IgLaZIki0s="},    // [8] [57] 市场不可用
	{Username: "naotqp7801", Password: "ja9C5LZelku0", SharedSecret: "g+kIH7JuL98R5O00j87379CkFus="},   // [9]
	{Username: "zszvlv6362", Password: "ejuj7Rnof1BB", SharedSecret: "mQI147JxRz78GWjDdQEBoL7aaBc="},   // [10]
	{Username: "fbrdz08225", Password: "NewNWnME1R6", SharedSecret: "VjYAPygKL4jxwSu69HeyzW58r3M="},    // [11]
	{Username: "rwfio67235", Password: "JzBvNCICYfFx", SharedSecret: "0C4hU7ieyVyYFvdDPKoTII20xMc="},   // [12]
	{Username: "ejvp732231", Password: "myz2bzwCzFYQ", SharedSecret: "KHzBIonDKW8enmoCUYgLN+oYQ4M="},   // [13]
	{Username: "mcg9ipxd04nl", Password: "AAlLQXPdDy3U", SharedSecret: "zUN+RyvAQZjHnyT+5guHBPB2NOg="}, // [14]
	{Username: "pfze6stttee", Password: "P4A9ydYvzGmq", SharedSecret: "qNSLkP8OjsD/VuHG5eFGUjSupCs="},  // [15]
	{Username: "krqk10ik5qk", Password: "vPfPgdUqVX76", SharedSecret: "1daUVnJtxNg4pB2gxV2l10jfz1U="},  // [16]
	{Username: "bhrcrnulng5", Password: "YF8TX9fAxWsq", SharedSecret: "X1z0h/KTJmns1um4ThZdRGrrNps="},  // [17]
	{Username: "dih9u8nad", Password: "BJdENppgNHxH", SharedSecret: "TFuheV7W4oPoH4Q2EH8EEi9vmKU="},    // [18]
	{Username: "tvyij7pxdasz", Password: "0DLuZvp5MSEI", SharedSecret: "YGnTkbpo/uOGFtNeFRhGlIQxrEg="}, // [19]
	{Username: "hnysh898sg", Password: "mC43y8o8irxT", SharedSecret: "ynrDLQop7KGLFe0DfiFcW8lOy6A="},   // [20]
	{Username: "tzjn5e5xnz5y", Password: "wBvFtVZkCpsZ", SharedSecret: "RGcG25ZSJiswJpwz56DTaPWR+nI="}, // [21]
	{Username: "tgf21ra7e3", Password: "50Z2DJJMNnfI", SharedSecret: "FWhvNXuMuhPGj4V1FLBD31Fzks8="},   // [22]
	{Username: "acjjzz1twx", Password: "H3s5wvDLgL4d", SharedSecret: "+xtZKqLMsIMu7T4LgP3rO6wNV2Q="},   // [23]
	{Username: "akji2nh66u0f", Password: "t88j8RZs37BY", SharedSecret: "KjTGnop2fyCu2E7hBWflE2TdEO4="}, // [24]
	{Username: "kneao1dmw", Password: "5Dtu8Kk9JIAg", SharedSecret: "GjkmsGH1utG1QVyg06NQc1wNTwg="},    // [25]
	{Username: "plzhhgt075", Password: "ESvgOZnLTjKb", SharedSecret: "WxqsqlawbHEB7Pjah0wutOpLKE0="},   // [26]
	{Username: "iuuwhmusxdv7", Password: "STQI7NOal7l6", SharedSecret: "+xtZKqLMsIMu7T4LgP3rO6wNV2Q="}, // [27]
	{Username: "bmlgbjot5hz", Password: "1Z5pkOTuZigf", SharedSecret: "GZKHLVfxwjYPFMBF33l7Vu3NMY4="},  // [28]
	{Username: "aeuybz0905", Password: "f4J5Cs6cHnHP", SharedSecret: "UGDYQfigAc47yH/wPcL0E3PCHPY="},   // [29]
}

var config *Steam.Config = Steam.NewConfig("your_username:your_password@8.217.238.29:8080")

// var config *Steam.Config = Steam.NewConfig("")
// var config *Steam.Config = Steam.DefaultConfig()

func main() {
	accountIndex := 6
	// TestLogin(accountIndex)
	// TestGetTokenCode(accountIndex)
	// TestGetFriendInfoByLink(accountIndex)
	// TestGetFriendInfoByLinkAndAddFriend(accountIndex)
	// TestGetProductByAppUrl(accountIndex)
	// TestGetSteamGift(accountIndex)
	// TestUnsendAllGift(accountIndex)
	// TestConcurrentPayment(accountIndex)
	// TestTransactionStatus(accountIndex)
	// TestUnsendGift(accountIndex)
	// TestGetTokenCode(accountIndex)
	// TestSetLanguage(accountIndex)
	// TestClearCart(accountIndex)
	// TestGetCart(accountIndex)
	// TestAddItemToCart(accountIndex)
	// TestInitTransaction(accountIndex)
	// TestAddItemToCartAndInitTransaction(accountIndex)
	// TestAccess(accountIndex)
	// TestValidateCart(accountIndex)
	// TestCancelTransaction(accountIndex)
	// TestGetFinalPrice(accountIndex)
	// TestTestGetPayLinkAgain(accountIndex)
	// TestAddFriendByFriendCode(accountIndex)
	// TestAddFriendByLink(accountIndex)
	// TestCheckIsFriend(accountIndex)
	// TestCheckFriendStatus(accountIndex)
	// TestRemoveFriend(accountIndex)
	// TestCheckAccountAvailable(accountIndex)
	// TestGetSummary(accountIndex)
	TestGetInventory(accountIndex)
	// TestGetMyListings(accountIndex)
	// TestGetMyListings(accountIndex)
	// TestGetConfirmations(accountIndex)
	// TestGetConfirmations(accountIndex)
	// TestPutList(accountIndex)
	// TestBuyListing(accountIndex)
	// TestPutList2(accountIndex)
	// TestGetConfirmations(accountIndex)
	// TestRemoveMyListings(accountIndex)
	// TestGetBalance(accountIndex)
	// TestGetWaitBalance(accountIndex)
	// TestGetInventoryAndPutList(accountIndex)
	// TestCreateOrder(accountIndex)
	// TestGetGameUpdateInofs(1879330)
}

func TestTestGetPayLinkAgain(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	client.SetProxy("8.217.238.29:8080")
	Logger.Info(client.TestGetPayLinkAgain())
}

func TestValidateCart(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.ValidateCart())
}

func TestTransactionStatus(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.TransactionStatus("483707716887528185", 1))
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

func TestAccess(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	checkoutURL, err := client.AccessCheckoutURL("56990384110504959")
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info("获取支付页面成功: ", checkoutURL)
}

func TestGetFinalPrice(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetFinalPrice("56990384110504959"))
}

func TestCancelTransaction(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	if err := client.CancelTransaction("113286647533192916"); err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info("取消交易成功")
}

func TestInitTransaction(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.InitTransaction())
}

func TestAddItemToCart(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	addCartItems := make([][]Model.AddCartItem, 0)
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 827941, AccountidGiftee: 352956450, Message: "Apewar"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 489963, AccountidGiftee: 352956450, Message: "霓虹深渊 - 游戏原声"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 181611, AccountidGiftee: 352956450, Message: "Slay the Spire"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1011400, AccountidGiftee: 352956450, Message: "坤坤轮盘"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 96096, AccountidGiftee: 352956450, Message: "Mind Games"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 605518, AccountidGiftee: 352956450, Message: "Funny Truck"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 170869, AccountidGiftee: 352956450, Message: "Trivia Night"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 169006, AccountidGiftee: 352956450, Message: "Dead Drop"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 377271, AccountidGiftee: 352956450, Message: "TTV3"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1298871, AccountidGiftee: 352956450, Message: "Gladiator Fights"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 272173, AccountidGiftee: 352956450, Message: "Bighead Runner"}})

	for _, addCartItem := range addCartItems {
		if err := client.AddItemToCart(addCartItem); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("添加购物车成功")
	}
}

func TestAddItemToCartAndInitTransaction(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	addCartItems := make([][]Model.AddCartItem, 0)
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 827941, AccountidGiftee: 352956450, Message: "Apewar"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 489963, AccountidGiftee: 352956450, Message: "霓虹深渊 - 游戏原声"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 181611, AccountidGiftee: 352956450, Message: "Slay the Spire"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1011400, AccountidGiftee: 352956450, Message: "坤坤轮盘"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 96096, AccountidGiftee: 352956450, Message: "Mind Games"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 605518, AccountidGiftee: 352956450, Message: "Funny Truck"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 170869, AccountidGiftee: 352956450, Message: "Trivia Night"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 169006, AccountidGiftee: 352956450, Message: "Dead Drop"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 377271, AccountidGiftee: 352956450, Message: "TTV3"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1298871, AccountidGiftee: 352956450, Message: "Gladiator Fights"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 272173, AccountidGiftee: 352956450, Message: "Bighead Runner"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 203389, AccountidGiftee: 352956450, Message: "Thief Simulator"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 8183, AccountidGiftee: 352956450, Message: "Terraria"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 379292, AccountidGiftee: 352956450, Message: "Ranch Simulator"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 121418, AccountidGiftee: 352956450, Message: "Don't Starve Together"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 240331, AccountidGiftee: 352956450, Message: "Swarmlake"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 85819, AccountidGiftee: 352956450, Message: "Plantera"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 195824, AccountidGiftee: 352956450, Message: "Polygoneer"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1150347, AccountidGiftee: 1753831820, Message: "Microsoft Flight Simulator 2024"}, {PackageID: 625300, AccountidGiftee: 1753831820, Message: "BLUE REFLECTION"}})

	payLinks := make([]string, 0)

	for _, addCartItem := range addCartItems {
		if err := client.AddItemToCart(addCartItem); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("添加购物车成功")

		transID, err := client.InitTransaction()
		if err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("初始化交易成功: ", transID)

		total, err := client.GetFinalPrice(transID)
		if err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("获取最终价格成功: ", total)

		checkoutURL, err := client.AccessCheckoutURL(transID)
		if err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("获取支付页面成功: ", checkoutURL)

		payLinks = append(payLinks, checkoutURL)

		if err := client.CancelTransaction(transID); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("取消交易成功")

		if err := client.ClearCart(); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("清空购物车成功")
	}
	for _, payLink := range payLinks {
		Logger.Info("支付链接: ", payLink)
	}
}

func TestConcurrentPayment(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	addCartItems := make([][]Model.AddCartItem, 0)
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 73106, AccountidGiftee: 352956450, Message: "超级鸡马"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 198332, AccountidGiftee: 352956450, Message: "Arcadian Atlas"}})

	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1107210, AccountidGiftee: 352956450, Message: "球跳塔"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1101478, AccountidGiftee: 352956450, Message: "恐怖之眼"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 984094, AccountidGiftee: 352956450, Message: "纸片大作战2"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 110604, AccountidGiftee: 352956450, Message: "Antisphere"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1033686, AccountidGiftee: 352956450, Message: "咪子不要! - 金缮之美"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1033677, AccountidGiftee: 352956450, Message: "咪子不要! - 日常小物"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 675628, AccountidGiftee: 352956450, Message: "Risen Soundtrack"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 279451, AccountidGiftee: 352956450, Message: "A Sky Full of Stars - Original Sound Track"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 984063, AccountidGiftee: 352956450, Message: "黑洞大作战"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 979135, AccountidGiftee: 352956450, Message: "炮弹人冲冲冲"}})

	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 445776, AccountidGiftee: 352956450, Message: "BIOMUTANT - Soundtrack"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 984097, AccountidGiftee: 352956450, Message: "神枪手强尼"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 834807, AccountidGiftee: 352956450, Message: "奔跑吧，香肠！"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 616514, AccountidGiftee: 352956450, Message: "Farm Kitten - Puzzle Pipes"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1013759, AccountidGiftee: 352956450, Message: "来切我鸭"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 197731, AccountidGiftee: 352956450, Message: "1bitHeart Original Soundtrack"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 366861, AccountidGiftee: 352956450, Message: "Sudoku 9x16x25"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 990639, AccountidGiftee: 352956450, Message: "毒液入侵者"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1383441, AccountidGiftee: 352956450, Message: "Merlin Survivors"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 375305, AccountidGiftee: 352956450, Message: "Kakuro"}})

	transIDs := make([]string, 0)

	for i, addCartItem := range addCartItems {

		if err := client.AddItemToCart(addCartItem); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("添加购物车成功: ", i+1, " / ", len(addCartItems))

		transID, err := client.InitConcurrentTransaction()
		if err != nil {
			Logger.Error(err)
			Logger.Error("初始化交易失败: ", i+1, " / ", len(addCartItems))
			return
		}
		Logger.Info("初始化交易成功: ", transID, " ", i+1, " / ", len(addCartItems))

		total, err := client.GetFinalPrice(transID)
		if err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("获取最终价格成功: ", total, " ", i+1, " / ", len(addCartItems))

		if err := client.ClearCart(); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("清空购物车成功: ", i+1, " / ", len(addCartItems))

		if err := client.CancelTransaction("231481419234113013"); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("取消交易成功: ", i+1, " / ", len(addCartItems))

		transIDs = append(transIDs, transID)
	}

	totalDuration := 10 * time.Second
	startTime := time.Now()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		elapsed := time.Since(startTime)
		remaining := totalDuration - elapsed

		if remaining <= 0 {
			Logger.Info(fmt.Sprintf(
				"休眠完成！总耗时: %v",
				elapsed.Round(time.Second),
			))
			break
		}

		Logger.Info(fmt.Sprintf(
			"已休眠: %v, 还需休眠: %v",
			elapsed.Round(time.Second),
			remaining.Round(time.Second),
		))

		<-ticker.C
	}

	wg := sync.WaitGroup{}
	wg.Add(len(transIDs))
	for _, transID := range transIDs {
		go func(transID string) {
			if err := client.FinalizeTransaction(transID); err != nil {
				Logger.Error(err)
				return
			}
			wg.Done()
			Logger.Info("完成最终支付: ", transID)
		}(transID)
	}
	wg.Wait()
	Logger.Info("同时付交易完成")
}

func TestGetCart(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetCart())
}

func TestClearCart(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	if err := client.ClearCart(); err != nil {
		Logger.Error("清空购物车失败: ", err)
		return
	}
	// 你TestClearCart

	Logger.Info("清空购物车成功")
}

func TestRemoveFriend(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	if err := client.RemoveFriend(76561198313222178); err != nil {
		Logger.Error("删除好友失败: ", err)
		return
	}
	Logger.Info("删除好友成功")
}

func TestAddFriendByFriendCode(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	if err := client.AddFriendByFriendCode(352956450); err != nil {
		Logger.Error("添加好友失败: ", err)
		return
	}
	Logger.Info("添加好友成功")
}

func TestCheckIsFriend(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	isFriend, err := client.CheckIsFriend("76561199794056204")
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info("是否是好友: ", isFriend)
}

func TestCheckFriendStatus(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	if err = client.CheckFriendStatus("https://s.team/p/jtgt-jrbr/CDPKKTCF"); err != nil {
		Logger.Error("检查好友状态失败: ", err)
		return
	}
	Logger.Info("检查好友状态成功")
}

func TestAddFriendByLink(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	steamID, err := client.AddFriendByLink("https://s.team/p/chbn-qbdd/RCVBGMHJ")
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.Info("添加好友成功: ", steamID)
}

func TestGetFriendInfoByLink(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	friendInfo, inviteToken, err := client.GetFriendInfoByLink("https://s.team/p/chbn-qbdd/BHMGGQBR")
	if err != nil {
		Logger.Info(friendInfo)
		Logger.Error(err)
		return
	}
	Logger.Info(friendInfo)
	Logger.Info(inviteToken)
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

func TestLogin(accountIndex int) {
	if _, err := os.Stat("temp/session_" + strconv.Itoa(accountIndex) + ".json"); os.IsNotExist(err) {
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

func TestGetMyListings(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	client.GetMyListings()
	// Logger.Infof("已上架物品 (%d 个) -> %+v\n", len(activeListings), activeListings)
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

	// items, err := client.GetInventory(Constants.PrimalCarnage, Constants.PrimalCarnageCategory)
	// if err != nil {
	// 	Logger.Error(err)
	// 	return
	// }

	// if len(items) == 0 {
	// 	Logger.Error("无可用库存")
	// 	return
	// }

	// // 随机
	// randomIndex := rand.Intn(len(items))
	// randomItem := items[randomIndex]

	data, err := os.ReadFile("mafiles/" + client.GetUsername() + ".maFile")
	if err != nil {
		Logger.Error(err)
		return
	}

	_, err = client.PutList(Constants.PrimalCarnage, Constants.PrimalCarnageCategory, "123123", 0.14, 23, string(data))
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.Info("上架成功")
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
	Logger.Info(client.BuyListing("321360", "9079938361156157936", "", 0.16, 0.14, maFileContent).Error())
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

	Logger.Info(client.CreateOrder("Giftapult", 0.12, 15, maFileContent))
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

func TestGetProductByAppUrl(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(client.GetProductByAppUrl("https://store.steampowered.com/app/2669320"))
}

func TestGetFriendInfoByLinkAndAddFriend(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}
	friendInfo, inviteToken, err := client.GetFriendInfoByLink("https://s.team/p/chbn-qbdd/QVRQMRJK")
	if err != nil {
		Logger.Error(err)
		return
	}

	_, err = client.AddFriendByInviteTokenAndSteamID(inviteToken, friendInfo.AbuseID)
	if err != nil {
		Logger.Error("通过邀请token和steamID添加好友失败,错误:", err)
		return
	}

	Logger.Info("通过邀请token和steamID添加好友成功")
}

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
