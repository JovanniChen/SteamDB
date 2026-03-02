# SteamDB

`SteamDB` 是一个用于与 Steam 平台交互的 Go 项目，既可以作为库被其他项目引用，也包含 `main.go` 作为本地测试入口。

当前代码已覆盖登录、会话恢复、好友、库存与市场、购物车与交易流程、更新事件检测等能力。

## 功能概览

- 账号登录与会话管理
- Steam Guard 令牌生成
- 账户信息获取（昵称、余额、国家、语言）
- 积分系统（GetSummary、反应配置、加反应）
- 好友相关（好友码、链接、状态检查、删除）
- 库存/礼物获取、市场上架/下架/购买/订单
- 购物车与支付交易流程（初始化、最终支付、取消、价格查询、支付链接）
- 游戏更新事件抓取与本地 SQLite 去重保存

## 项目结构

```text
.
├── main.go               # 测试入口（按需调用 TestXxx）
├── Steam/
│   ├── client.go         # 对外客户端 API
│   ├── Dao/              # 底层 HTTP/认证/业务实现
│   ├── Model/            # 数据结构
│   ├── Protoc/           # protobuf 定义与生成代码
│   ├── Constants/        # 常量与 API 端点
│   └── Utils/            # 工具函数（如令牌生成）
├── mafiles/              # Steam Guard maFile（本地文件）
├── temp/                 # 会话缓存（session_*.json）
├── steam.db              # 更新事件 SQLite（运行时生成）
└── Makefile              # 构建与开发命令
```

## 环境要求

- Go 1.24+
- 可访问 Steam 相关域名的网络环境
- 部分功能需要有效的 `maFile`（如市场确认相关）

## 快速开始

### 1) 安装依赖

```bash
go mod download
```

### 2) 构建

```bash
make build
```

### 3) 运行测试入口

```bash
go run main.go
```

`main.go` 默认通过注释切换 `TestXxx` 方法。你可以按需要修改 `accountIndex`、代理配置和调用的方法。

## 作为库使用

### 基础登录

```go
package main

import (
	"fmt"
	"log"

	"github.com/JovanniChen/SteamDB/Steam"
)

func main() {
	client, err := Steam.NewClient(Steam.NewConfig("127.0.0.1:7890"))
	if err != nil {
		log.Fatal(err)
	}

	user, err := client.Login(&Steam.LoginCredentials{
		Username:     "your_username",
		Password:     "your_password",
		SharedSecret: "base64_shared_secret",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("登录成功: %s (%d)\n", user.Nickname, user.SteamID)
}
```

### 会话恢复（免重复登录）

```go
client, _ := Steam.NewClient(Steam.DefaultConfig())

client.SetLoginInfo(
	"username",
	7656119xxxxxxxxxx,
	"nickname",
	"CN",
	"access_token",
	"refresh_token",
	loginCookies,
	steamOffset,
	"schinese",
)

summary, err := client.GetPointsSummary(client.GetSteamID())
if err != nil {
	// handle error
}
_ = summary
```

## 常用 API（按场景）

- 登录与状态: `Login`, `SetLoginInfo`, `CheckLoginStatus`, `GetTokenCode`
- 账户信息: `GetUserInfo`, `GetBalance`, `GetWaitBalance`, `GetSteamID`, `GetNickname`, `GetAccessToken`
- 好友: `AddFriendByLink`, `AddFriendByFriendCode`, `CheckIsFriend`, `CheckFriendStatus`, `RemoveFriend`
- 市场与库存: `GetInventory`, `GetSteamGift`, `PutList`, `BuyListing`, `CreateOrder`, `GetMyListings`, `RemoveMyListings`, `GetConfirmations`
- 购物车与交易: `AddItemToCart`, `GetCart`, `ClearCart`, `ValidateCart`, `InitTransaction`, `InitConcurrentTransaction`, `GetFinalPrice`, `AccessCheckoutURL`, `GetAlipayURL`, `FinalizeTransaction`, `CancelTransaction`, `TransactionStatus`, `UnsendGift`, `UnsendAllGift`
- 更新检测: `GetGameUpdateEvents`（带数据库去重判断）

## Makefile 命令

```bash
make help         # 查看全部命令
make build        # 构建当前平台
make build-all    # 构建多平台产物
make test         # 运行测试（当前仓库无 *_test.go 时会快速结束）
make fmt          # go fmt
make vet          # go vet
make mod-tidy     # 整理依赖
make run          # go run .
```

## 说明与注意事项

- `main.go` 是本仓库的测试入口，不代表稳定 CLI 接口。
- 代理可通过 `Steam.NewConfig("host:port")` 或 `client.SetProxy(...)` 动态切换。
- `Config.Timeout` 字段已定义，但当前实现中底层 HTTP 超时由 `Dao` 内固定值控制（`10s`）。
- 游戏更新检测会在项目根目录自动创建/使用 `steam.db`。
- 该项目依赖 Steam 非公开接口行为，接口结构变化可能导致功能失效，建议在调用侧做好重试与降级。
