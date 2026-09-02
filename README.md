# SteamDB

SteamDB 是一个使用 Go 编写的 Steam 非官方客户端项目。项目既可以作为 Go 库调用，也提供了一个以 `--case` 为入口的本地调试程序，用于验证登录、Steam Guard、好友、购物车、购买、赠送、库存、社区市场、钱包、消费历史和游戏更新等流程。

> 重要说明：本项目调用了多个 Steam Web/API 非公开或非稳定接口。Steam 页面结构、protobuf 定义、Cookie 策略、风控规则和接口参数随时可能发生变化。本项目不保证接口长期稳定，也不应直接用于无人值守的生产资金操作。

## 当前能力

目前代码中已经实现或接入的主要能力如下。

### 登录与会话

- Steam 用户名和密码登录。
- RSA 公钥获取和密码加密。
- Steam Guard 共享密钥验证码生成。
- 登录轮询、访问令牌和刷新令牌获取。
- `steamLoginSecure`、`sessionid` 等多域名 Cookie 保存。
- 登录状态保存到本地 JSON session 文件。
- 从 session 文件恢复客户端状态，避免每次调用都重新登录。
- 获取 SteamID、昵称、国家代码、语言、钱包余额和待处理余额。
- 设置语言、退出所有设备和修改隐私设置。

### 网络与代理

- 无代理直连。
- HTTP/HTTPS 代理。
- 带用户名和密码的 HTTP/HTTPS 代理。
- SOCKS5/SOCKS5H 代理。
- 运行时通过 `Client.SetProxy` 切换代理。
- 按代理地址缓存 `http.Transport`，切换回已使用的代理时复用 transport。
- HTTP keep-alive 和连接池复用，不会为每个请求强制重新创建 TCP 连接。
- 请求失败重试和请求成功回调。

### 商店、购物车与结算

- 从 Steam 商店页面解析 package、bundle、价格和购买入口。
- 清空购物车、读取购物车、验证购物车。
- 使用 protobuf 调用 `IAccountCartService/AddItemsToCart`。
- 使用 protobuf 调用 `IAccountCartService/ModifyLineItem`。
- 给自己购买游戏。
- 通过 Steam AccountID 赠送游戏。
- 通过邮箱地址赠送游戏。
- 设置礼物留言和计划发送时间。
- 初始化普通交易和同时付交易。
- 获取最终价格和价格明细。
- 获取结算页面或支付宝跳转地址。
- 完成、取消和查询交易状态。
- Steam 钱包充值及指定国家充值。
- 修改商店国家和结算国家。
- 撤回单个或全部未领取礼物。

### 消费历史

- 请求 `/account/history/` 消费历史页面。
- 使用 XPath 解析交易 ID、日期、类型、支付方式、价格、钱包变更和退款状态。
- 支持同一笔交易中包含多个礼物商品和多个收礼人。
- 从最新记录开始筛选“连续、未退款、类型为礼物购买”的记录。
- 保存原始响应到 `store_purchase_history.html`，便于排查 Steam 页面结构变化。

### 好友、库存和社区市场

- 通过好友码或邀请链接添加好友。
- 接受好友请求、删除好友和检查好友关系。
- 创建好友邀请链接和解析好友链接信息。
- 获取库存并按可交易、可市场交易、commodity 等属性筛选。
- 获取 Steam 礼物库存。
- 获取自己的市场上架列表。
- 查询市场商品列表。
- 上架、下架和购买市场物品。
- 创建市场求购订单。
- 获取和处理移动确认。
- 获取交易对象库存并发送交易报价。
- 检查账号市场状态和按市场报价估算汇率。

### 积分和游戏更新

- 获取 Steam 积分摘要、反应配置和反应记录。
- 添加积分反应。
- 抓取 Steam 商店游戏更新事件。
- 提取 `event_type == 12` 的更新事件。
- 使用 SQLite 保存每个游戏的最新事件，并判断是否出现新更新。

## 技术栈

- Go `1.24.0`
- `net/http`：HTTP 请求、Cookie、连接池和代理
- `golang.org/x/net/proxy`：SOCKS5 代理
- `google.golang.org/protobuf`：Steam protobuf 请求和响应
- `github.com/antchfx/htmlquery`：HTML XPath 解析
- `github.com/mattn/go-sqlite3`：游戏更新事件持久化

## 环境要求

最低要求：

- Go `1.24+`
- 可访问 Steam 相关域名的网络环境
- Git

根据所使用的功能，可能还需要：

- `protoc`：修改 `.proto` 文件并重新生成 Go 文件时需要
- `protoc-gen-go`：生成 `*.pb.go` 时需要
- C 编译器：`go-sqlite3` 使用 CGO，完整构建和运行 SQLite 功能时需要
- Steam Mobile Authenticator 的 `.maFile`：市场上架、购买、交易报价和移动确认等操作需要
- Windows 交叉编译器：在 macOS/Linux 上生成支持 SQLite 的 Windows 二进制时需要

检查本机环境：

```bash
go version
protoc --version
protoc-gen-go --version
```

当前项目的 `go.mod` 指定：

```text
go 1.24.0
```

## 安全警告

本项目会处理以下高敏感数据：

- Steam 用户名和密码
- Steam Guard `shared_secret`
- maFile 中的 `identity_secret`、设备信息和会话信息
- access token 和 refresh token
- `steamLoginSecure` 和 `sessionid`
- 代理用户名和密码
- 好友邀请 token、交易报价 token 和交易 ID

请遵守以下原则：

1. 不要在 README、日志、截图、Issue 或聊天记录中粘贴真实令牌和账号密码。
2. 不要把 `mafiles/`、`temp/`、数据库或 session 文件提交到版本库。
3. 不要把真实代理密码直接硬编码到准备共享的源码中。
4. 如果敏感文件已经被 Git 跟踪，仅添加 `.gitignore` 不会将它从历史记录中移除。
5. 已经公开过的密码、shared secret、access token、refresh token 和代理密码应立即轮换。
6. 运行购买、充值、市场或交易报价 case 前，必须检查代码中硬编码的商品、价格、数量、收礼人和交易 ID。

当前 `.gitignore` 已包含 `*.maFile`、`temp/`、`*.db`、`main.go` 和 `main_accounts.go` 等规则，但已经被 Git 跟踪的文件仍然需要单独处理。

## 项目结构

```text
.
├── main.go                         # CLI 参数解析和程序入口，本地配置文件
├── main_accounts.go                # 本地调试账号列表，包含敏感信息
├── main_runner.go                  # case 注册表
├── main_cases_auth.go              # 登录、语言、国家、隐私等 case
├── main_cases_checkout.go          # 购物车、购买、交易和充值 case
├── main_cases_friend.go            # 好友操作 case
├── main_cases_market.go            # 库存、市场、确认和交易报价 case
├── main_cases_store.go             # 商店、积分、消费历史和游戏更新 case
├── main_cases_wallet.go            # 钱包和用户信息 case
├── main_session.go                 # session 序列化、保存和恢复
├── Makefile                        # 本地构建、测试和格式化命令
├── go.mod                          # Go 模块和依赖版本
├── Steam/
│   ├── client.go                   # 对外公开的 Client API
│   ├── Constants/
│   │   ├── api.go                  # Steam API 地址
│   │   ├── constants.go            # 游戏、库存 context 等常量
│   │   └── countries.go            # 国家、货币和 Steam currency ID
│   ├── Dao/
│   │   ├── dao.go                  # HTTP 客户端、代理、重试和 Cookie
│   │   ├── login.go                # 登录、令牌和会话状态
│   │   ├── cart.go                 # 购物车和 protobuf 接口
│   │   ├── checkout.go             # 结算、支付、充值和撤回礼物
│   │   ├── friend.go               # 好友相关接口
│   │   ├── market.go               # 库存、市场、确认和交易报价
│   │   ├── point.go                # Steam 积分系统
│   │   ├── store.go                # 商品信息和消费历史
│   │   ├── update.go               # 游戏更新事件抓取
│   │   ├── db.go                   # SQLite 更新事件存储
│   │   ├── user.go                 # 用户信息
│   │   └── time.go                 # Steam 时间同步
│   ├── Model/                      # 业务数据结构
│   ├── Protoc/                     # proto 定义和生成的 pb.go
│   ├── Param/                      # URL/Form 参数工具
│   ├── Errors/                     # Steam 错误码和错误定义
│   ├── Logger/                     # 日志工具
│   └── Utils/                      # maFile、确认参数和通用工具
├── mafiles/                        # 本地 maFile，禁止提交
├── temp/                           # session 缓存，禁止提交
├── build/                          # 构建产物
├── steam.db                        # 游戏更新事件数据库
└── store_purchase_history.html     # 消费历史原始响应调试文件
```

## 安装依赖

```bash
git clone <repository-url>
cd SteamDB
go mod download
```

验证项目：

```bash
go test ./...
go vet ./...
```

如果运行环境限制了默认 Go 缓存目录，可以临时指定缓存：

```bash
GOCACHE=/tmp/steamdb-go-cache go test ./...
```

## 客户端配置

### 无代理

```go
config := Steam.DefaultConfig()
client, err := Steam.NewClient(config)
```

也可以传入 `nil`，`NewClient` 会创建默认配置：

```go
client, err := Steam.NewClient(nil)
```

### HTTP 代理

无认证代理：

```go
config := Steam.NewConfig("127.0.0.1:7890")
```

无协议前缀时，项目为了兼容旧配置，会默认按 HTTP 代理处理。等价写法：

```go
config := Steam.NewConfig("http://127.0.0.1:7890")
```

带认证代理：

```go
config := Steam.NewConfig("http://username:password@proxy.example.com:8080")
```

HTTPS 代理 URL 也可以被解析：

```go
config := Steam.NewConfig("https://username:password@proxy.example.com:8443")
```

### SOCKS5 代理

无认证 SOCKS5：

```go
config := Steam.NewConfig("socks5://127.0.0.1:1080")
```

带认证 SOCKS5：

```go
config := Steam.NewConfig("socks5://username:password@proxy.example.com:1080")
```

项目同时接受 `socks5h://`：

```go
config := Steam.NewConfig("socks5h://username:password@proxy.example.com:1080")
```

### 运行时切换代理

```go
client.SetProxy("socks5://127.0.0.1:1080")

// 切回直连
client.SetProxy("")
```

每个 `Dao` 会按代理字符串缓存 transport。相同代理再次使用时会复用 transport 和空闲连接。HTTP 客户端不会在每次访问时主动新建完整连接；在服务端关闭连接、连接超时、网络变化或连接池无可用连接时才会建立新连接。

### 当前网络实现参数

代码中的实际网络参数如下：

- TCP 建连超时：15 秒
- TCP keep-alive：120 秒
- TLS 握手超时：3 秒
- HTTP 客户端整体超时：10 秒
- 默认请求重试次数：3 次
- 最大空闲连接数：5
- 每个主机最大空闲连接数：10
- 每个主机最大连接数：60
- HTTP/2：当前被禁用
- TLS 证书验证：当前使用 `InsecureSkipVerify: true`

需要注意：`Steam.Config.Timeout` 当前虽然存在并默认设置为 30 秒，但 `Dao.New` 尚未使用该字段，实际 HTTP 客户端超时仍是固定的 10 秒。TLS 证书验证被关闭也不适合生产环境，尤其不适合在不可信代理上发送登录凭据。

## 配置本地账号

CLI case 通过 `main_accounts.go` 中的账号索引获取账号：

```go
var accounts = []Account{
    {
        Username:     "your_username",
        Password:     "your_password",
        SharedSecret: "your_base64_shared_secret",
    },
}
```

`--account 0` 对应数组中的第 0 个账号。

推荐后续将账号配置改为环境变量或本地加密配置，而不是继续维护硬编码账号数组。至少应确保该文件不再进入提交历史。

## maFile

移动确认、市场操作和交易报价相关 case 会按以下路径读取 maFile：

```text
mafiles/<Steam用户名>.maFile
```

例如：

```text
mafiles/example_user.maFile
```

maFile 通常包含移动认证器的重要密钥。不要上传、分享或记录完整内容。

当前 `Steam.LoginCredentials` 也保留了 `MaFile` 字段，但 `Client.Login` 的实际底层认证调用使用的是用户名、密码和 `SharedSecret`；`MaFile` 会被放入返回的 `UserInfo`，登录流程本身不会直接解析该字段。CLI 的 `TestLogin` 仍要求对应 maFile 文件存在，这是当前调试入口的实现约束。

## 登录与 session

### 第一次登录

```bash
go run . --case TestLogin --account 0
```

登录成功后会保存：

```text
temp/session_0.json
```

session 内容包括：

- 账号索引
- 用户名
- SteamID
- 昵称
- 国家代码
- access token
- refresh token
- 各 Steam 域名 Cookie
- 登录时间
- Steam 时间偏移
- Steam 语言

### session 有效期策略

当前 `TestLogin` 使用文件修改时间判断 session 是否需要重新登录：

- session 文件不存在：执行登录
- session 文件修改时间不超过 4 小时：跳过登录
- session 文件超过 4 小时：删除文件并重新登录

这个 4 小时是本地调试策略，不代表 Steam token 或 Cookie 的官方有效期。Steam 仍可能提前使 session 失效。如果接口返回登录页、Cookie 不存在或认证失败，应删除对应 session 后重新执行登录：

```bash
rm temp/session_0.json
go run . --case TestLogin --account 0
```

除 `TestLogin`、`TestGetTokenCode`、`TestGetPackageDetails` 和游戏更新等少数无需恢复登录态的 case 外，大多数 case 都会调用 `loadFromSession`。

## CLI 使用方式

查看所有可用 case：

```bash
go run . --list-cases
```

通用格式：

```bash
go run . --case <CaseName> --account <AccountIndex>
```

示例：

```bash
go run . --case TestLogin --account 0
go run . --case TestGetUserInfo --account 0
go run . --case TestGetStorePurchaseHistory --account 0
```

游戏更新 case 使用 `--game-id`：

```bash
go run . --case TestGetGameUpdateInofs --game-id 1879330
```

未指定 `--account` 时默认值是 `3`；未指定 `--game-id` 时默认值是 `1879330`。

## Case 列表

下面的 case 是调试入口，不是稳定的命令行产品接口。很多参数仍直接写在 `main_cases_*.go` 中，运行前必须修改为自己的测试数据。

风险标记：

- 只读：原则上只查询数据，但仍可能触发登录风控或限流。
- 修改：会改变好友、语言、国家、隐私、购物车或订单状态。
- 资金：可能创建真实购买、充值、市场订单或完成支付。
- 高风险：可能撤回礼物、退出所有设备、发送交易报价或批量操作。

### 登录、账号和钱包

| Case | 风险 | 说明 |
| --- | --- | --- |
| `TestGetTokenCode` | 只读 | 使用 shared secret 生成当前 Steam Guard 验证码 |
| `TestLogin` | 修改 | 登录并将完整 session 保存到 `temp/session_<index>.json` |
| `TestCheckAccountAvailable` | 只读 | 检查社区市场可用状态 |
| `TestSetLanguage` | 修改 | 将 Steam 语言设置为代码中指定值 |
| `TestGetCountryCode` | 只读 | 输出当前账号国家代码 |
| `TestGetUserInfo` | 只读 | 获取用户信息、钱包余额和待处理余额 |
| `TestGetBalance` | 只读 | 获取钱包余额，单位为 Steam 返回的最小货币单位 |
| `TestGetWaitBalance` | 只读 | 获取待处理钱包余额 |
| `TestLogoutAll` | 高风险 | 退出所有设备，可能使现有 session 全部失效 |
| `TestSetPrivacy` | 修改 | 修改当前账号的隐私设置 |

### 好友

| Case | 风险 | 说明 |
| --- | --- | --- |
| `TestCreateFriendLink` | 修改 | 创建新的好友邀请链接 |
| `TestAddFriendByFriendCode` | 修改 | 使用代码中硬编码的好友码添加好友 |
| `TestAddFriendByLink` | 修改 | 使用代码中的邀请链接添加好友 |
| `TestGetFriendInfoByLink` | 只读 | 解析邀请链接对应的用户和 invite token |
| `TestGetFriendInfoByLinkAndAddFriend` | 修改 | 解析链接后添加对应好友 |
| `TestAcceptFriend` | 修改 | 接受指定 SteamID 的好友请求 |
| `TestRemoveFriend` | 修改 | 删除指定 SteamID 的好友 |
| `TestCheckIsFriend` | 只读 | 检查与指定 SteamID 的好友关系 |
| `TestCheckFriendStatus` | 只读 | 检查好友邀请链接状态 |

### 购物车、购买和充值

| Case | 风险 | 说明 |
| --- | --- | --- |
| `TestGetProductByAppUrl` | 只读 | 解析商店 App 页面中的 package/bundle 和价格 |
| `TestGetPackageDetails` | 只读 | 请求 package details，并直接打印服务器响应 |
| `TestGetCart` | 只读 | 获取购物车，当前实现主要用于调试输出 |
| `TestClearCart` | 修改 | 清空当前账号购物车 |
| `TestAddItemToCart` | 修改 | 向购物车加入礼物商品，可按 AccountID 或邮箱赠送 |
| `TestAddItemToCartWithSentTime` | 修改 | 添加带计划发送时间的礼物 |
| `TestValidateCart` | 只读 | 请求 Checkout ValidateCart，当前实现打印原始结果 |
| `TestBuyGameToSelf` | 资金 | 创建给自己购买的交易并生成结算链接，示例中随后取消交易 |
| `TestBuyGameToOther` | 资金 | 创建赠送他人的交易并生成结算链接，示例中随后取消交易 |
| `TestConcurrentPayment` | 高风险 | 初始化同时付交易并并发完成交易，必须先审查代码 |
| `TestGetFinalPrice` | 只读 | 获取硬编码交易 ID 的最终价格 |
| `TestAccess` | 资金 | 获取硬编码交易的结算页面或外部支付入口 |
| `TestCancelTransaction` | 修改 | 取消硬编码交易 ID |
| `TestTransactionStatus` | 只读 | 轮询交易状态 |
| `TestTestGetPayLinkAgain` | 资金 | 重新获取外部支付链接，并在 case 内切换代理 |
| `TestAddFunds` | 资金 | 创建 Steam 钱包充值流程 |
| `TestAddFundsWithCountry` | 资金 | 使用指定国家创建钱包充值流程 |
| `TestSetCountry` | 修改 | 修改商店国家和结算国家 |

### 库存、市场和交易报价

| Case | 风险 | 说明 |
| --- | --- | --- |
| `TestGetInventory` | 只读 | 获取并筛选库存 |
| `TestGetSteamGift` | 只读 | 获取 Steam 礼物库存、名称、市场名称和收礼人 |
| `TestGetMyListings` | 只读 | 获取已上架和待确认物品 |
| `TestGetMarketListings` | 只读 | 获取指定市场商品的挂单列表 |
| `TestGetSteamRate` | 只读 | 按不同国家和货币市场价格估算汇率 |
| `TestIsAccountBanned` | 只读 | 检查账号是否存在市场限制或红信状态 |
| `TestPutList` | 资金 | 随机选择符合条件的库存物品并按硬编码价格上架 |
| `TestRemoveMyListings` | 修改 | 下架硬编码 creator/listing ID |
| `TestBuyListing` | 资金 | 购买硬编码市场挂单并处理移动确认 |
| `TestCreateOrder` | 资金 | 创建指定商品、价格、数量的求购订单 |
| `TestGetConfirmations` | 只读 | 获取移动确认列表，需要 maFile |
| `TestUnsendGift` | 高风险 | 遍历礼物库存并逐个撤回 |
| `TestUnsendAllGift` | 高风险 | 撤回所有符合条件的未发送礼物 |
| `TestGetPartnerInventory` | 只读 | 获取交易对象库存 |
| `TestSendGift` | 高风险 | 从库存选取物品并向指定交易链接发送报价 |

### 商店、积分和游戏更新

| Case | 风险 | 说明 |
| --- | --- | --- |
| `TestGetSummary` | 只读 | 获取 Steam 积分摘要，同时会输出部分 Cookie 调试信息 |
| `TestGetStorePurchaseHistory` | 只读 | 解析最新消费历史，并打印连续未退款礼物购买记录 |
| `TestGetGameUpdateInofs` | 只读 | 抓取游戏更新事件，写入 SQLite 并判断是否有新事件 |

## 作为 Go 库使用

模块路径：

```text
github.com/JovanniChen/SteamDB
```

创建客户端：

```go
package main

import (
    "log"

    "github.com/JovanniChen/SteamDB/Steam"
)

func main() {
    client, err := Steam.NewClient(Steam.DefaultConfig())
    if err != nil {
        log.Fatal(err)
    }

    _ = client
}
```

### 登录示例

```go
package main

import (
    "fmt"

    "github.com/JovanniChen/SteamDB/Steam"
)

func run(client *Steam.Client) error {
	maFileContent := "..."

	userInfo, err := client.Login(&Steam.LoginCredentials{
		Username:     "your_username",
		Password:     "your_password",
		SharedSecret: "your_shared_secret",
		MaFile:       maFileContent,
	})
	if err != nil {
		return err
	}

	fmt.Println(userInfo.SteamID)
	fmt.Println(userInfo.Nickname)
	return nil
}
```

`Client.Login` 本身不会自动写入 `temp/session_*.json`。session 文件的保存是 CLI 中 `TestLogin` 和 `SteamSession.Save` 实现的应用层逻辑。

### 请求计数回调

```go
import "sync/atomic"

var requests atomic.Int64

client.SetRequestCallback(func() {
    requests.Add(1)
})
```

回调会在 `RetryRequest` 获得 HTTP 响应后执行，可用于请求统计或代理计费统计。

## 购物车与礼物

`Model.AddCartItem` 当前结构：

```go
type AddCartItem struct {
    PackageID       uint32
    BundleID        uint32
    AccountidGiftee uint32
    Message         string
    EmailGiftee     string
}
```

### 给 Steam 用户赠送

```go
err := client.AddItemToCart([]Model.AddCartItem{
    {
        PackageID:       645485,
        AccountidGiftee: 739009475,
        Message:         "Enjoy the game",
    },
})
```

这里的 `AccountidGiftee` 是 Steam AccountID，不是 SteamID64。

### 通过邮箱赠送

```go
err := client.AddItemToCart([]Model.AddCartItem{
    {
        PackageID:   645485,
        EmailGiftee: "recipient@example.com",
        Message:     "Enjoy the game",
    },
})
```

收礼人规则：

- `AccountidGiftee != 0` 且 `EmailGiftee` 为空：按 Steam AccountID 赠送。
- `EmailGiftee` 非空且 `AccountidGiftee == 0`：按邮箱赠送。
- 两者同时填写：返回错误。
- 两者都不填写：礼物接口返回错误。
- 邮箱会先执行 `strings.TrimSpace`。

给自己购买应使用 `AddItemToCartSelf`，不需要填写收礼人：

```go
err := client.AddItemToCartSelf([]Model.AddCartItem{
    {PackageID: 645485},
})
```

### 计划发送礼物

```go
sentTime := int32(time.Now().Add(24 * time.Hour).Unix())

err := client.AddItemToCartWithSentTime(
    []Model.AddCartItem{
        {
            PackageID:       645485,
            AccountidGiftee: 739009475,
            Message:         "Happy birthday",
        },
    },
    sentTime,
)
```

`sentTime` 是 Unix 时间戳，当前 protobuf 字段类型是 `int32`。调用方需要注意 2038 年问题和整数溢出。

### 修改购物车行项目

`ModifyLineItem` 使用 protobuf 调用：

```go
response, err := client.ModifyLineItem(&Protoc.ModifyCartSend{
    LineItemId:  lineItemID,
    UserCountry: "CN",
    GiftInfo: &Protoc.GiftInfo{
        EmailGiftee: "recipient@example.com",
        GiftMessage: &Protoc.GiftMessage{
            Message: "Updated message",
        },
    },
    Flag: &Protoc.Flag{
        IsGift: true,
    },
})
```

如果 `UserCountry` 为空，DAO 会使用当前登录账号的国家代码。请求体会被 protobuf 序列化后以 Base64 形式放入 `input_protobuf_encoded` 表单字段，access token 放在查询参数中。

## 购买流程

一个典型购买流程如下：

1. `ClearCart()` 清理旧购物车。
2. `AddItemToCartSelf()` 或 `AddItemToCart()` 添加商品。
3. `ValidateCart()` 验证购物车。
4. `InitTransaction()` 初始化交易并获得 transaction ID。
5. `GetFinalPriceWithDetails()` 核对最终金额、钱包支付额和外部支付额。
6. `AccessCheckoutURL()` 或 `GetAlipayURL()` 获取支付入口。
7. 根据业务选择 `FinalizeTransaction()` 或 `CancelTransaction()`。
8. 使用 `TransactionStatus()` 查询交易最终状态。

示例：

```go
if err := client.ClearCart(); err != nil {
    return err
}

if err := client.AddItemToCartSelf([]Model.AddCartItem{
    {PackageID: 645485},
}); err != nil {
    return err
}

transactionID, err := client.InitTransaction("alipay", "CN", 0)
if err != nil {
    return err
}

details, err := client.GetFinalPriceWithDetails(transactionID)
if err != nil {
    return err
}

fmt.Printf("%+v\n", details)
```

Steam 可能因账号国家、钱包币种、购买限制、待处理订单、价格变化或风控返回错误。不要在未核对最终价格时自动完成交易。

## 消费历史

调用：

```go
result, err := client.GetStorePurchaseHistory()
```

返回结构：

```go
type StorePurchaseHistoryResult struct {
    TotalRecords                  int
    LatestUnrefundedGiftPurchases []StorePurchaseHistoryRecord
}

type StorePurchaseHistoryRecord struct {
    Index           int
    TransactionID   string
    Date            string
    Items           []string
    Receivers       []string
    TransactionType string
    Payment         string
    BasePrice       string
    Tax             string
    Shipping        string
    Total           string
    WalletChange    string
    WalletBalance   string
    Refunded        bool
}
```

同一笔消费记录可能包含多个礼物：

```go
for _, record := range result.LatestUnrefundedGiftPurchases {
    for i := 0; i < len(record.Items) || i < len(record.Receivers); i++ {
        if i < len(record.Items) {
            fmt.Println("item:", record.Items[i])
        }
        if i < len(record.Receivers) {
            fmt.Println("receiver:", record.Receivers[i])
        }
    }
}
```

当前筛选规则严格按页面从新到旧执行：

```text
TransactionType == "礼物购买" && Refunded == false
```

遇到第一条不满足条件的记录后立即停止，因此返回的是“从最新一条开始连续的一段未退款礼物购买”，不是页面中所有未退款礼物记录。

当前限制：

- 只解析首次请求 `/account/history/` 返回的记录。
- 尚未调用 `/account/AjaxLoadMoreHistory/` 加载后续分页。
- `TotalRecords` 表示当前页面已解析的行数，不是账号全部历史记录数。
- 原始 HTML 会写入项目根目录的 `store_purchase_history.html`。
- 如果 session 失效，Steam 可能返回登录页，解析器会提示未获取到消费历史表格。
- HTML 结构变化可能导致 XPath 失效，需要使用保存的原始 HTML 调整解析器。

项目包含多礼物记录回归测试：

```bash
go test ./Steam/Dao -run TestParseStorePurchaseHistoryRowParsesMultipleGifts -v
```

## Protobuf

proto 文件位于：

```text
Steam/Protoc/login.proto
Steam/Protoc/cart.proto
Steam/Protoc/store.proto
Steam/Protoc/point.proto
```

生成文件位于同一目录下的 `*.pb.go`。这些文件应由 `protoc` 生成，不应长期手工修改。

安装 Go 插件：

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

重新生成单个文件：

```bash
cd Steam/Protoc
protoc --go_out=. cart.proto
```

重新生成全部 proto：

```bash
cd Steam/Protoc
protoc --go_out=. login.proto cart.proto store.proto point.proto
```

`cart.proto` 当前 `GiftInfo` 定义：

```proto
message GiftInfo {
    int32 accountid_giftee = 1;
    GiftMessage gift_message = 2;
    int32 time_scheduled_send = 3;
    string email_giftee = 4;
}
```

修改字段后必须重新生成 `cart.pb.go`，并执行：

```bash
gofmt -w Steam/Protoc/cart.pb.go
go test ./...
```

生成文件头中的 `protoc` 和 `protoc-gen-go` 版本可能因开发机版本不同而变化，这是生成器版本变化造成的正常差异，但提交前应检查是否产生了无关的大范围 diff。

## 游戏更新和 SQLite

`GetGameUpdateEvents` 会：

1. 请求 Steam 商店新闻页面。
2. 从页面 `data-initialevents` 属性解析 JSON。
3. 提取 `event_type == 12` 的事件。
4. 读取 SQLite 中该游戏最新事件。
5. 比较 `UniqueID`。
6. 有变化时写入 `game_update_events` 表并返回 `needUpdate == true`。

默认数据库：

```text
steam.db
```

表结构由程序自动创建：

```sql
CREATE TABLE IF NOT EXISTS game_update_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL,
    unique_id TEXT NOT NULL,
    app_id INTEGER NOT NULL,
    start_time INTEGER NOT NULL,
    event_name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(game_id, unique_id)
);
```

可以在数据库首次初始化前修改路径：

```go
Dao.SetDBPath("data/steam.db")
```

`dbOnce` 会保证数据库只初始化一次，因此应在第一次调用更新查询前设置路径。

## Makefile

查看帮助：

```bash
make help
```

当前目标：

| 命令 | 作用 |
| --- | --- |
| `make run` | 执行 `go run .`，仍需通过参数指定 case 才会执行具体逻辑 |
| `make build` | 构建当前平台二进制到 `build/SteamDB` |
| `make test` | 执行 `go test ./...` |
| `make fmt` | 执行 `go fmt ./...` |
| `make vet` | 执行 `go vet ./...` |
| `make mod-tidy` | 执行 `go mod tidy` 和 `go mod verify` |
| `make clean` | 删除 `build/` |

`make run` 不会自动附加 `--case`，直接运行时会提示使用 `--list-cases`。

## 构建

### 当前平台

```bash
make build
```

或者：

```bash
go build -trimpath -o build/SteamDB .
```

### Windows 原生构建

项目依赖 `go-sqlite3`，因此完整功能需要 CGO 和 C 编译器。在 Windows 上安装可用的 GCC 工具链后执行：

```powershell
$env:CGO_ENABLED = "1"
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -trimpath -o build\SteamDB-windows-amd64.exe .
```

### 在 macOS 上交叉编译 Windows amd64

安装 MinGW-w64：

```bash
brew install mingw-w64
```

构建：

```bash
mkdir -p build
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CC=x86_64-w64-mingw32-gcc \
go build -trimpath -o build/SteamDB-windows-amd64.exe .
```

### 在 Linux 上交叉编译 Windows amd64

Debian/Ubuntu 安装工具链：

```bash
sudo apt-get update
sudo apt-get install -y gcc-mingw-w64-x86-64
```

构建：

```bash
mkdir -p build
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CC=x86_64-w64-mingw32-gcc \
go build -trimpath -o build/SteamDB-windows-amd64.exe .
```

只设置 `CGO_ENABLED=0` 可能生成二进制，但 `go-sqlite3` 的真实数据库功能不可用，游戏更新持久化会在运行时失败。因此当前项目不建议使用禁用 CGO 的 Windows 构建作为完整版本。

## 测试和代码检查

完整测试：

```bash
go test ./...
```

详细输出：

```bash
go test -v ./...
```

仅测试 DAO：

```bash
go test -v ./Steam/Dao
```

格式化和静态检查：

```bash
go fmt ./...
go vet ./...
```

整理依赖：

```bash
go mod tidy
go mod verify
```

注意：多个 case 会访问真实 Steam 服务，不能把它们当作普通单元测试自动运行。`go test ./...` 当前主要覆盖不需要真实账号的解析逻辑。

## 日志与调试文件

调试时可能产生：

- `store_purchase_history.html`：消费历史原始页面
- `index.html` 或其他手工保存的 Steam 页面
- `steam.db`：游戏更新数据库
- `temp/session_*.json`：登录状态
- 日志文件：由 Logger 配置决定

这些文件可能包含账号标识、交易记录、Cookie、余额或购买信息。分享日志前必须脱敏。

## 当前已知限制

1. Steam 非公开接口和 HTML 页面可能随时变化。
2. 消费历史仅解析第一页，尚未实现 `AjaxLoadMoreHistory` 分页。
3. 多个 DAO 方法仍直接打印原始响应，尚未全部改为结构化返回值。
4. 多个 CLI case 包含硬编码 SteamID、AccountID、好友链接、商品 ID、价格、数量和交易 ID。
5. `Config.Timeout` 当前没有传入底层 HTTP 客户端。
6. HTTP 客户端当前关闭 TLS 证书校验，不适合生产安全要求。
7. HTTP/2 当前被禁用。
8. 重试逻辑主要处理网络错误；非 200 状态重试代码目前被注释。
9. session 文件使用明文 JSON 存储 token 和 Cookie，没有加密。
10. session 是否重新登录只按文件修改时间 4 小时判断，不能保证 Steam session 实际有效。
11. 市场和购买流程会受到地区、币种、账号限制、Steam Guard、风控和速率限制影响。
12. SQLite 驱动依赖 CGO，跨平台构建需要对应 C 工具链。
13. 游戏更新抓取使用固定 XPath，Steam 页面布局变化后可能失效。
14. `GetStorePurchaseHistory` 当前会覆盖写入根目录调试 HTML，且写文件错误未向上返回。
15. 部分命名保留了历史拼写，例如 `GetGameUpdateInofs` 和 `Catetory`，调用时需要使用代码中的实际名称。

## 常见问题

### 提示 `Cookie not exist`

当前 session 中缺少目标 Steam 域名的 Cookie。删除 session 并重新登录，确认登录流程为 `store.steampowered.com`、`steamcommunity.com` 和 `checkout.steampowered.com` 保存了 Cookie。

### 消费历史返回登录页

通常是 session 失效。重新执行：

```bash
rm temp/session_<index>.json
go run . --case TestLogin --account <index>
```

### SOCKS5 代理无法连接

检查以下项目：

- URL 是否包含 `socks5://` 或 `socks5h://`
- host 和端口是否正确
- 用户名和密码是否需要 URL 编码
- 代理是否允许访问 Steam 域名和 443 端口
- 代理出口国家是否与账号商店国家和付款流程冲突

### HTTP 代理旧格式还能否使用

可以。以下格式仍按 HTTP 代理解析：

```text
host:port
username:password@host:port
```

### 每次请求都会重新建立连接吗

不会。transport 启用了连接池和 keep-alive，正常情况下会复用现有连接。但连接空闲超时、代理断开、服务端关闭、网络切换或连接池无可用连接时会建立新连接。

### protobuf 返回解析错误

优先检查：

- proto 字段号和字段类型是否与 Steam 当前接口一致
- 请求是否真的返回 protobuf，而不是 HTML 登录页或错误页
- `X-Eresult` 是否为 `1`
- access token 是否有效
- 是否使用正确的 Base64 编码和表单字段名
- 代理是否替换或截断响应

### Windows 构建提示 C 编译器错误

`go-sqlite3` 依赖 CGO。Windows 本机构建需要 GCC，macOS/Linux 交叉编译需要 MinGW-w64，并且要为 `CC` 指定目标平台编译器。

## 开发建议

- 新增 Steam 接口时，把 endpoint 放入 `Steam/Constants/api.go`。
- 业务参数和响应结构放入 `Steam/Model`。
- protobuf 定义放入 `Steam/Protoc` 并重新生成 `*.pb.go`。
- 网络请求和解析实现放入对应的 `Steam/Dao` 文件。
- 在 `Steam/client.go` 暴露稳定的客户端方法。
- 仅用于人工验证的调用放入对应 `main_cases_*.go`。
- 新增 case 后在 `main_runner.go` 注册。
- HTML 解析至少加入一个本地 fixture 或最小 HTML 单元测试。
- 涉及真实资金的流程应拆分为“准备、核价、确认”阶段，不要默认自动 finalize。
- 新增日志时避免输出 token、Cookie、密码、maFile 和完整支付 URL。

## 免责声明

本项目仅用于技术研究和自有账号测试。使用者需要自行承担账号限制、交易损失、地区或支付合规、Steam 服务条款以及接口变更带来的风险。请勿用于未经授权的账号、欺诈、绕过平台限制或其他违法违规用途。
