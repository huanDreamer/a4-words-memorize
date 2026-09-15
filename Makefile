# A4 背单词法 —— 常用命令集合
#
# 最常用：make help / make run / make test / make check
# 本地覆盖变量：复制下面任意变量到 local.mk（已被 include，且不会进 git）

# ===== 可覆盖变量 =====
GO          ?= go
BIN         ?= words
CMD_PKG     ?= ./cmd
IMPORT_PKG  ?= ./cmd/import

DATA_DIR    ?= data
SOURCE_DIR  ?= source
ADDR        ?= :8900
PORT        ?= 8900
HOST        ?= 127.0.0.1

IMAGE       ?= words
TAG         ?= latest

CGO_ENABLED ?= 0
GOFLAGS     ?=

# 生成数据集时额外传给 cmd/import 的参数，例如：make dataset IMPORT_ARGS="-sample 200"
IMPORT_ARGS ?= -plans

SHELL := /bin/sh
.DEFAULT_GOAL := help

# ===== 帮助 =====
.PHONY: help
help: ## 显示所有可用命令
	@echo "A4 背单词法 —— 可用命令："
	@echo ""
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ===== 构建 =====
.PHONY: build
build: ## 编译当前平台可执行文件到 ./words
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) -o $(BIN) $(CMD_PKG)
	@echo "已生成 ./$(BIN)"

.PHONY: build-linux
build-linux: ## 交叉编译 Linux/amd64 可执行文件（Docker 镜像使用）
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BIN) $(CMD_PKG)
	@echo "已生成 linux/amd64 版 ./$(BIN)"

.PHONY: install
install: ## 安装可执行文件到 $GOBIN（默认 $GOPATH/bin）
	CGO_ENABLED=$(CGO_ENABLED) $(GO) install $(CMD_PKG)
	@echo "已安装到 $$($(GO) env GOBIN)"

# ===== 运行 =====
.PHONY: run
run: build ## 编译并启动服务（读取 ./data 数据目录）
	./$(BIN) -addr $(ADDR) -data $(DATA_DIR)

.PHONY: dev
dev: ## 开发模式直接 go run 启动（gin debug、热改代码后重跑）
	$(GO) run $(CMD_PKG) -addr $(ADDR) -data $(DATA_DIR)

.PHONY: release
release: build ## 以 release(生产) 模式启动
	GIN_MODE=release ./$(BIN) -addr $(ADDR) -data $(DATA_DIR)

# ===== 数据 =====
.PHONY: dataset
dataset: ## 生成测试数据集：导入 source/ 词典到 ./data（可加 IMPORT_ARGS="-sample 200"）
	$(GO) run $(IMPORT_PKG) -data $(DATA_DIR) -source $(SOURCE_DIR) $(IMPORT_ARGS)

.PHONY: reset-data
reset-data: ## 清空并重新生成数据集（慎用：会删除 ./data）
	rm -rf $(DATA_DIR)
	$(MAKE) dataset

.PHONY: show-data
show-data: ## 查看当前数据集的规模概况
	@for f in $(DATA_DIR)/*.json; do \
		[ -f "$$f" ] || continue; \
		printf '  %-24s %10s bytes\n' "$$f" "$$(wc -c < "$$f" | tr -d ' ')"; \
	done

# ===== 测试与检查 =====
.PHONY: test
test: ## 运行全部单元测试
	$(GO) test $(GOFLAGS) ./...

.PHONY: test-v
test-v: ## 运行全部单元测试（详细输出）
	$(GO) test $(GOFLAGS) -v ./...

.PHONY: race
race: ## 运行单元测试并开启竞态检测
	$(GO) test $(GOFLAGS) -race ./...

.PHONY: cover
cover: ## 运行测试并生成覆盖率报告 coverage.html
	$(GO) test $(GOFLAGS) -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -func=coverage.out | tail -n 1
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "已生成 coverage.html"

.PHONY: fmt
fmt: ## 格式化代码（gofmt -s -w）
	$(GO) fmt ./...
	@gofmt -s -w .

.PHONY: fmt-check
fmt-check: ## 检查代码格式（CI 用，不修改文件）
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then echo "以下文件未通过 gofmt："; echo "$$out"; exit 1; fi
	@echo "gofmt 检查通过"

.PHONY: vet
vet: ## 运行 go vet 静态检查
	$(GO) vet ./...

.PHONY: lint
lint: vet ## 运行静态检查（有 staticcheck 时一并运行）
	@if command -v staticcheck >/dev/null 2>&1; then staticcheck ./...; \
	else echo "未安装 staticcheck，已跳过（go install honnef.co/go/tools/cmd/staticcheck@latest）"; fi

.PHONY: tidy
tidy: ## 整理依赖（go mod tidy）
	$(GO) mod tidy

.PHONY: check
check: fmt-check vet test ## 提交前检查：格式 + 静态检查 + 测试

# ===== Docker =====
.PHONY: docker-build
docker-build: build-linux ## 构建 Docker 镜像 words:latest
	docker build -t $(IMAGE):$(TAG) .

.PHONY: docker-run
docker-run: ## 启动容器（挂载 ./data 持久化数据）
	docker run --rm -it \
		-p $(PORT):8900 \
		-v "$(CURDIR)/$(DATA_DIR):/apps/data" \
		-e WORDS_DATA_DIR=/apps/data \
		$(IMAGE):$(TAG)

.PHONY: docker-push
docker-push: ## 推送镜像到仓库 words:latest
	docker push $(IMAGE):$(TAG)

# ===== 冒烟测试 =====
.PHONY: smoke
smoke: ## 对已启动的服务做快速冒烟测试（需先 make run）
	@printf "对 http://%s:%s 冒烟...\n" "$(HOST)" "$(PORT)"
	@code=$$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 http://$(HOST):$(PORT)/); \
	if [ "$$code" != "200" ]; then \
		echo "  失败：期望 200，实际 $$code（服务是否已启动？先执行 make run）"; exit 1; \
	fi
	@printf "  GET  /                -> 200\n"
	@printf "  GET  /study?bookId=CET4luan_1 -> "
	@curl -s -o /dev/null -w '%{http_code}\n' --max-time 5 "http://$(HOST):$(PORT)/study?bookId=CET4luan_1"
	@echo "冒烟通过"

# ===== 清理 =====
.PHONY: clean
clean: ## 清理编译产物与覆盖率文件（保留数据）
	rm -f $(BIN) coverage.out coverage.html
	$(GO) clean -cache -testcache
	@echo "已清理编译产物"

.PHONY: clean-all
clean-all: clean ## 清理编译产物 **并删除** ./data 数据
	rm -rf $(DATA_DIR)
	@echo "已删除 $(DATA_DIR)/"

# 本地变量覆盖（可选，不纳入版本库）
-include local.mk
