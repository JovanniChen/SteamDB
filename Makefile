# SteamDB Makefile
# 支持多平台编译的Makefile

# 项目信息
PROJECT_NAME := SteamDB
VERSION := 1.0.0
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go相关变量
GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
CGO_ENABLED := 0

# 构建目录
BUILD_DIR := build
DIST_DIR := dist

# 编译标志
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -s -w"
BUILD_FLAGS := -trimpath

# 默认目标
.PHONY: all
all: clean build-all

# 帮助信息
.PHONY: help
help:
	@echo "SteamDB 构建系统"
	@echo ""
	@echo "可用目标:"
	@echo "  build          - 构建当前平台版本"
	@echo "  build-all      - 构建所有平台版本"
	@echo "  build-windows  - 构建Windows版本 (amd64)"
	@echo "  build-mac      - 构建macOS版本 (amd64, arm64)"
	@echo "  build-linux    - 构建Linux版本 (amd64, arm64)"
	@echo "  clean          - 清理构建文件"
	@echo "  test           - 运行测试"
	@echo "  fmt            - 格式化代码"
	@echo "  vet            - 代码检查"
	@echo "  mod-tidy       - 整理依赖"
	@echo "  run            - 运行程序"
	@echo "  help           - 显示此帮助信息"
	@echo ""
	@echo "当前平台: $(GOOS)/$(GOARCH)"

# 构建当前平台版本
.PHONY: build
build:
	@echo "构建 $(GOOS)/$(GOARCH) 版本..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME) .

# 构建所有平台版本
.PHONY: build-all
build-all: build-windows build-mac build-linux
	@echo "所有平台构建完成!"

# 构建Windows版本
.PHONY: build-windows
build-windows:
	@echo "构建Windows版本..."
	@mkdir -p $(DIST_DIR)/windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/windows/$(PROJECT_NAME).exe .
	@echo "Windows版本构建完成: $(DIST_DIR)/windows/$(PROJECT_NAME).exe"

# 构建macOS版本
.PHONY: build-mac
build-mac:
	@echo "构建macOS版本..."
	@mkdir -p $(DIST_DIR)/macos
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/macos/$(PROJECT_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/macos/$(PROJECT_NAME)-darwin-arm64 .
	@echo "macOS版本构建完成:"
	@echo "  - $(DIST_DIR)/macos/$(PROJECT_NAME)-darwin-amd64"
	@echo "  - $(DIST_DIR)/macos/$(PROJECT_NAME)-darwin-arm64"

# 构建Linux版本
.PHONY: build-linux
build-linux:
	@echo "构建Linux版本..."
	@mkdir -p $(DIST_DIR)/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/linux/$(PROJECT_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(DIST_DIR)/linux/$(PROJECT_NAME)-linux-arm64 .
	@echo "Linux版本构建完成:"
	@echo "  - $(DIST_DIR)/linux/$(PROJECT_NAME)-linux-amd64"
	@echo "  - $(DIST_DIR)/linux/$(PROJECT_NAME)-linux-arm64"

# 运行程序
.PHONY: run
run:
	@echo "运行程序..."
	$(GO) run .

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	$(GO) test -v ./...

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	$(GO) fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "代码检查..."
	$(GO) vet ./...

# 整理依赖
.PHONY: mod-tidy
mod-tidy:
	@echo "整理依赖..."
	$(GO) mod tidy
	$(GO) mod verify

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@echo "清理完成!"

# 安装依赖
.PHONY: deps
deps:
	@echo "下载依赖..."
	$(GO) mod download

# 显示项目信息
.PHONY: info
info:
	@echo "项目信息:"
	@echo "  名称: $(PROJECT_NAME)"
	@echo "  版本: $(VERSION)"
	@echo "  构建时间: $(BUILD_TIME)"
	@echo "  Git提交: $(GIT_COMMIT)"
	@echo "  Go版本: $(shell go version)"
	@echo "  当前平台: $(GOOS)/$(GOARCH)"

# 创建发布包
.PHONY: package
package: build-all
	@echo "创建发布包..."
	@mkdir -p $(DIST_DIR)/packages
	@cd $(DIST_DIR) && \
		tar -czf packages/$(PROJECT_NAME)-windows-amd64.tar.gz windows/ && \
		tar -czf packages/$(PROJECT_NAME)-macos.tar.gz macos/ && \
		tar -czf packages/$(PROJECT_NAME)-linux.tar.gz linux/
	@echo "发布包创建完成:"
	@ls -la $(DIST_DIR)/packages/

# 开发环境设置
.PHONY: dev-setup
dev-setup: deps mod-tidy fmt vet
	@echo "开发环境设置完成!"

# 持续集成目标
.PHONY: ci
ci: fmt vet test build-all
	@echo "CI构建完成!"
