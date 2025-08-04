# SteamDB Go客户端库

一个用于与Steam平台交互的Go语言第三方库，提供简单易用的API接口。

## 主要功能

- **用户认证**: 支持Steam用户名密码登录，包括Steam Guard双因素认证
- **积分系统**: 查询用户积分余额、等级信息
- **反应系统**: 获取反应配置，为用户添加反应（表情、奖励等）
- **令牌生成**: 生成Steam Guard验证码
- **会话管理**: 自动管理登录会话和Cookie
- **错误处理**: 完善的错误处理和重试机制

## 快速开始

### 安装

```bash
go get github.com/steamdb/steamdb-go
```

### 基本使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/steamdb/steamdb-go/Steam"
)

func main() {
    // 创建客户端
    client, err := Steam.NewClient(Steam.DefaultConfig())
    if err != nil {
        log.Fatal(err)
    }
    
    // 登录
    credentials := &Steam.LoginCredentials{
        Username:     "your_username",
        Password:     "your_password", 
        SharedSecret: "your_shared_secret", // 可选，Steam Guard共享密钥
    }
    
    userInfo, err := client.Login(credentials)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("登录成功！Steam ID: %d\n", userInfo.SteamID)
}
```

## API 参考

### 客户端配置

#### Steam.Config

客户端配置结构体，用于初始化Steam客户端。

```go
type Config struct {
    Proxy   string        // 代理服务器地址，格式: "host:port"
    Timeout time.Duration // 请求超时时间
}
```

#### Steam.DefaultConfig()

返回默认配置。

```go
config := Steam.DefaultConfig()
// 默认配置：无代理，30秒超时
```

#### Steam.NewClient(config *Config)

创建新的Steam客户端实例。

```go
client, err := Steam.NewClient(config)
```

### 用户认证

#### LoginCredentials

登录凭据结构体。

```go
type LoginCredentials struct {
    Username     string // Steam用户名
    Password     string // Steam密码
    SharedSecret string // Steam Guard共享密钥(base64编码)，可选
}
```

#### UserInfo

用户信息结构体，包含登录成功后的用户详细信息。

```go
type UserInfo struct {
    SteamID      uint64 // Steam ID
    Username     string // 用户名
    Nickname     string // 昵称
    AccessToken  string // 访问令牌
    RefreshToken string // 刷新令牌
    CountryCode  string // 国家代码
}
```

#### Client.Login(credentials *LoginCredentials)

执行Steam登录。

```go
userInfo, err := client.Login(credentials)
```

### Steam Guard

#### Client.GetTokenCode(sharedSecret string)

生成Steam Guard验证码。

```go
code, err := client.GetTokenCode("your_shared_secret")
fmt.Println("验证码:", code) // 输出6位数字验证码
```

### 积分系统

#### PointsSummary

积分摘要信息结构体。

```go
type PointsSummary struct {
    SteamID uint64 // Steam ID
    Points  int64  // 当前积分数量
    Level   int32  // 用户等级
}
```

#### Client.GetPointsSummary(steamID uint64)

获取用户积分摘要。

```go
summary, err := client.GetPointsSummary(userInfo.SteamID)
fmt.Printf("积分: %d, 等级: %d\n", summary.Points, summary.Level)
```

### 反应系统

#### ReactionConfig

反应配置信息结构体。

```go
type ReactionConfig struct {
    ReactionID       uint32   // 反应ID
    PointsCost       int64    // 消耗积分数量
    ValidTargetTypes []uint32 // 可用的目标类型列表
}
```

#### Client.GetReactionConfig()

获取所有可用的反应配置。

```go
reactions, err := client.GetReactionConfig()
for _, reaction := range reactions {
    fmt.Printf("反应ID: %d, 消耗积分: %d\n", 
        reaction.ReactionID, reaction.PointsCost)
}
```

#### AddReactionResult

添加反应结果结构体。

```go
type AddReactionResult struct {
    Success        bool  // 操作是否成功
    PointsConsumed int64 // 消耗的积分数量
}
```

#### Client.AddReaction(targetSteamID uint64, reactionType uint32, reactionID uint32)

为指定用户添加反应。

```go
result, err := client.AddReaction(targetSteamID, 1, 23)
if err != nil {
    log.Printf("添加反应失败: %v", err)
} else if result.Success {
    fmt.Printf("成功添加反应，消耗积分: %d\n", result.PointsConsumed)
}
```

#### Client.GetReactions(steamID uint64, reactionType uint32)

获取用户的反应记录。

```go
reactions, err := client.GetReactions(steamID, 0) // 0表示获取所有类型
```

### 辅助方法

#### Client.GetSteamID()

获取当前登录用户的Steam ID。

```go
steamID := client.GetSteamID()
```

#### Client.GetAccessToken()

获取当前有效的访问令牌。

```go
token, err := client.GetAccessToken()
```

#### Client.GetNickname()

获取用户昵称。

```go
nickname := client.GetNickname()
```

#### Client.GetCountryCode()

获取用户国家代码。

```go
countryCode := client.GetCountryCode()
```

## 高级用法

### 使用代理

```go
config := &Steam.Config{
    Proxy:   "127.0.0.1:8080", // SOCKS5或HTTP代理
    Timeout: 30 * time.Second,
}
client, err := Steam.NewClient(config)
```

### 错误处理

该库提供了详细的错误信息，建议根据不同的错误类型进行处理：

```go
userInfo, err := client.Login(credentials)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "密码错误"):
        fmt.Println("用户名或密码错误")
    case strings.Contains(err.Error(), "需要手机验证码"):
        fmt.Println("需要Steam Guard验证")
    case strings.Contains(err.Error(), "请求过于频繁"):
        fmt.Println("请求过于频繁，请稍后重试")
    default:
        fmt.Printf("登录失败: %v\n", err)
    }
    return
}
```

### Steam Guard设置

要使用Steam Guard功能，你需要获取共享密钥：

1. 在Steam移动应用中启用Steam Guard
2. 提取共享密钥（通常为base64编码的字符串）
3. 在登录时提供该密钥

```go
credentials := &Steam.LoginCredentials{
    Username:     "your_username",
    Password:     "your_password",
    SharedSecret: "abcdefghijk123456789==", // base64编码的密钥
}
```

## 注意事项

1. **频率限制**: Steam对API请求有频率限制，建议在请求之间添加适当的延迟
2. **账户安全**: 请妥善保管你的登录凭据，特别是SharedSecret
3. **合规使用**: 请遵守Steam的使用条款，仅用于个人学习和合规目的
4. **错误重试**: 库内置了重试机制，但仍建议在应用层面进行适当的错误处理

## 许可证

本项目采用MIT许可证，详见LICENSE文件。

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 支持

如果你在使用过程中遇到问题，可以：

1. 查看examples目录中的示例代码
2. 提交Issue描述你的问题
3. 查阅Steam官方API文档