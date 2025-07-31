# Claude Code 提问建议指南

## 概述
本文档提供了与Claude Code交互的最佳实践和建议，帮助您更有效地使用AI助手进行软件开发工作。

## 提问方式建议

### 1. 明确具体的任务
**好的提问方式：**
```
请帮我在Steam/Model/目录下添加一个Point.go文件，包含积分管理的结构体
```

**避免的提问方式：**
```
帮我做点什么
```

### 2. 提供上下文信息
**好的提问方式：**
```
我的Go项目中有一个Steam模块，现在需要在Steam/Dao/point.go中添加查询用户积分的方法，
请参考Steam/Dao/user.go中的现有模式
```

**相关命令：**
```bash
# 查看项目结构
find . -name "*.go" -type f | head -20

# 查看特定文件内容
cat Steam/Dao/user.go

# 搜索相关代码模式
grep -r "func.*Get" Steam/Dao/
```

### 3. 指定期望的输出格式
**好的提问方式：**
```
请帮我重构main.go文件，并将修改保存到文件中
```

**相关命令：**
```bash
# 备份原文件（建议操作）
cp main.go main.go.backup

# 查看当前文件状态
git status

# 查看文件差异
git diff main.go
```

### 4. 分步骤的复杂任务
**好的提问方式：**
```
请帮我：
1. 分析Steam/Protoc/point.proto文件
2. 根据proto定义生成对应的Go结构体
3. 在Steam/Model/目录创建Point.go文件
4. 添加相应的数据库操作方法
```

**相关命令：**
```bash
# 查看proto文件
cat Steam/Protoc/point.proto

# 生成protobuf代码
protoc --go_out=. Steam/Protoc/point.proto

# 检查生成的文件
ls -la Steam/Model/

# 运行测试验证
go test ./Steam/Model/
```

### 5. 代码审查和优化
**好的提问方式：**
```
请审查Steam/Utils/utils.go文件，重点关注性能和安全性，
并提供具体的改进建议
```

**相关命令：**
```bash
# 代码静态分析
go vet ./Steam/Utils/

# 运行格式化
go fmt ./Steam/Utils/

# 性能分析
go test -bench=. ./Steam/Utils/

# 安全检查（需要安装gosec）
gosec ./Steam/Utils/
```

### 6. 错误诊断和修复
**好的提问方式：**
```
我运行`go build`时出现编译错误，错误信息如下：
[粘贴具体错误信息]
请帮我分析并修复这些问题
```

**相关命令：**
```bash
# 详细编译信息
go build -v ./...

# 检查依赖
go mod tidy
go mod verify

# 查看Go版本
go version

# 清理模块缓存
go clean -modcache
```

### 7. 测试相关任务
**好的提问方式：**
```
请为Steam/Dao/point.go中的GetUserPoints方法编写单元测试，
使用项目中已有的测试模式
```

**相关命令：**
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./Steam/Dao/

# 生成测试覆盖率报告
go test -cover ./Steam/Dao/

# 详细测试输出
go test -v ./Steam/Dao/

# 运行基准测试
go test -bench=. ./Steam/Dao/
```

### 8. Git操作相关
**好的提问方式：**
```
请帮我提交当前的代码变更，包括新增的Point.go文件，
提交信息为"add point management features"
```

**相关命令：**
```bash
# 查看当前状态
git status

# 查看变更详情
git diff

# 添加特定文件
git add Steam/Model/Point.go

# 添加所有变更
git add .

# 提交变更
git commit -m "add point management features"

# 查看提交历史
git log --oneline -5
```

## 高效交互技巧

### 1. 使用文件路径引用
```
请修改Steam/Config.go:25行的数据库连接配置
```

### 2. 指定编程语言和框架
```
这是一个Go项目，使用GORM作为ORM，请帮我...
```

### 3. 提供期望的代码风格
```
请遵循项目现有的代码风格，参考Steam/Dao/user.go的写法
```

### 4. 请求解释和文档
```
请解释Steam/Utils/phoneToken.go中的令牌生成逻辑，
并添加必要的注释
```

## 项目特定的有用命令

### Go项目管理
```bash
# 初始化/更新依赖
go mod init
go mod tidy

# 构建项目
go build -o steamdb main.go

# 运行项目
go run main.go

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o steamdb-linux main.go
```

### 代码质量检查
```bash
# 格式化代码
go fmt ./...

# 代码检查
go vet ./...

# 查找未使用的代码
go list -f '{{.Dir}} {{.GoFiles}}' ./... | xargs grep -l "^func"
```

### 性能分析
```bash
# CPU性能分析
go test -cpuprofile=cpu.prof -bench=.

# 内存分析
go test -memprofile=mem.prof -bench=.

# 查看性能报告
go tool pprof cpu.prof
```

## 注意事项

1. **明确性** - 越具体的问题，越能得到准确的答案
2. **上下文** - 提供足够的背景信息
3. **分步骤** - 复杂任务分解为小步骤
4. **验证** - 要求运行测试或检查来验证结果
5. **安全性** - 避免在代码中暴露敏感信息

## 示例对话流程

```
用户：请帮我在Steam/Model/目录下创建Point.go文件，定义积分相关的结构体

Claude：我来帮您创建Point.go文件。首先让我查看一下现有的模型结构...
[执行相关操作]

用户：很好，现在请在Steam/Dao/目录下添加对应的数据库操作方法

Claude：现在我来添加数据库操作方法，我会参考现有的user.go文件的模式...
[执行相关操作]

用户：请为这些方法编写单元测试

Claude：我来为这些方法编写单元测试...
[执行相关操作]
```

这样的交互方式能够确保每一步都清晰明确，并且能够得到高质量的代码实现。