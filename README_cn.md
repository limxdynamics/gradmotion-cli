# 中文 | [English](README.md)


# Flux CLI (`gm`)

[![npm](https://img.shields.io/npm/v/@limxdynamics/gm-cli)](https://www.npmjs.com/package/@limxdynamics/gm-cli)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](./LICENSE)
[![Node](https://img.shields.io/badge/node-%3E%3D16-brightgreen)](https://nodejs.org)
[![Platforms](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)]()

跨平台命令行工具，用于管理 LimX Dynamics 训练平台 [Flux](https://gm-spaces.limxdynamics.com) 上的项目（project）与训练任务（task）：创建、运行、监控、拉取日志与模型 checkpoint 等。

CLI 默认输出结构化 JSON，面向脚本自动化与 AI Agent 集成设计——可以把 `gm` 交给编程 Agent 去操作训练任务，不需要额外写 SDK。

> **关于本仓库**：本仓库当前用作 Flux CLI 的公开主页（仅含 `LICENSE`，GitHub 上没有发布 Release 或二进制文件）。工具的**唯一发布渠道是 npm 包** [`@limxdynamics/gm-cli`](https://www.npmjs.com/package/@limxdynamics/gm-cli)，其中打包了 macOS / Linux / Windows（amd64 / arm64）共 6 个平台的 Go 编译产物。下文所有命令与输出，均基于对 npm 最新发布版本（`0.1.14`）的实际安装与运行验证。

## 功能特性

- **跨平台单文件二进制**：macOS / Linux / Windows，amd64 / arm64。`npm install` 只是根据平台选出 `vendor/` 里对应的二进制并直接调用（执行阶段不依赖 Node.js，只有安装阶段用到）。
- **Agent 友好**：默认输出结构化 JSON（`{success, data, meta, error}`），并附带 [`skills/gm-cli/SKILL.md`](#ai-agent-集成) 技能文件（配套文档写的是配合 **Cursor 的 Agent Skills** 使用；这个格式理论上也可能适配其他编程 Agent 的技能系统，但目前没有文档明确支持）。
- **完整项目/任务生命周期**：project 与 task 的增删改查、运行/停止、日志、超参数、checkpoint 模型列表、图表数据、批量操作等。
- **多 profile 配置**：可在多个环境（dev / staging / prod）间快速切换 `base_url`、超时、并发等参数。
- **安全鉴权存储**：`auth login` 保存的 API Key 会写入系统密钥库（macOS Keychain / Windows 凭据管理器 / Linux Secret Service），仅在本地无密钥库时才会回退到配置文件。
- **健壮请求**：内置重试（指数退避）、可配置超时与并发；所有写操作支持 `--dry-run` 预览请求而不实际执行。
- **明确的退出码**：不同错误类型对应不同退出码，方便脚本/CI 判断执行结果。

## 安装

```bash
npm install -g @limxdynamics/gm-cli
```

> 需要 Node.js >= 16（仅用于 npm 分发与选择二进制，不影响运行速度）。安装脚本会根据 `process.platform` / `process.arch` 自动挑选 `vendor/gm-<os>-<arch>` 对应的可执行文件。

安装完成后验证：

```bash
$ gm --version
{
  "success": true,
  "data": {
    "version": "0.1.14",
    "commit": "b2ef166fbf899eca738815dfb558599cdf0b9caf",
    "date": "2026-06-12T03:49:51Z"
  }
}
```

## 快速开始

```bash
# 1. 登录（API Key 会写入系统密钥库 / 配置文件）
gm auth login --api-key "gm_sk_your_key"

# 2. 本地检查登录状态 & 向服务端确认身份
gm auth status
gm auth whoami

# 3. 开始使用
gm project list
gm task list --page 1 --limit 20
```

> 在 Flux 中创建 API Key：左下角头像菜单 → **API Key 管理**。完整 Key 仅在创建时显示一次，请妥善保存。

## 命令总览

```text
gm
├── auth        login / logout / status / whoami
├── config      set / get / profile (list / set / use)
├── project     list / create / info / edit / delete
├── task        create / edit / copy / list / info / run / stop / delete
│               logs / image / storage / data / hp / env
│               params / tag / note / model / resource / batch
└── completion  bash / zsh / fish / powershell
```

用 `gm --help` 或 `gm <command> [子命令] --help` 查看任意层级的详细说明，每个子命令自带 `Examples` 区块。

### 全局参数

| Flag | 说明 |
|---|---|
| `--api-key` | 临时指定 API Key（不落盘，优先级最高） |
| `--base-url` | 临时指定服务端地址 |
| `--profile` | 使用指定的配置 profile |
| `--dry-run` | 只打印将要发出的请求（method / endpoint / body），不实际执行 |
| `--yes` | 跳过交互式二次确认（Agent / CI 等非交互场景下危险操作**必须**加此参数） |
| `--human` | 人类可读格式输出（默认是 JSON） |
| `--quiet` | 只输出关键字段 |
| `--debug` | 打印调试日志 |
| `--timeout` / `--retry` / `--concurrency` | 请求超时、重试次数、并发度 |
| `--log-file` | 将日志写入指定文件 |

## 常用示例

```bash
# —— 任务管理 ——
gm task list --status "3,5"                     # 3=运行中 5=已完成
gm task info --task-id "task_xxx"
gm task create --file ./create.json              # 推荐用 JSON 文件描述任务
gm task create --file ./create.json --dry-run    # 先预览，确认无误再去掉 --dry-run
gm task run --task-id "task_xxx"
gm task logs --task-id "task_xxx" --follow --interval 2s --timeout 5m
gm --yes task stop --task-id "task_xxx"
gm task model list --task-id "task_xxx"          # 查看训练产出的 checkpoint

# —— 项目管理 ——
gm project list --page 1 --limit 50
gm project create --data '{"projectName":"my-project"}'

# —— Profile 切换 ——
gm config profile set dev --base-url "https://dev.example.com/prod-api" --timeout 30s
gm config profile use dev

# —— 输出格式 ——
gm task list --human     # 表格
gm task list --quiet     # 只要关键字段
```

`task create` 所需的最小字段示例：

```bash
gm task create --data '{
  "taskBaseInfo": {
    "projectId": "PRO_xxx", "taskType": "1", "trainType": "1",
    "taskName": "my-task", "goodsId": "ESKxxx",
    "imageId": "BJXxxx", "imageVersion": "V0000xx"
  },
  "taskCodeInfo": {
    "codeType": "2",
    "codeUrl": "[{\"codeUrl\":\"https://github.com/org/repo.git\",\"versionType\":\"1\",\"versionName\":\"main\"}]",
    "isOpen": "1",
    "startScript": "gm-run xxx/train.py --task=xxx --headless --max_iterations=3000"
  }
}'
```

## 配置

**优先级：** 命令行参数 > 环境变量 > 配置文件

| 环境变量 | 说明 |
|---|---|
| `GM_BASE_URL` | 服务地址 |
| `GM_API_KEY` | API Key（临时使用，不落盘） |
| `GM_TIMEOUT` | 请求超时时间（如 `30s`） |
| `GM_RETRY` | 重试次数 |
| `GM_CONCURRENCY` | 并发数限制 |
| `GM_PROFILE` | 临时 profile |

配置文件位置：
- macOS / Linux：`~/.config/gradmotion/config.yaml`
- Windows：`%APPDATA%\gradmotion\config.yaml`

文件内容示例（实测 `gm auth login` 生成，系统密钥库可用的情况下，API Key 本身不会落到这个文件里）：

```yaml
current: prod
profiles:
    prod:
        base_url: ""
        api_key: ""
        timeout: 30s
        retry: 3
        concurrency: 4
```

## 退出码

| 码 | 含义 |
|---|---|
| `0` | 成功 |
| `1` | 通用错误（网络、服务端 5xx、未知异常） |
| `2` | 参数错误（flag 缺失/非法，或非交互环境下危险操作未加 `--yes`） |
| `3` | 资源不存在（404） |
| `4` | 权限不足（401 / 403） |
| `5` | 冲突 / 已存在（409） |
| `10` | dry-run 校验通过，可安全去掉 `--dry-run` 正式执行 |

## AI Agent 集成

npm 包内随附 [`skills/gm-cli/SKILL.md`](https://www.npmjs.com/package/@limxdynamics/gm-cli) 与 `docs/GradmotionCLISkill-Guide.md`。官方教程文档写的是配合 **Cursor 的 Agent Skills** 安装使用（`~/.cursor/skills/gm-cli`）；技能文件的格式理论上也可能适配其他编程 Agent 的技能系统，但文档里没有明确写支持。技能文件里的关键约束：

- **非交互环境会自动检测**：CLI 检测到 stdin 非终端时，危险操作（create/edit/run/stop/delete 等）不会挂起等待确认，而是直接以退出码 `2` 失败——Agent 执行这类命令时**必须显式加 `--yes`**。
- **写操作建议先 `--dry-run`**：预览请求 body / endpoint 确认无误后，再去掉 `--dry-run` 正式执行。
- 涉及文件建议使用相对路径（如 `--file ./payload.json`），执行成功后清理临时 JSON 文件。
- 所有输出默认是 JSON，便于 Agent 直接解析而无需匹配人类可读文本。

## License

本仓库根目录 [LICENSE](./LICENSE) 为 **Apache License 2.0**。

---

<sub>本 README 基于对 npm 包 `@limxdynamics/gm-cli@0.1.14` 的实际下载、解包与运行验证编写，不是照抄其他文档翻译的。如果后续版本行为有差异，以 `gm <command> --help` 的实时输出为准。</sub>
