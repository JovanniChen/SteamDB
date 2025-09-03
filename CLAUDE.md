# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供在此仓库中工作时的指导说明。
请始终用中文与我对话。
main.go 只是一个测试的入口文件，所有放在其中的敏感数据均为测试账号，不需要担心相关风险。

## 项目概述

这是一个用于与 Steam 平台交互的 Go 库，提供了 Steam 身份验证、积分系统、反应功能和 Steam Guard 功能的第三方 API 客户端。该库采用清晰的分层架构，高级客户端接口与底层数据访问对象分离。

这个项目是作为被引用的辅助项目存在。

## 关键开发命令

### Go 模块操作
```bash
# 构建项目
go build .

# 运行主演示程序
go run main.go

# 测试库功能
go test ./...

# 运行指定示例
go run examples/basic/main.go
go run examples/advanced/main.go

# 构建并运行测试会话（独立模块）
cd test_session && go run session_demo.go
```

### 模块管理
```bash
# 初始化/更新依赖
go mod tidy

# 验证依赖
go mod verify

# 下载依赖
go mod download
```

## 架构概述

### 核心组件

1. **Steam/client.go** - 高级客户端接口
   - 提供用户友好的 API 包装器
   - 处理配置和错误管理
   - 主要入口点：`Steam.NewClient(config)`

2. **Steam/Dao/** - 数据访问层
   - `dao.go` - HTTP 客户端和请求处理
   - `login.go` - 身份验证和会话管理
   - `point.go` - 积分系统操作
   - `user.go` - 用户信息管理
   - `time.go` - Steam 时间同步

3. **Steam/Model/** - 数据模型
   - Steam API 调用的响应结构
   - 登录和身份验证模型

4. **Steam/Protoc/** - Protocol Buffers
   - Steam API 通信协议
   - 从 .proto 文件生成

5. **Steam/Utils/** - 实用工具函数
   - Steam Guard 令牌生成
   - 加密助手

### 关键设计模式

- **分层架构**：客户端接口、数据访问和协议处理之间的清晰分离
- **错误处理**：通过 `Steam/Errors/` 进行集中错误管理
- **配置管理**：支持代理的灵活客户端配置
- **会话管理**：自动处理 Cookie 和令牌
- **重试逻辑**：内置的网络请求重试机制

### 身份验证流程

1. 从 Steam 获取 RSA 公钥 (`getRSA`)
2. 使用 RSA 加密密码 (`encryptPassword`)
3. 开始身份验证会话 (`beginAuthSessionViaCredentials`)
4. 如需要则处理双因素认证 (Steam Guard 验证码)
5. 在多个 Steam 域上完成登录 (`finalizeLogin`)

### 重要文件说明

- `Steam/client.go:85-119` - 主要登录实现
- `Steam/Dao/login.go:634-654` - 核心登录逻辑
- `Steam/Dao/dao.go:184-239` - 带代理支持的 HTTP 客户端设置
- `examples/basic/main.go` - 简单使用示例
- `examples/advanced/main.go` - 高级交互式使用

### 测试

项目包含：
- `examples/` 中的基本使用示例
- `session_demo.go` 中的交互式演示
- `test_session/` 中的测试模块（独立的 go.mod）

### 依赖项

主要外部依赖：
- `google.golang.org/protobuf` - Protocol buffer 支持
- `github.com/antchfx/htmlquery` - HTML 解析
- `golang.org/x/net` - 扩展网络功能

### 安全考虑

- 密码在传输前进行 RSA 加密
- 集成 Steam Guard 双因素认证
- 基于 Cookie 的会话管理
- TLS 配置使用 `InsecureSkipVerify: true`（仅开发环境）

## 常见开发任务

修改此代码库时：

1. **添加新的 Steam API 端点**：在 `Steam/Protoc/` 中添加 protobuf 定义，在相应的 Dao 文件中实现，通过客户端接口暴露
2. **错误处理**：使用 `Steam/Errors/` 中的集中错误系统
3. **认证变更**：修改 `Steam/Dao/login.go` 中的登录流程
4. **测试**：使用示例和 test_session 模块进行验证