// 高级使用示例：Steam积分管理工具
// 演示如何使用SteamDB库管理Steam积分和反应系统
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JovanniChen/SteamDB/Steam"
)

func main() {
	fmt.Println("=== Steam积分管理工具 ===\n")

	// 从环境变量或用户输入获取登录信息
	username := getInput("请输入Steam用户名: ")
	password := getInput("请输入Steam密码: ")
	sharedSecret := getInput("请输入Steam Guard共享密钥(可选，直接回车跳过): ")

	// 创建客户端配置
	config := Steam.DefaultConfig()
	config.Timeout = 60 * time.Second // 增加超时时间

	// 创建客户端
	client, err := Steam.NewClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 登录
	fmt.Println("\n正在登录Steam...")
	credentials := &Steam.LoginCredentials{
		Username:     username,
		Password:     password,
		SharedSecret: sharedSecret,
	}

	userInfo, err := client.Login(credentials)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}

	fmt.Printf("登录成功！欢迎 %s (Steam ID: %d)\n\n", userInfo.Nickname, userInfo.SteamID)

	// 主菜单循环
	for {
		showMenu()
		choice := getInput("请选择操作: ")

		switch choice {
		case "1":
			showUserInfo(client, userInfo)
		case "2":
			showPointsSummary(client, userInfo.SteamID)
		case "3":
			showReactionConfig(client)
		case "4":
			addReactionToUser(client)
		case "5":
			generateTokenCode(client, sharedSecret)
		case "6":
			showUserReactions(client, userInfo.SteamID)
		case "0":
			fmt.Println("感谢使用，再见！")
			return
		default:
			fmt.Println("无效选择，请重试。")
		}
		fmt.Println()
	}
}

func showMenu() {
	fmt.Println("=== 主菜单 ===")
	fmt.Println("1. 查看用户信息")
	fmt.Println("2. 查看积分摘要")
	fmt.Println("3. 查看反应配置")
	fmt.Println("4. 为用户添加反应")
	fmt.Println("5. 生成Steam Guard代码")
	fmt.Println("6. 查看用户反应记录")
	fmt.Println("0. 退出")
	fmt.Print("选择: ")
}

func showUserInfo(client *Steam.Client, userInfo *Steam.UserInfo) {
	fmt.Println("=== 用户信息 ===")
	fmt.Printf("Steam ID: %d\n", userInfo.SteamID)
	fmt.Printf("用户名: %s\n", userInfo.Username)
	fmt.Printf("昵称: %s\n", userInfo.Nickname)
	fmt.Printf("国家代码: %s\n", userInfo.CountryCode)

	// 获取访问令牌状态
	if token, err := client.GetAccessToken(); err == nil && token != "" {
		fmt.Printf("访问令牌状态: 有效 (长度: %d)\n", len(token))
	} else {
		fmt.Println("访问令牌状态: 无效或为空")
	}
}

func showPointsSummary(client *Steam.Client, steamID uint64) {
	fmt.Println("=== 积分摘要 ===")
	fmt.Println("正在获取积分信息...")

	summary, err := client.GetPointsSummary(steamID)
	if err != nil {
		fmt.Printf("获取积分摘要失败: %v\n", err)
		return
	}

	fmt.Printf("Steam ID: %d\n", summary.SteamID)
	fmt.Printf("当前积分: %d\n", summary.Points)
	fmt.Printf("用户等级: %d\n", summary.Level)
}

func showReactionConfig(client *Steam.Client) {
	fmt.Println("=== 反应配置 ===")
	fmt.Println("正在获取反应配置...")

	reactions, err := client.GetReactionConfig()
	if err != nil {
		fmt.Printf("获取反应配置失败: %v\n", err)
		return
	}

	fmt.Printf("共找到 %d 种反应类型:\n\n", len(reactions))
	for i, reaction := range reactions {
		fmt.Printf("%d. 反应ID: %d\n", i+1, reaction.ReactionID)
		fmt.Printf("   消耗积分: %d\n", reaction.PointsCost)
		fmt.Printf("   可用目标类型: %v\n", reaction.ValidTargetTypes)
		fmt.Println()

		// 限制显示数量，避免输出过多
		if i >= 9 {
			fmt.Printf("... 还有 %d 种反应类型\n", len(reactions)-i-1)
			break
		}
	}
}

func addReactionToUser(client *Steam.Client) {
	fmt.Println("=== 添加反应 ===")

	targetSteamIDStr := getInput("请输入目标用户Steam ID: ")
	targetSteamID, err := strconv.ParseUint(targetSteamIDStr, 10, 64)
	if err != nil {
		fmt.Printf("无效的Steam ID: %v\n", err)
		return
	}

	reactionTypeStr := getInput("请输入反应类型 (1=用户档案, 2=用户内容): ")
	reactionType, err := strconv.ParseUint(reactionTypeStr, 10, 32)
	if err != nil {
		fmt.Printf("无效的反应类型: %v\n", err)
		return
	}

	reactionIDStr := getInput("请输入反应ID: ")
	reactionID, err := strconv.ParseUint(reactionIDStr, 10, 32)
	if err != nil {
		fmt.Printf("无效的反应ID: %v\n", err)
		return
	}

	fmt.Printf("正在为用户 %d 添加反应 (类型: %d, ID: %d)...\n", 
		targetSteamID, reactionType, reactionID)

	// 这里需要从反应配置中获取积分消耗，暂时使用默认值
	pointsCost := int64(1) // 默认消耗1积分
	result, err := client.AddReaction(targetSteamID, uint32(reactionType), uint32(reactionID), pointsCost)
	if err != nil {
		fmt.Printf("添加反应失败: %v\n", err)
		return
	}

	if result.Success {
		fmt.Printf("成功添加反应！消耗积分: %d\n", result.PointsConsumed)
	} else {
		fmt.Println("添加反应失败，请检查参数或积分余额")
	}
}

func generateTokenCode(client *Steam.Client, sharedSecret string) {
	fmt.Println("=== 生成Steam Guard代码 ===")

	if sharedSecret == "" {
		fmt.Println("未配置Steam Guard共享密钥，无法生成验证码")
		return
	}

	fmt.Println("正在生成验证码...")
	code, err := client.GetTokenCode(sharedSecret)
	if err != nil {
		fmt.Printf("生成验证码失败: %v\n", err)
		return
	}

	fmt.Printf("当前Steam Guard验证码: %s\n", code)
	fmt.Println("注意: 验证码每30秒更新一次")
}

func showUserReactions(client *Steam.Client, steamID uint64) {
	fmt.Println("=== 用户反应记录 ===")
	fmt.Printf("正在获取用户 %d 的反应记录...\n", steamID)

	reactions, err := client.GetReactions(steamID, 0) // 0表示获取所有类型
	if err != nil {
		fmt.Printf("获取反应记录失败: %v\n", err)
		return
	}

	if reactions == nil {
		fmt.Println("未找到反应记录")
		return
	}

	fmt.Printf("反应记录获取成功，数据类型: %T\n", reactions)
	// 由于反应记录的具体结构取决于Steam API响应，这里只显示基本信息
	// 在实际使用中，可以根据具体的Protoc.ReactionsReceive结构来解析和显示详细信息
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}