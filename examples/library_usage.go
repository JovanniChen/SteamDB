// SteamDB库使用示例
// 展示如何在其他项目中引用和使用SteamDB库
package main

import (
	"fmt"
	"log"

	"github.com/steamdb/steamdb-go/Steam/Config"
	"github.com/steamdb/steamdb-go/Steam/Dao"
)

// LibraryUsageExample 库使用示例
func LibraryUsageExample() {
	// 方式1：直接创建配置（最灵活，推荐用于库引用）
	config := Config.NewConfigWithCredentials(
		"your_username",
		"your_password", 
		"", // SharedSecret可选
	)
	
	// 可选：调整网络配置
	config.RequestTimeout = 60
	config.MaxRetries = 3
	config.EnableTLSVerify = true
	
	// 创建Steam客户端
	steamClient := Dao.NewWithConfig(
		config.Proxy,
		config.EnableTLSVerify,
		config.RequestTimeout,
		config.MaxRetries,
	)
	
	// 使用客户端
	err := steamClient.Login(config.Username, config.Password, config.SharedSecret)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	
	fmt.Println("Steam登录成功")
}

// DaoLayerUsageExample DAO层使用示例
func DaoLayerUsageExample() {
	// 创建配置
	config := Config.NewConfig()
	
	// 从调用方自己的配置系统设置凭据
	// 例如：从数据库、命令行参数、自定义配置文件等
	config.Username = getUsernameFromYourSystem()
	config.Password = getPasswordFromYourSystem()
	
	// 创建DAO实例
	dao := Dao.NewWithConfig(
		config.Proxy,
		config.EnableTLSVerify,
		config.RequestTimeout,
		config.MaxRetries,
	)
	
	// 使用DAO层API
	err := dao.Login(config.Username, config.Password, config.SharedSecret)
	if err != nil {
		log.Fatalf("DAO层登录失败: %v", err)
	}
	
	fmt.Println("DAO层登录成功")
}

// EnvironmentBasedUsageExample 基于环境变量的使用示例  
func EnvironmentBasedUsageExample() {
	// 创建默认配置
	config := Config.NewConfig()
	
	// 让调用方决定是否从环境变量加载
	config.LoadFromEnv()
	
	// 验证配置
	if err := config.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}
	
	// 使用配置...
	fmt.Println("基于环境变量的配置完成")
}

// 模拟调用方自己的配置系统
func getUsernameFromYourSystem() string {
	// 这里可以是从数据库、配置服务、命令行等任何地方获取
	return "your_steam_username"
}

func getPasswordFromYourSystem() string {
	// 这里可以是从安全存储、密钥管理服务等获取
	return "your_steam_password"
}

func MainExample() {
	fmt.Println("SteamDB库使用示例")
	fmt.Println("=================")
	
	fmt.Println("\n1. 基本库使用示例:")
	LibraryUsageExample()
	
	fmt.Println("\n2. DAO层使用示例:")
	DaoLayerUsageExample()
	
	fmt.Println("\n3. 环境变量配置示例:")
	EnvironmentBasedUsageExample()
}

// 避免main函数冲突的初始化函数
func init() {
	// 这是示例代码，实际使用时请调用MainExample()
}