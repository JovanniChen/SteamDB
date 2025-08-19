// SteamDB库使用示例
// 展示如何在其他项目中引用和使用SteamDB库
package main

import (
	"fmt"
	"log"

	"github.com/JovanniChen/SteamDB/Steam"
	"github.com/JovanniChen/SteamDB/Steam/Dao"
)

// LibraryUsageExample 库使用示例
func LibraryUsageExample() {
	// 创建Steam客户端
	client, err := Steam.NewClient(Steam.DefaultConfig())
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	
	// 登录凭据
	credentials := &Steam.LoginCredentials{
		Username:     "your_username",
		Password:     "your_password", 
		SharedSecret: "", // Steam Guard共享密钥，可选
	}
	
	// 执行登录
	userInfo, err := client.Login(credentials)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	
	fmt.Printf("Steam登录成功！用户: %s, Steam ID: %d\n", userInfo.Username, userInfo.SteamID)
}

// DaoLayerUsageExample DAO层使用示例
func DaoLayerUsageExample() {
	// 创建DAO实例，参数为代理地址，空字符串表示不使用代理
	dao := Dao.New("")
	
	// 执行登录
	err := dao.Login("your_username", "your_password", "")
	if err != nil {
		log.Fatalf("DAO层登录失败: %v", err)
	}
	
	fmt.Println("DAO层登录成功")
}

// AdvancedUsageExample 高级使用示例
func AdvancedUsageExample() {
	// 创建带代理的客户端配置
	config := &Steam.Config{
		Proxy:   "127.0.0.1:8080", // 可选：使用代理
		Timeout: 60,               // 60秒超时
	}
	
	client, err := Steam.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	
	// 使用Steam Guard
	credentials := &Steam.LoginCredentials{
		Username:     "your_username",
		Password:     "your_password",
		SharedSecret: "your_shared_secret", // base64编码的共享密钥
	}
	
	userInfo, err := client.Login(credentials)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	
	// 获取用户积分信息
	summary, err := client.GetPointsSummary(userInfo.SteamID)
	if err != nil {
		log.Printf("获取积分信息失败: %v", err)
	} else {
		fmt.Printf("用户积分: %d, 等级: %d\n", summary.Points, summary.Level)
	}
	
	fmt.Println("高级功能使用完成")
}

func MainExample() {
	fmt.Println("SteamDB库使用示例")
	fmt.Println("=================")
	
	fmt.Println("\n1. 基本库使用示例:")
	LibraryUsageExample()
	
	fmt.Println("\n2. DAO层使用示例:")
	DaoLayerUsageExample()
	
	fmt.Println("\n3. 高级功能示例:")
	AdvancedUsageExample()
}

// 避免main函数冲突的初始化函数
func init() {
	// 这是示例代码，实际使用时请调用MainExample()
}