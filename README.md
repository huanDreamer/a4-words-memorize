## A4背单词法

一个用 Go 编写的背单词 Web 应用：以「单词书」为单位，每次生成 25 个新单词的学习计划，并记录学习进度。

### 技术栈

- Go 1.19+ / Gin
- **本地 JSON 文件存储**（无需 MongoDB 等外部数据库）
- 数据来源词典：https://github.com/huanDreamer/dict

### 存储说明

数据保存在 `data/` 目录下，每个「集合」对应一个 JSON 文件：

| 文件 | 内容 |
| --- | --- |
| `data/book.json` | 单词书列表 |
| `data/words.json` | 单词及释义 |
| `data/study_plan.json` | 学习计划（最近学习记录） |

写入采用「先写临时文件再原子 rename」，并保留一份 `.bak` 备份；读操作无需加锁。
数据目录可通过 `-data` 参数或环境变量 `WORDS_DATA_DIR` 覆盖（默认 `data`）。

仓库内已附带一份**开箱可用的测试数据集**（四级词汇 1162 词 + 雅思词汇 3427 词，共 4589 词，
并含 10 条示例学习记录），clone 后无需导入即可直接启动。
运行时产生的 `.bak` / `.tmp` 已在 `.gitignore` 中忽略，不会污染工作区。
若词典源有更新，或想要一份全新的数据，重新生成即可（会覆盖现有数据）：

```bash
make reset-data                          # 清空并重建整套数据集
make dataset IMPORT_ARGS="-sample 200"   # 或只造一份小数据集
```

### 快速开始

> 项目提供了 Makefile 封装常用命令，执行 `make help` 可查看全部目标。

1. 启动服务（数据集已随仓库提供，无需额外准备）：

   ```bash
   make run       # 编译并启动，默认监听 :8900
   make dev       # go run 直接启动（开发模式，gin debug 日志）
   make release   # 以 release 模式启动
   ```

2. 打开 http://localhost:8900

### 常用命令（Makefile）

| 命令 | 说明 |
| --- | --- |
| `make help` | 列出全部可用目标 |
| `make build` / `make build-linux` | 编译本机 / Linux-amd64 可执行文件 |
| `make run` / `make dev` / `make release` | 启动服务（生产 / 开发 / release 模式） |
| `make dataset` / `make reset-data` / `make show-data` | 生成 / 重建 / 查看数据集 |
| `make test` / `make race` / `make cover` | 单元测试 / 竞态检测 / 覆盖率报告 |
| `make fmt` / `make fmt-check` / `make vet` / `make lint` | 格式化与静态检查 |
| `make check` | 提交前检查：格式 + vet + 测试 |
| `make smoke` | 对已启动的服务做一次冒烟测试 |
| `make docker-build` / `make docker-run` / `make docker-push` | 构建 / 运行 / 推送镜像 |
| `make clean` / `make clean-all` | 清理编译产物 / 连同数据一起清理 |

可用变量（可命令行覆盖，如 `make run ADDR=:9000`）：`ADDR`、`DATA_DIR`、`IMAGE`、`TAG`、`PORT`、`IMPORT_ARGS` 等。
若需要固定本地配置，可新建 `local.mk`（已被 include，且不会提交到 git）。

### 常用参数

| 参数 | 说明 | 默认值 |
| --- | --- | --- |
| `-addr` | HTTP 监听地址（或 `WORDS_ADDR`） | `:8900` |
| `-data` | 数据目录（或 `WORDS_DATA_DIR`） | `data` |
| `-release` | production 模式（或 `GIN_MODE=release`） | false |

### cmd/import 参数

| 参数 | 说明 | 默认值 |
| --- | --- | --- |
| `-data` | 数据目录 | `data` |
| `-source` | 词典源目录（`*.json`，一行一个单词） | `source` |
| `-sample` | 每本书最多导入多少个单词，0 表示全部 | 0 |
| `-plans` | 是否额外生成示例学习计划 | false |
| `-plan-count` | 每本书生成的示例学习计划数量 | 5 |

### 测试

```bash
make test            # 等价于 go test ./...
make race            # 开启竞态检测
make cover           # 生成 coverage.html
```

### Docker 部署

```bash
make docker-build                      # 交叉编译 + 构建镜像 words:latest
make docker-run                        # 运行容器，并挂载 ./data 持久化
# 或手动：
docker build -t words:v1.0 .
docker run -p 8900:8900 -v $(pwd)/data:/apps/data words:v1.0
```

### 系统功能介绍

1. 单词本列表
2. 最近学习记录
3. 单词学习页面
