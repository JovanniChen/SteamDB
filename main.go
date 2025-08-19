// Steam数据库操作主程序
// 本程序用于连接Steam平台，进行用户登录、获取令牌代码和添加反应等操作
package main

import (
	"fmt"

	"github.com/steamdb/steamdb-go/Steam/Dao"
)

// main 主函数，程序入口点
// 执行Steam平台相关操作的演示流程
func main() {
	// 创建Dao实例，参数为代理地址，空字符串表示不使用代理
	// proxy := "your_username:your_password@52.53.230.117:8080"
	// proxy := "your_username:your_password@13.57.219.39:8080"
	// proxy := "your_username:your_password@3.101.191.201:8080"
	// proxy := "your_username:your_password@52.53.183.38:8080"
	d := Dao.New("")

	// 获取令牌代码，用于双因素认证
	// 参数是base64编码的身份标识符
	code, _ := d.GetTokenCode("yKuRsv/OmI584XxMt2LUWWbCM+Y=")
	fmt.Println(code)

	// 用户登录Steam平台
	// 参数：用户名，密码，令牌代码
	// err := d.Login("rgckq82191", "vxlu26493E", "")
	err := d.Login("za0ww9ml4xl2", "HLHxGyRMm6Zi", "F54xOr9Tpyd5fAxgKx+RHR7vHik=") // xv6753
	// err := d.Login("zytmnd2097", "awtekBcEkXz9", "vNVDHuqBle/rnsG7EQW2xQUqlME=") // 4wzwg
	// err := d.Login("zwrvsq6897", "5uoIBclSSBI8", "kUcQLn0pJutKt9oeh8yRDG7t+o8=") // wqrmhz
	// err := d.Login("zuzuaw8238", "uYj035ynLA5N", "yKuRsv/OmI584XxMt2LUWWbCM+Y=") // kxweoq
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取用户Cookie信息（已注释）
	// str, err := d.GetUserCookies()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(str))

	// 获取用户信息（已注释）
	// err = d.UserInfo()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// 设置语言为简体中文（已注释）
	//err = d.SetLanguage("schinese")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 以下是一些其他功能的示例调用（已注释）
	// result, err := d.GetReacionts(76561198313222178, 3)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(len(result.Reactionids))
	// fmt.Println(result) // 获取反应
	// fmt.Println(d.GetSteamID())
	fmt.Println(d.GetSummary(d.GetSteamID())) // 获取摘要
	// fmt.Println(d.GetReactionConfig()) // 获取反应配置
	// fmt.Println(d.GetSteamIDByFriendLink("https://steamcommunity.com/user/chbn-qbdd/KCDRCPRT/"))

	// 为指定用户添加反应
	// 参数：用户SteamID，反应类型，反应ID
	// fmt.Println(d.AddReaction(76561198313222178, 3, 23))
	// fmt.Println(d.CheckLogin("https://steamcommunity.com/"))
}
