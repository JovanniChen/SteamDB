// SteamDB Go客户端库使用示例
// 演示如何使用重构后的SteamDB库进行Steam平台交互
package main

import (
	"fmt"
	"log"

	"github.com/JovanniChen/SteamDB/Steam"
)

func main() {
	// 创建客户端配置
	config := Steam.DefaultConfig()
	// 如果需要使用代理，可以设置：
	// config.Proxy = "127.0.0.1:8080"
	
	// 创建Steam客户端
	client, err := Steam.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 准备登录凭据
	credentials := &Steam.LoginCredentials{
		Username:     "your_username",      // 替换为你的Steam用户名
		Password:     "your_password",      // 替换为你的Steam密码
		SharedSecret: "your_shared_secret", // 替换为你的Steam Guard共享密钥(base64编码)，如果没有2FA可以留空
	}

	// 执行登录
	fmt.Println("正在登录Steam...")
	userInfo, err := client.Login(credentials)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}

	// 显示登录成功信息
	fmt.Printf("登录成功！\n")
	fmt.Printf("Steam ID: %d\n", userInfo.SteamID)
	fmt.Printf("用户名: %s\n", userInfo.Username)
	fmt.Printf("昵称: %s\n", userInfo.Nickname)
	fmt.Printf("国家代码: %s\n", userInfo.CountryCode)

	// 获取Steam Guard令牌代码（如果配置了共享密钥）
	if credentials.SharedSecret != "" {
		code, err := client.GetTokenCode(credentials.SharedSecret)
		if err != nil {
			fmt.Printf("获取令牌代码失败: %v\n", err)
		} else {
			fmt.Printf("当前Steam Guard代码: %s\n", code)
		}
	}

	// 获取积分摘要
	fmt.Println("\n正在获取积分摘要...")
	summary, err := client.GetPointsSummary(userInfo.SteamID)
	if err != nil {
		fmt.Printf("获取积分摘要失败: %v\n", err)
	} else {
		fmt.Printf("积分余额: %d\n", summary.Points)
		fmt.Printf("用户等级: %d\n", summary.Level)
	}

	// 获取反应配置
	fmt.Println("\n正在获取反应配置...")
	reactions, err := client.GetReactionConfig()
	if err != nil {
		fmt.Printf("获取反应配置失败: %v\n", err)
	} else {
		fmt.Printf("可用反应类型数量: %d\n", len(reactions))
		for i, reaction := range reactions {
			if i < 3 { // 只显示前3个作为示例
				fmt.Printf("反应ID: %d, 消耗积分: %d\n", reaction.ReactionID, reaction.PointsCost)
			}
		}
	}

	fmt.Println("\n示例程序执行完成！")
}