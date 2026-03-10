# Gradmotion CLI 快速上手教程

> 本教程覆盖从 API Key 创建、CLI 本地安装，到在 Cursor Agent 中配置 CLI Skill 的完整流程。

---
## 最简流程：从 API Key 到 Skill 使用（约 5 分钟）

按下面顺序做完即可在 Cursor 里用自然语言操作 Gradmotion 任务。

| 步骤 | 操作 |
|------|------|
| **1. 创建 API Key** | 登录 Gradmotion → 左下角头像 → **API Key 管理** → **新建** → 填写名称并**立即复制**生成的 Key（格式 `gm_sk_xxxxxxxx`，仅显示一次）。 |
| **2. 安装 CLI** | 执行 `npm install -g @limxdynatic-gradmotion/gradmotion-cli`（需 Node.js >= 16）。安装后执行 `gm --version` 确认。 |
| **3. 配置** | `gm config set base_url "https://spaces.gradmotion.com/prod-api"`<br>`gm auth login --api-key "gm_sk_你的Key"` |
| **4. 验证** | `gm auth status`（本地）→ `gm auth whoami`（服务端）。 |
| **5. 安装 Skill** | `mkdir -p ~/.cursor/skills/gm-cli`，将 [SKILL.md](../npm/skills/gm-cli/SKILL.md) 放到该目录；或在 Cursor **Settings → Features → Agent Skills** 确认 Skills 路径为 `~/.cursor/skills/`。 |
| **6. 使用** | 在 Cursor Agent 里直接说「帮我列出任务」「查看任务 task_xxx 详情」等，或输入 `gm task list`；提到 `gm`、Gradmotion、任务、项目时 Agent 会自动使用 gm-cli Skill。 |

详细说明与常见问题见下文各章节。

---


## 目录

1. [创建 API Key](#1-创建-api-key)
2. [安装 CLI](#2-安装-cli)
3. [配置 CLI](#3-配置-cli)
4. [验证连接](#4-验证连接)
5. [在 Cursor Agent 中安装 CLI Skill](#5-在-cursor-agent-中安装-cli-skill)
6. [通过 Agent 使用 CLI](#6-通过-agent-使用-cli)
7. [常见问题](#7-常见问题)

---

## 1. 创建 API Key

### 1.1 登录 Gradmotion 平台

打开浏览器，访问 Gradmotion 平台并登录你的账号。

### 1.2 进入 API Key 管理页面

登录后，点击左下角用户头像 → **API Key 管理**。

### 1.3 创建新的 API Key

1. 点击 **新建** 按钮。
2. 填写 Key 的名称（如 `local-dev`、`cursor-agent`），便于区分用途。
3. 点击 **确认**。
4. **立即复制并妥善保存生成的 Key**（格式为 `gm_sk_xxxxxxxxxxxxxxxx`），该 Key 仅在创建时完整显示一次。

> **安全提示**：API Key 具有账号级别的操作权限，请勿提交到代码仓库、分享给他人，或在公开场合展示。

---

## 2. 安装 CLI

推荐通过 npm 一行命令安装，自动匹配当前系统和架构（macOS / Linux / Windows，x64 / arm64）。

### 2.1 通过 npm 安装（推荐）

```bash
npm install -g @limxdynatic-gradmotion/gradmotion-cli
```

> 需要 Node.js >= 16，安装后会自动选择对应平台的二进制。

### 2.2 验证安装

```bash
gm --version
```

输出示例：

```
gradmotion-cli v0.1.0 (darwin/arm64)
```

---

## 3. 配置 CLI

### 3.1 设置服务地址（base_url）

```bash
gm config set base_url "https://spaces.gradmotion.com/prod-api"
```

> `base_url` 为 Gradmotion 后端服务地址，CLI 会在此基础上拼接 `/api` 前缀发起请求。

### 3.2 保存 API Key

```bash
gm auth login --api-key "gm_sk_xxxxxxxxxxxxxxxx"
```

CLI 会优先将 Key 存入系统 Keychain（macOS Keychain / Windows Credential Manager / Linux Secret Service），若 Keychain 不可用则回退存储到本地配置文件 `~/.config/gradmotion/config.yaml`。

成功输出示例：

```json
{
  "success": true,
  "data": {
    "profile": "default",
    "saved_to": "keychain"
  }
}
```

### 3.3 配置文件说明

配置文件默认路径：

- **macOS / Linux**：`~/.config/gradmotion/config.yaml`
- **Windows**：`%APPDATA%\gradmotion\config.yaml`

配置结构示例：

```yaml
profiles:
  default:
    base_url: https://spaces.gradmotion.com/prod-api
    timeout: 30s
    retry: 3
    concurrency: 4
current: default
```

### 3.4 多环境 Profile（可选）

如果你需要同时管理多个环境（如开发、测试、生产），可以使用 Profile：

```bash
# 创建 dev profile
gm config profile set dev \
  --base-url "https://dev.gradmotion.com/prod-api" \
  --timeout 30s \
  --retry 3

# 切换到 dev profile
gm config profile use dev

# 查看所有 profile
gm config profile list

# 临时使用某个 profile（不切换默认）
gm --profile dev task list
```

---

## 4. 验证连接

### 4.1 检查本地认证状态

```bash
gm auth status
```

输出示例：

```json
{
  "success": true,
  "data": {
    "profile": "default",
    "base_url": "https://spaces.gradmotion.com/prod-api",
    "has_api_key": true,
    "key_source": "keychain"
  }
}
```

### 4.2 验证服务端连接

```bash
gm auth whoami
```

成功后会返回当前账号的用户信息，确认 CLI 与 Gradmotion 服务端通信正常。

### 4.3 快速任务列表测试

```bash
gm task list
```

### 4.4 执行前快速探测（推荐）

在执行创建/编辑/删除前，建议先确认当前 CLI 具备完整命令：

```bash
gm --help
gm task --help
gm project --help
```

---

## 5. 在 Cursor Agent 中安装 CLI Skill

CLI Skill 让 Cursor 中的 AI Agent 能够理解并正确调用 `gm` 命令，包括参数校验、安全约束和最佳实践。

### 5.1 找到 Skills 存放目录

Cursor 的 Skills 通常存放在 `~/.cursor/skills/` 目录下。

```bash
# 查看当前 skills 目录
ls ~/.cursor/skills/
```

### 5.2 创建 gm-cli Skill 目录

```bash
mkdir -p ~/.cursor/skills/gm-cli
```

### 5.3 下载 SKILL.md 文件

点击下方链接下载 gm-cli Skill 文件：

**[下载 SKILL.md](../npm/skills/gm-cli/SKILL.md)**

将下载的 `SKILL.md` 放到 `~/.cursor/skills/gm-cli/` 目录下：

```bash
mkdir -p ~/.cursor/skills/gm-cli
mv ~/Downloads/SKILL.md ~/.cursor/skills/gm-cli/SKILL.md
```

> 该 Skill 文件包含 `gm` 命令的完整操作规范、参数校验规则和安全约束，供 Cursor Agent 在执行 CLI 操作时自动参考。

### 5.4 在 Cursor 中注册 Skill

打开 Cursor，进入 **Settings（齿轮图标）→ Features → Agent Skills**，确认 Skills 目录路径已正确配置为 `~/.cursor/skills/`。

或者，直接在对话框中 `@` 引用该 Skill 即可激活。

> **确认 Skill 已生效**：在 Cursor 的 Agent 对话框中输入 `gm task list`，Agent 会自动识别 gm-cli 上下文并给出正确的操作建议。

---

## 6. 通过 Agent 使用 CLI

安装 CLI Skill 后，你可以直接用自然语言让 Cursor Agent 帮你操作 `gm` 命令。

### 6.1 典型工作流示例

**场景一：查看任务列表**

在 Cursor Agent 中输入：

```
帮我查看所有任务列表
```

Agent 会执行：

```bash
gm task list --page 1 --limit 50
```

---

**场景二：查看指定任务详情**

```
查看任务 task_123 的详情和当前状态
```

Agent 会执行：

```bash
gm task info --task-id "task_123"
```

---

**场景三：创建并运行任务**

```
帮我创建一个训练任务，项目 ID 是 proj_001，任务名称是 test-run
```

Agent 会先提示你补充必要参数（如 `codeType`、`goodsId` 等），然后生成：

```bash
gm task create --file ./create.json
gm task run --task-id "task_xxx"
```

---

**场景四：实时查看任务日志**

```
追踪任务 task_123 的实时日志，最多等 5 分钟
```

Agent 会执行：

```bash
gm task logs --task-id "task_123" --follow --interval 2s --timeout 5m
```

---

**场景五：停止任务（高风险，会确认）**

```
停止任务 task_123
```

Agent 会先执行状态预检：

```bash
gm task info --task-id "task_123"
```

确认任务可停止后，提示你二次确认，再执行：

```bash
gm task stop --task-id "task_123"
```

---

**场景六：先选项目，再创建任务**

```
先帮我列出项目，然后在 proj_001 下面创建一个训练任务
```

Agent 通常会按顺序执行：

```bash
gm project list --page 1 --limit 50
gm task create --file ./create.json
gm task run --task-id "task_xxx"
```

---

### 6.2 环境变量方式（适用于 CI/脚本场景）

如不想将 Key 持久化到本地，可通过环境变量临时注入：

```bash
export GM_BASE_URL="https://spaces.gradmotion.com/prod-api"
export GM_API_KEY="gm_sk_xxxxxxxxxxxxxxxx"

gm task list
```

---

## 7. 常见问题

### Q1：`gm auth status` 显示 `has_api_key: false`

**原因**：API Key 尚未保存或已被清除。

**解决**：

```bash
gm auth login --api-key "gm_sk_xxxxxxxxxxxxxxxx"
```

---

### Q2：`gm auth whoami` 返回 `PERMISSION_DENIED`

**原因**：API Key 无效、已过期或已被撤销。

**解决**：前往 Gradmotion 平台 → API Key 管理，检查 Key 状态，必要时重新创建并重新登录。

```bash
gm auth logout
gm auth login --api-key "gm_sk_新的Key"
```

---

### Q3：`base_url` 报错或请求失败

**原因**：`base_url` 未设置或地址不正确。

**解决**：

```bash
gm config set base_url "https://spaces.gradmotion.com/prod-api"
gm auth status  # 确认 base_url 已生效
```

---

### Q4：macOS Keychain 弹窗请求访问权限

**原因**：CLI 首次将 API Key 写入系统 Keychain 时，macOS 会弹出授权确认。

**解决**：点击 **允许（Allow）** 即可，后续不再弹出。

---

### Q5：CLI Skill 在 Agent 中没有生效

**检查步骤**：

1. 确认 `~/.cursor/skills/gm-cli/SKILL.md` 文件存在且内容正确。
2. 重启 Cursor。
3. 在 Agent 对话中明确提及 `gm` 或 `gm-cli` 关键词，触发 Skill 激活。

---

### Q6：`gm task stop / delete` 执行时需要输入确认

**这是正常行为**。高风险操作默认需要二次确认，防止误操作。

如需在脚本中跳过确认（无人值守模式）：

```bash
gm --yes task stop --task-id "task_xxx"
gm --yes task delete --task-id "task_xxx"
```

> **注意**：仅在确认安全的自动化场景中使用 `--yes`。

---

### Q7：`gm task storage list` 和其他接口路径规则不一致？

这是预期行为。大多数命令走 `/api/...`，但 `gm task storage list` 使用的是绝对路径接口 `/gm/storage/list`。  
CLI 已内置兼容，无需手动处理路径拼接。

---

## 附录：常用命令速查

### 帮助与版本

| 操作 | 命令 |
|------|------|
| 总览帮助 | `gm --help` |
| 子命令帮助 | `gm task --help`、`gm project --help` |
| 查看版本 | `gm --version` 或 `gm -v` |

### 认证

| 操作 | 命令 |
|------|------|
| 保存 API Key | `gm auth login --api-key "<KEY>"` |
| 检查本地认证 | `gm auth status` |
| 验证服务端连接 | `gm auth whoami` |
| 退出登录 | `gm auth logout` |

### 配置与 Profile

| 操作 | 命令 |
|------|------|
| 设置服务地址 | `gm config set base_url "<URL>"` |
| Profile 列表 | `gm config profile list` |
| 创建/更新 Profile | `gm config profile set <name> --base-url "<URL>" --timeout 30s --retry 3 --concurrency 4` |
| 切换 Profile | `gm config profile use <name>` |
| 临时指定 Profile | `gm --profile <name> task list` |

### 项目

| 操作 | 命令 |
|------|------|
| 项目列表 | `gm project list --page 1 --limit 50` |
| 创建项目 | `gm project create --file ./project-create.json` |
| 项目详情 | `gm project info --project-id "<ID>"` |
| 编辑项目 | `gm project edit --file ./project-edit.json` 或 `--data '{"projectId":"proj_xxx","projectName":"新名称"}'` |
| 删除项目 | `gm project delete --project-id "<ID>"`（需确认，可加 `--yes`） |

### 任务（创建 / 编辑 / 生命周期）

| 操作 | 命令 |
|------|------|
| 查看任务列表 | `gm task list --page 1 --limit 50` |
| 查看任务详情 | `gm task info --task-id "<ID>"` |
| 创建任务 | `gm task create --file ./create.json` 或 `--data '{"...":"..."}'` |
| 编辑任务 | `gm task edit --file ./edit.json`（需先 `task info` 取全量再改） |
| 复制任务 | `gm task copy --file ./task-copy.json` |
| 运行任务 | `gm task run --task-id "<ID>"` |
| 停止任务 | `gm task stop --task-id "<ID>"`（需确认，可加 `--yes`） |
| 删除任务 | `gm task delete --task-id "<ID>"`（需确认，可加 `--yes`） |

### 任务日志

| 操作 | 命令 |
|------|------|
| 查看任务日志 | `gm task logs --task-id "<ID>"` |
| 实时追踪日志 | `gm task logs --task-id "<ID>" --follow --interval 2s --timeout 1m` |
| 仅输出日志正文（不包 JSON，便于管道） | `gm task logs --task-id "<ID>" --raw` |
| 不向 stderr 输出请求元数据 | `gm task logs ... --no-request-log` |

### 资源 / 镜像 / 存储

| 操作 | 命令 |
|------|------|
| 资源列表（算力） | `gm task resource list --goods-back-category 3 --page-num 1 --page-size 10` |
| 官方镜像 | `gm task image official` |
| 个人镜像 | `gm task image personal --version-status 1 --page-num 1 --page-size 50` |
| 镜像版本 | `gm task image versions --image-id "<ID>"` |
| 个人存储列表 | `gm task storage list --folder-path "personal/"` |

### 数据 / 超参 / 环境

| 操作 | 命令 |
|------|------|
| 图表 Key 列表 | `gm task data keys --task-id "<ID>"` |
| 图表数据查询 | `gm task data get --task-id "<ID>" --data-key "train/loss"` |
| 图表数据下载 | `gm task data download --task-id "<ID>"` |
| 超参读取 | `gm task hp get --task-id "<ID>"` |
| 超参提交 | `gm task params submit --task-id "<ID>" --file ./params.json` |
| 超参更新 | `gm task params update --task-id "<ID>" --data '{"...":"..."}'` |
| 运行环境 | `gm task env get --task-id "<ID>"` |

### 任务标签

| 操作 | 命令 |
|------|------|
| 更新任务标签 | `gm task tag update --task-id "<ID>" --tags "tag1,tag2"` 或 `--file ./tag.json` |
| 查看任务标签 | `gm task tag get --task-id "<ID>"` |
| 用户历史标签列表 | `gm task tag list --limit 200` |

### 批量操作（高风险，默认需确认）

| 操作 | 命令 |
|------|------|
| 批量停止 | `gm task batch stop --task-ids "id1,id2,id3"` |
| 批量删除 | `gm task batch delete --task-ids "id1,id2,id3"` |

### 输出与调试

| 操作 | 命令 |
|------|------|
| 人类可读输出 | `gm task list --human` |
| 仅关键字段 | `gm task list --quiet` |
| 开启调试日志 | `gm task list --debug` |
| 日志写入文件 | `gm task list --log-file ./gm.log` |
| 临时覆盖 base_url / api-key（不落盘） | `gm --base-url "<URL>" --api-key "<KEY>" auth whoami` |
