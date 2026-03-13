# SteamDB Makefile (Minimal)

PROJECT_NAME := SteamDB
GO := go
BUILD_DIR := build

.PHONY: help run test fmt vet build clean mod-tidy

help:
	@echo "SteamDB Makefile (Minimal)"
	@echo ""
	@echo "可用目标:"
	@echo "  run        - 运行程序 (go run .)"
	@echo "  build      - 构建当前平台二进制到 build/$(PROJECT_NAME)"
	@echo "  test       - 运行测试 (go test ./...)"
	@echo "  fmt        - 格式化代码 (go fmt ./...)"
	@echo "  vet        - 代码检查 (go vet ./...)"
	@echo "  mod-tidy   - 整理并校验依赖 (go mod tidy && go mod verify)"
	@echo "  clean      - 清理 build 目录"

run:
	$(GO) run .

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -trimpath -o $(BUILD_DIR)/$(PROJECT_NAME) .

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

mod-tidy:
	$(GO) mod tidy
	$(GO) mod verify

clean:
	rm -rf $(BUILD_DIR)
