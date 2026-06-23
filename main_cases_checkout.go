package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/JovanniChen/SteamDB/Steam/Logger"
	"github.com/JovanniChen/SteamDB/Steam/Model"
)

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

func TestAddItemToCart(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	addCartItems := make([][]Model.AddCartItem, 0)
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1535794, AccountidGiftee: 1535794, Message: "Apewar"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1430343, AccountidGiftee: 352956450, Message: "Apewar"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 489963, AccountidGiftee: 352956450, Message: "霓虹深渊 - 游戏原声"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 181611, AccountidGiftee: 352956450, Message: "Slay the Spire"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1011400, AccountidGiftee: 352956450, Message: "坤坤轮盘"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 96096, AccountidGiftee: 352956450, Message: "Mind Games"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 605518, AccountidGiftee: 352956450, Message: "Funny Truck"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 170869, AccountidGiftee: 352956450, Message: "Trivia Night"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 169006, AccountidGiftee: 352956450, Message: "Dead Drop"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 377271, AccountidGiftee: 352956450, Message: "TTV3"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 1298871, AccountidGiftee: 352956450, Message: "Gladiator Fights"}})
	// addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 272173, AccountidGiftee: 352956450, Message: "Bighead Runner"}})

	for _, addCartItem := range addCartItems {
		if err := client.AddItemToCartSelf(addCartItem); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("添加购物车成功")
	}
}

func TestBuyGameToSelf(accountIndex int) {
	Logger.Info("-------------------------------- 给自己购买游戏 --------------------------------")

	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	addCartItems := make([][]Model.AddCartItem, 0)
	// addCartItems = append(addCartItems, []Model.AddCartItem{{BundleID: 13013, AccountidGiftee: 352956450, Message: "怪物猎人"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 645485}}) // Barro 22

	payLinks := make([]string, 0)

	for _, addCartItem := range addCartItems {
		if err := client.AddItemToCartSelf(addCartItem); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("添加购物车成功")

		// err := client.SetStoreCountry("CN")
		// if err != nil {
		// 	Logger.Info("设置商店国家失败")
		// 	Logger.Error(err)
		// 	return
		// }
		// Logger.Info("设置商店国家成功")

		// err = client.SetCheckoutCountry("CN")
		// if err != nil {
		// 	Logger.Info("设置结算国家失败")
		// 	Logger.Error(err)
		// 	return
		// }
		// Logger.Info("设置结算国家成功")

		transID, err := client.InitTransaction("alipay", "CN", 0)
		if err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("初始化交易成功: ", transID)

		result, err := client.GetFinalPriceWithDetails(transID)
		if err != nil {
			Logger.Error(err)
			return
		}

		fmt.Printf("%+v\n", result)

		Logger.Info("获取最终价格成功: ", result.Total)

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

func TestBuyGameToOther(accountIndex int) {
	Logger.Info("-------------------------------- 给他人赠送游戏 --------------------------------")
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	addCartItems := make([][]Model.AddCartItem, 0)
	// addCartItems = append(addCartItems, []Model.AddCartItem{{BundleID: 13013, AccountidGiftee: 352956450, Message: "怪物猎人"}})
	addCartItems = append(addCartItems, []Model.AddCartItem{{PackageID: 645485, AccountidGiftee: 739009475, Message: "Barro 22"}}) // Barro 22

	payLinks := make([]string, 0)

	for _, addCartItem := range addCartItems {
		if err := client.AddItemToCart(addCartItem); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("添加购物车成功")

		transID, err := client.InitTransaction("alipay", "CN", 1)
		if err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("初始化交易成功: ", transID)

		result, err := client.GetFinalPriceWithDetails(transID)
		if err != nil {
			Logger.Error(err)
			return
		}

		fmt.Printf("%+v\n", result)

		Logger.Info("获取最终价格成功: ", result.Total)

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

func TestAddFunds(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	payLink, err := client.AddFunds(3000)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(payLink)
}

func TestAddFundsWithCountry(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	payLink, err := client.AddFundsWithCountry(4000, "HK")
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Info(payLink)
}

func TestSetCountry(accountIndex int) {
	client, err := loadFromSession(accountIndex)
	if err != nil {
		Logger.Error(err)
		return
	}

	Logger.Info(client.SetStoreCountry("HK"))
}
