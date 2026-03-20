**[English](./README.md) | 中文**

# Gradmotion CLI

跨平台命令行工具，用于管理 [Gradmotion](https://spaces.gradmotion.com) 平台上的训练任务与项目。默认输出结构化 JSON，天然适配 AI Agent 集成。

## 特性

- **跨平台** — 单一二进制，支持 macOS / Linux / Windows（amd64 & arm64）
- **安全认证** — API Key 优先存入系统 Keychain（macOS Keychain / Windows Credential Manager / Linux Secret Service），自动降级至本地配置
- **多环境 Profile** — 通过 Profile 在开发、测试、生产等多环境间快速切换
- **完整任务生命周期** — 创建、编辑、运行、停止、删除、日志追踪、超参管理、批量操作一站搞定
- **Agent 友好** — 默认 JSON 输出 + JSONL 结构化日志，支持 Cursor Agent Skill 集成
- **可观测** — 每请求自动注入 Trace ID，支持 `--debug` 调试和 `--log-file` 日志落盘
- **健壮请求** — 内置指数退避重试、可配超时与并发控制

## 安装

### npm（推荐）

```bash
npm install -g @limxdynamics/gm-cli
```

> 需要 Node.js >= 16，安装时自动匹配当前平台架构。

### 从 GitHub Release 下载

前往 [Releases](https://github.com/limxdynamics/gradmotion-cli/releases) 下载对应平台的压缩包，解压后将 `gm` 放入 `PATH`：

```bash
# macOS / Linux
sudo install -m 0755 gm /usr/local/bin/gm
```

### 从源码编译

```bash
git clone https://github.com/limxdynamics/gradmotion-cli.git
cd gradmotion-cli
make build      # 编译当前平台 → ./gm
make install    # 安装到 /usr/local/bin（需 sudo）
```

### 验证

```bash
gm --version
```

## 快速开始

```bash
# 1. 设置服务地址
gm config set base_url "https://spaces.gradmotion.com/prod-api"

# 2. 登录（API Key 存入系统 Keychain）
gm auth login --api-key "gm_sk_your_key"

# 3. 验证
gm auth status          # 检查本地认证状态
gm auth whoami          # 验证服务端连接

# 4. 开始使用
gm project list
gm task list
```

> API Key 在 Gradmotion 平台 -> 左下角头像 -> **API Key 管理** 中创建，仅在创建时完整显示一次，请妥善保存。

## 命令概览

```
gm
├── auth                        # 认证管理
│   ├── login                   # 保存 API Key
│   ├── logout                  # 清除 API Key
│   ├── status                  # 查看本地认证状态
│   └── whoami                  # 查询服务端用户信息
├── config                      # 配置管理
│   ├── set <key> <value>       # 设置配置项
│   ├── get <key>               # 读取配置项
│   └── profile                 # 多环境 Profile
│       ├── list                # 列出所有 Profile
│       ├── set <name>          # 创建/更新 Profile
│       └── use <name>          # 切换当前 Profile
├── project                     # 项目管理
│   ├── list                    # 项目列表
│   ├── create                  # 创建项目
│   ├── info                    # 项目详情
│   ├── edit                    # 编辑项目
│   └── delete                  # 删除项目
└── task                        # 任务管理
    ├── create / edit / copy    # 创建、编辑、复制任务
    ├── list / info             # 列表与详情
    ├── run / stop / delete     # 运行、停止、删除
    ├── logs                    # 查看/追踪任务日志
    ├── resource list           # 算力资源列表
    ├── image                   # 镜像管理（official / personal / versions）
    ├── storage list            # 个人存储列表
    ├── data                    # 训练数据（keys / get / download）
    ├── hp get                  # 超参读取
    ├── env get                 # 运行环境
    ├── params                  # 超参管理（submit / update）
    ├── tag                     # 标签管理（update / get / list）
    └── batch                   # 批量操作（stop / delete）
```

使用 `gm --help` 或 `gm <command> --help` 查看详细用法。

## 常用示例

### 任务管理

```bash
# 查看任务列表
gm task list --page 1 --limit 50

# 查看任务详情
gm task info --task-id "task_xxx"

# 创建任务
gm task create --file ./create.json

# 运行任务
gm task run --task-id "task_xxx"

# 查看实时日志
gm task logs --task-id "task_xxx" --follow --interval 2s --timeout 5m

# 停止任务（需确认）
gm task stop --task-id "task_xxx"
```

### 项目管理

```bash
gm project list --page 1 --limit 50
gm project create --file ./project-create.json
gm project info --project-id "proj_xxx"
```

### 多环境切换

```bash
# 创建 dev Profile
gm config profile set dev --base-url "https://dev.gradmotion.com/prod-api" --timeout 30s

# 切换到 dev
gm config profile use dev

# 临时使用某 Profile（不切换默认）
gm --profile dev task list
```

### 输出控制

```bash
gm task list --human       # 人类可读表格
gm task list --quiet       # 仅关键字段
gm task list --debug       # 开启调试日志
gm task logs --task-id "task_xxx" --raw --no-request-log   # 纯净日志流，适合管道
```

## 配置

### 优先级

```
CLI flags  >  环境变量  >  配置文件
```

### 环境变量

| 变量 | 说明 |
|------|------|
| `GM_BASE_URL` | 服务地址 |
| `GM_API_KEY` | API Key（临时覆盖，不落盘） |
| `GM_TIMEOUT` | 请求超时（如 `30s`） |
| `GM_RETRY` | 重试次数 |
| `GM_CONCURRENCY` | 并发数 |
| `GM_PROFILE` | 临时指定 Profile |

### 配置文件

- **macOS / Linux**：`~/.config/gradmotion/config.yaml`
- **Windows**：`%APPDATA%\gradmotion\config.yaml`

```yaml
profiles:
  default:
    base_url: https://spaces.gradmotion.com/prod-api
    timeout: 30s
    retry: 3
    concurrency: 4
current: default
```

## Cursor Agent 集成

Gradmotion CLI 提供 [Cursor Agent Skill](./npm/skills/gm-cli/SKILL.md)，安装后可用自然语言操作任务与项目。

```bash
# 安装 Skill
mkdir -p ~/.cursor/skills/gm-cli
cp npm/skills/gm-cli/SKILL.md ~/.cursor/skills/gm-cli/SKILL.md
```

在 Cursor Agent 中直接说：

- "帮我列出所有任务"
- "查看任务 task_xxx 的详情"
- "创建一个训练任务"
- "追踪任务 task_xxx 的实时日志"

详见 [快速上手教程](./docs/GETTING-STARTED.md)。

## 开发

### 前置条件

- Go 1.23+
- Git
- [GoReleaser](https://goreleaser.com)（发布时需要）

### 编译与测试

```bash
make build          # 编译当前平台
make build-all      # 编译所有平台（macOS/Linux/Windows × amd64/arm64）
make test           # 运行测试
make clean          # 清理产物
```

### 本地打包

```bash
make release-local  # GoReleaser 本地打包（生成压缩包至 dist/）
```

### 项目结构

```
gradmotion-cli/
├── cmd/gradmotion/          # 程序入口
├── internal/
│   ├── commands/            # 子命令实现（auth / config / project / task）
│   ├── config/              # 配置读取与合并
│   ├── auth/                # Keychain 管理
│   ├── client/              # HTTP Client（重试、超时、Trace ID）
│   ├── output/              # JSON / 人类可读输出
│   └── log/                 # JSONL 结构化日志
├── npm/                     # npm 发布包
│   ├── bin/                 # npm 入口脚本
│   ├── scripts/             # postinstall 平台适配
│   └── skills/gm-cli/       # Cursor Agent Skill
├── scripts/                 # 构建脚本
├── docs/                    # 文档
├── .goreleaser.yaml         # GoReleaser 配置
├── .github/workflows/       # GitHub Actions CI/CD
├── Makefile
└── go.mod
```

## 发布

推送符合语义化版本的 tag 后，[GitHub Actions](./.github/workflows/release.yml) 自动完成：

1. GoReleaser 多平台编译打包
2. 创建 GitHub Release 并上传产物
3. 发布 npm 包 `@limxdynamics/gm-cli`

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

详见 [编译与打包指南](./docs/BUILD.md)。

## 文档

| 文档 | 说明 |
|------|------|
| [快速上手](./docs/GETTING-STARTED.md) | 从 API Key 创建到 Cursor Agent Skill 配置的完整流程 |
| [编译与打包](./docs/BUILD.md) | 本地编译、跨平台打包、GoReleaser 发布指南 |
| [技术规格](./docs/Gradmotion-CLI-SPEC.md) | 架构设计、API 映射、输出规范 |
| [Agent Skill](./npm/skills/gm-cli/SKILL.md) | Cursor Agent 操作规范与参数校验规则 |
| [英文 README](./README.md) | 英文版本说明文档 |

## 许可证

[Apache License 2.0](./LICENSE)
