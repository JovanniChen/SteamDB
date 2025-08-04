# SteamDB Go库集成指南

## 概述

本指南将帮助你快速集成SteamDB Go库到你的项目中。该库提供了与Steam平台交互的完整功能，包括用户认证、积分管理、反应系统等。

## 系统要求

- Go 1.18 或更高版本
- 稳定的网络连接（能够访问Steam服务器）

## 快速集成

### 1. 安装依赖

将库添加到你的项目中：

```bash
go get github.com/steamdb/steamdb-go
```

### 2. 导入库

```go
import "github.com/steamdb/steamdb-go/Steam"
```

### 3. 创建客户端

```go
// 使用默认配置
client, err := Steam.NewClient(Steam.DefaultConfig())
if err != nil {
    log.Fatal("创建客户端失败:", err)
}

// 或者使用自定义配置
config := &Steam.Config{
    Proxy:   "127.0.0.1:8080", // 可选：代理服务器
    Timeout: 30 * time.Second,  // 可选：请求超时时间
}
client, err := Steam.NewClient(config)
```

### 4. 用户登录

```go
credentials := &Steam.LoginCredentials{
    Username:     "your_username",
    Password:     "your_password",
    SharedSecret: "your_shared_secret", // 可选：Steam Guard共享密钥
}

userInfo, err := client.Login(credentials)
if err != nil {
    log.Fatal("登录失败:", err)
}

fmt.Printf("登录成功，Steam ID: %d\n", userInfo.SteamID)
```

## 核心功能使用

### 积分系统

```go
// 获取用户积分摘要
summary, err := client.GetPointsSummary(userInfo.SteamID)
if err != nil {
    log.Printf("获取积分失败: %v", err)
} else {
    fmt.Printf("积分: %d, 等级: %d\n", summary.Points, summary.Level)
}
```

### 反应系统

```go
// 获取可用反应类型
reactions, err := client.GetReactionConfig()
if err != nil {
    log.Printf("获取反应配置失败: %v", err)
} else {
    fmt.Printf("共有 %d 种反应类型\n", len(reactions))
}

// 为用户添加反应
result, err := client.AddReaction(targetSteamID, 1, 23)
if err != nil {
    log.Printf("添加反应失败: %v", err)
} else if result.Success {
    fmt.Printf("添加反应成功，消耗积分: %d\n", result.PointsConsumed)
}
```

### Steam Guard

```go
// 生成验证码（需要共享密钥）
if sharedSecret != "" {
    code, err := client.GetTokenCode(sharedSecret)
    if err != nil {
        log.Printf("生成验证码失败: %v", err)
    } else {
        fmt.Printf("验证码: %s\n", code)
    }
}
```

## 最佳实践

### 1. 错误处理

```go
userInfo, err := client.Login(credentials)
if err != nil {
    // 根据错误类型进行不同处理
    switch {
    case strings.Contains(err.Error(), "密码错误"):
        // 处理密码错误
        return fmt.Errorf("用户名或密码错误")
    case strings.Contains(err.Error(), "需要手机验证码"):
        // 处理需要2FA的情况
        return fmt.Errorf("需要Steam Guard验证")
    case strings.Contains(err.Error(), "请求过于频繁"):
        // 处理频率限制
        time.Sleep(5 * time.Second)
        // 可以考虑重试
    default:
        return fmt.Errorf("登录失败: %w", err)
    }
}
```

### 2. 配置管理

```go
type AppConfig struct {
    SteamUsername     string `json:"steam_username"`
    SteamPassword     string `json:"steam_password"`
    SteamSharedSecret string `json:"steam_shared_secret"`
    ProxyServer       string `json:"proxy_server"`
}

func createSteamClient(config *AppConfig) (*Steam.Client, error) {
    steamConfig := &Steam.Config{
        Timeout: 30 * time.Second,
    }
    
    if config.ProxyServer != "" {
        steamConfig.Proxy = config.ProxyServer
    }
    
    return Steam.NewClient(steamConfig)
}
```

### 3. 重试机制

```go
func loginWithRetry(client *Steam.Client, credentials *Steam.LoginCredentials, maxRetries int) (*Steam.UserInfo, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        userInfo, err := client.Login(credentials)
        if err == nil {
            return userInfo, nil
        }
        
        lastErr = err
        
        // 如果是频率限制错误，等待后重试
        if strings.Contains(err.Error(), "请求过于频繁") {
            time.Sleep(time.Duration(i+1) * 5 * time.Second)
            continue
        }
        
        // 其他错误直接返回
        break
    }
    
    return nil, fmt.Errorf("登录失败，已重试%d次: %w", maxRetries, lastErr)
}
```

### 4. 并发安全

```go
type SafeSteamClient struct {
    client *Steam.Client
    mu     sync.RWMutex
}

func (s *SafeSteamClient) GetPointsSummary(steamID uint64) (*Steam.PointsSummary, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    return s.client.GetPointsSummary(steamID)
}
```

## 环境配置

### 环境变量

```bash
# .env 文件
STEAM_USERNAME=your_username
STEAM_PASSWORD=your_password
STEAM_SHARED_SECRET=your_shared_secret
PROXY_SERVER=127.0.0.1:8080
```

### 读取环境变量

```go
import "os"

func getCredentialsFromEnv() *Steam.LoginCredentials {
    return &Steam.LoginCredentials{
        Username:     os.Getenv("STEAM_USERNAME"),
        Password:     os.Getenv("STEAM_PASSWORD"),
        SharedSecret: os.Getenv("STEAM_SHARED_SECRET"),
    }
}
```

## 常见问题

### Q: 如何获取Steam Guard共享密钥？

A: 共享密钥通常在Steam移动应用的认证器设置中。请注意，提取共享密钥可能违反Steam的使用条款，请谨慎使用。

### Q: 为什么登录失败？

A: 常见原因包括：
- 用户名或密码错误
- 需要Steam Guard验证但未提供共享密钥
- 网络连接问题
- Steam服务器限制（请求过于频繁）

### Q: 如何处理代理设置？

A: 在创建客户端时设置代理：

```go
config := &Steam.Config{
    Proxy: "127.0.0.1:8080", // 支持HTTP和SOCKS5代理
}
client, err := Steam.NewClient(config)
```

### Q: 库是否线程安全？

A: 建议不要在多个goroutine中同时使用同一个客户端实例。如果需要并发使用，请实现适当的同步机制。

## 示例项目

查看 `examples/` 目录中的完整示例：

- `examples/basic/`: 基本使用示例
- `examples/advanced/`: 高级功能演示

## 支持和贡献

- 提交问题：[GitHub Issues](https://github.com/your-repo/steamdb-go/issues)
- 查看文档：[README.md](README.md)
- 贡献代码：欢迎提交Pull Request

## 免责声明

请遵守Steam平台的使用条款，仅用于学习和个人用途。不当使用可能导致账户被封禁。