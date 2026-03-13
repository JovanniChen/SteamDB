# SteamDB

`SteamDB` 是一个用于与 Steam 平台交互的 Go 项目，可作为库使用，也内置了本地调试入口（`go run . --case ...`）。

当前主流程覆盖：

- 登录与会话恢复
- Steam Guard 令牌
- 好友操作
- 购物车/交易/支付流程
- 库存与市场相关操作
- 钱包余额查询
- 游戏更新事件抓取

## 环境要求

- Go `1.24+`
- 能访问 Steam 相关域名的网络环境
- 某些 case 需要本地 `mafiles/<username>.maFile`

## 当前项目结构（与入口相关）

```text
.
├── main.go                  # CLI 入口：--case/--account/--game-id
├── main_runner.go           # case 注册表
├── main_cases_auth.go       # 登录、令牌、语言、账号可用性等
├── main_cases_checkout.go   # 购物车、交易、支付流程
├── main_cases_friend.go     # 好友相关
├── main_cases_market.go     # 库存、礼物、市场、订单
├── main_cases_store.go      # 商店信息、更新事件等
├── main_cases_wallet.go     # 钱包余额相关
├── main_session.go          # session 读写与恢复
├── main_accounts.go         # 本地账号列表（调试用）
├── mafiles/                 # maFile 文件目录
└── temp/                    # session 缓存目录
```

## 快速开始

```bash
go mod download
go run . --list-cases
```

### 通用执行方式

```bash
go run . --case <CaseName> --account <AccountIndex>
```

示例：

```bash
go run . --case TestLogin --account 3
```

### 特殊 case：游戏更新事件

`TestGetGameUpdateInofs` 不走 `--account`，需要 `--game-id`：

```bash
go run . --case TestGetGameUpdateInofs --game-id 1879330
```

## 全部 case 的 go run 执行命令

以下命令可直接复制执行（示例都使用 `--account 3`，按需替换）：

```bash
go run . --case TestAccess --account 3
go run . --case TestAddFriendByFriendCode --account 3
go run . --case TestAddFriendByLink --account 3
go run . --case TestAddFunds --account 3
go run . --case TestAddItemToCart --account 3
go run . --case TestAddItemToCartAndInitTransaction --account 3
go run . --case TestBuyListing --account 3
go run . --case TestCancelTransaction --account 3
go run . --case TestCheckAccountAvailable --account 3
go run . --case TestCheckFriendStatus --account 3
go run . --case TestCheckIsFriend --account 3
go run . --case TestClearCart --account 3
go run . --case TestConcurrentPayment --account 3
go run . --case TestCreateOrder --account 3
go run . --case TestGetBalance --account 3
go run . --case TestGetCart --account 3
go run . --case TestGetConfirmations --account 3
go run . --case TestGetFinalPrice --account 3
go run . --case TestGetFriendInfoByLink --account 3
go run . --case TestGetFriendInfoByLinkAndAddFriend --account 3
go run . --case TestGetGameUpdateInofs --game-id 1879330
go run . --case TestGetInventory --account 3
go run . --case TestGetMyListings --account 3
go run . --case TestGetProductByAppUrl --account 3
go run . --case TestGetSteamGift --account 3
go run . --case TestGetSummary --account 3
go run . --case TestGetTokenCode --account 3
go run . --case TestGetWaitBalance --account 3
go run . --case TestInitTransaction --account 3
go run . --case TestLogin --account 3
go run . --case TestLogoutAll --account 3
go run . --case TestPutList --account 3
go run . --case TestRemoveFriend --account 3
go run . --case TestRemoveMyListings --account 3
go run . --case TestSetLanguage --account 3
go run . --case TestTestGetPayLinkAgain --account 3
go run . --case TestTransactionStatus --account 3
go run . --case TestUnsendAllGift --account 3
go run . --case TestUnsendGift --account 3
go run . --case TestValidateCart --account 3
```

## Makefile 命令

日常推荐（最常用）：

```bash
make run
make test
make build
make clean
```

全部可用目标：

```bash
make help
make run
make build
make test
make fmt
make vet
make mod-tidy
make clean
```

## 注意事项

- `main_accounts.go` 当前包含调试账号信息，建议改为本地配置文件并加入 `.gitignore`。
- 市场/确认类 case 依赖 `mafiles`，缺少文件会直接失败。
- 该项目依赖 Steam 非公开接口，接口变化可能导致 case 失效。
