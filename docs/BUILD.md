# Gradmotion CLI 编译与打包指南

本文档说明如何编译 Gradmotion CLI 并生成跨平台安装包。

---

## 快速开始

```bash
# 方式 1：Makefile（推荐）
make build          # 编译当前平台
make build-all      # 编译所有平台
make install        # 安装到系统（macOS）

# 方式 2：直接命令
go build -o gm ./cmd/gradmotion
./gm --version
```

---

## 一、单平台编译（开发用）

### 1.1 编译当前平台

```bash
cd /Users/limx/workspace/gradmotion-cli

# 方式 1：Makefile
make build

# 方式 2：go build
go build -o gm ./cmd/gradmotion

# 运行
./gm --version
```

### 1.2 安装到系统路径

```bash
# macOS / Linux
sudo install -m 0755 gm /usr/local/bin/gm
gm --version

# 或（用 Makefile）
make install
```

---

## 二、跨平台编译

### 2.1 用脚本一次编译所有平台

```bash
# 默认版本号 dev
./scripts/build-all.sh

# 指定版本号
VERSION=v0.1.0 ./scripts/build-all.sh

# 或用 Makefile
make build-all
VERSION=v0.1.0 make build-all
```

**产物**：`dist/` 目录
- `gm-darwin-amd64`（macOS Intel）
- `gm-darwin-arm64`（macOS Apple Silicon）
- `gm-linux-amd64`
- `gm-linux-arm64`
- `gm-windows-amd64.exe`
- `gm-windows-arm64.exe`

### 2.2 手动指定平台

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o gm-darwin-amd64 ./cmd/gradmotion

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o gm-darwin-arm64 ./cmd/gradmotion

# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o gm-linux-amd64 ./cmd/gradmotion

# Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o gm-windows-amd64.exe ./cmd/gradmotion
```

---

## 三、GoReleaser 打包（带压缩包、checksum）

### 3.1 前置条件

1. **初始化 Git 仓库**（如果还没有）：
```bash
git init
git add .
git commit -m "Initial commit"
```

2. **安装 GoReleaser**：
```bash
# macOS
brew install goreleaser

# 其他平台
go install github.com/goreleaser/goreleaser@latest
```

### 3.2 本地打包（不发布）

```bash
# 方式 1：用脚本
./scripts/release-local.sh

# 方式 2：直接命令
goreleaser release --clean --skip-publish --snapshot

# 方式 3：Makefile
make release-local
```

**产物**：`dist/` 目录
- `gradmotion-cli_<version>_darwin_amd64.tar.gz`
- `gradmotion-cli_<version>_darwin_arm64.tar.gz`
- `gradmotion-cli_<version>_linux_amd64.tar.gz`
- `gradmotion-cli_<version>_linux_arm64.tar.gz`
- `gradmotion-cli_<version>_windows_amd64.zip`
- `gradmotion-cli_<version>_windows_arm64.zip`
- `checksums.txt`

**参数说明**：
- `--clean`：清理旧的 `dist/`
- `--skip-publish`：不推送到 GitHub/GitLab
- `--snapshot`：不需要 tag，用当前 commit 生成版本号

### 3.3 正式发布（推送到 GitHub/GitLab）

#### 前置条件

1. 代码已推送到 GitHub 或 GitLab 远程仓库。
2. 已安装 GoReleaser：`brew install goreleaser` 或 `go install github.com/goreleaser/goreleaser/v2@latest`。
3. 工作区干净或只含本次要发布的提交（建议 `git status` 无未提交改动后再打 tag）。

---

#### 方式 A：发布到 GitHub

**1. 获取 GitHub Token**

- 打开 [GitHub → Settings → Developer settings → Personal access tokens](https://github.com/settings/tokens)。
- 新建 Token（Classic），勾选权限：`repo`（完整仓库权限）。
- 复制生成的 Token，不要提交到仓库。

**2. 设置环境变量并执行发布**

```bash
# 设置 Token（仅当前终端有效）
export GITHUB_TOKEN="ghp_xxxxxxxxxxxxxxxxxxxx"

# 1. 打 tag（版本号按语义化版本，如 v0.1.0）
git tag -a v0.1.0 -m "Release v0.1.0"

# 2. 推送 tag 到远程
git push origin v0.1.0

# 3. 执行 GoReleaser（会构建、打包、创建 Release、上传产物）
goreleaser release --clean
```

无需修改 `.goreleaser.yaml`，GoReleaser 会根据 `git remote` 自动识别 GitHub 仓库。

**3. 结果**

- 在 GitHub 仓库的 **Releases** 页会出现 `v0.1.0`。
- 包含各平台压缩包（tar.gz/zip）和 `checksums.txt`，changelog 由 git 历史生成。

---

#### 方式 B：发布到 GitLab

**1. 获取 GitLab Token**

- GitLab 项目 → **Settings → Access Tokens**（或用户 **Preferences → Access Tokens**）。
- 新建 Token，勾选 `api`、`read_repository`、`write_repository`。
- 复制 Token，不要提交到仓库。

**2. 配置 `.goreleaser.yaml`**

在 `.goreleaser.yaml` 末尾取消注释并填写 GitLab 仓库信息：

```yaml
release:
  gitlab:
    owner: your-username    # 用户名或组名
    name: gradmotion-cli   # 仓库名
```

若使用**自建 GitLab**（非 gitlab.com），还需加上：

```yaml
gitlab_urls:
  api: "https://gitlab.example.com/api/v4"
  download: "https://gitlab.example.com"
```

**3. 设置环境变量并执行发布**

```bash
export GITLAB_TOKEN="glpat-xxxxxxxxxxxxxxxxxxxx"

git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

goreleaser release --clean
```

**4. 结果**

- 在 GitLab 项目的 **Deploy → Releases** 中会出现 `v0.1.0`，并附带构建产物与 checksum。

---

#### 发布时 GoReleaser 会做什么

- 根据当前 tag 编译多平台二进制（darwin/linux/windows × amd64/arm64）。
- 打包为 tar.gz（Windows 为 zip）并生成 `checksums.txt`。
- 在 GitHub/GitLab 上创建对应版本的 Release。
- 上传所有压缩包与 checksum 到该 Release。

---

#### 常见问题

| 情况 | 处理 |
|------|------|
| 发布错版本 / 想重跑 | 在网页上删除该 Release 和 tag，删除本地 tag：`git tag -d v0.1.0`，修正后重新打 tag 并 `git push origin v0.1.0`，再执行 `goreleaser release --clean`。 |
| 只想本地打包、不推送 | 使用 `make release-local` 或 `goreleaser release --clean --skip-publish --snapshot`。 |
| CI 里自动发布 | 在 CI 中配置 `GITHUB_TOKEN` 或 `GITLAB_TOKEN`（敏感变量），在 tag 推送后执行 `goreleaser release --clean`。 |

---

## 四、配置发布到 GitLab（参考）

发布到 **GitLab** 时，需在 `.goreleaser.yaml` 中显式配置 `release.gitlab`（见上文 3.3 方式 B）。发布到 **GitHub** 时无需该配置。自建 GitLab 需同时配置 `gitlab_urls`。

---

## 五、测试安装包

### macOS

```bash
cd dist
tar -xzf gradmotion-cli_*_darwin_arm64.tar.gz  # 或 darwin_amd64
sudo install -m 0755 gm /usr/local/bin/gm
gm --version
```

### Linux

```bash
tar -xzf gradmotion-cli_*_linux_amd64.tar.gz
sudo install -m 0755 gm /usr/local/bin/gm
gm --version
```

### Windows

解压 `gradmotion-cli_*_windows_amd64.zip`，将 `gm.exe` 添加到 `PATH`。

---

## 六、推荐流程

| 场景 | 方法 | 命令 |
|------|------|------|
| 本地开发测试 | 单平台编译 | `make build` 或 `go build` |
| 交付多平台二进制 | 跨平台脚本 | `make build-all` 或 `./scripts/build-all.sh` |
| 正式发布（带压缩包） | GoReleaser | `make release-local`（本地）或 `goreleaser release`（推送） |
| CI/CD 自动发布 | GitLab CI + Tag | 推送 tag 后自动发布 Release 与 npm，见下方「八」 |

---

## 八、GitLab CI 流水线发布

推送符合语义化版本的 **tag**（如 `v0.1.1`）后，流水线会自动完成 **GitLab Release** 与 **npm** 发布。

### 8.1 配置 CI/CD 变量

在 GitLab 项目 **Settings → CI/CD → Variables** 中新增（建议勾选 Masked）：

| 变量名 | 说明 |
|--------|------|
| `GITLAB_TOKEN` | 用于创建 Release 和上传产物（Personal/Project Access Token，需 `api` 权限） |
| `NPM_TOKEN` | 用于 `npm publish`（npm 的 Automation Token） |

### 8.2 发布步骤

1. 在仓库中更新代码并提交到 `prod`（或当前发布分支）。
2. 打 tag 并推送：
   ```bash
   git tag -a v0.1.1 -m "Release v0.1.1"
   git push origin v0.1.1
   ```
3. 流水线自动运行：构建多平台二进制 → 创建 GitHub Release 并上传 → 发布 npm 包 `@limxdynatic-gradmotion/gradmotion-cli@0.1.0`。

### 8.3 触发规则

仅当推送的 tag 匹配 `v*.*.*`（如 `v0.1.0`、`v1.0.0`）时才会触发发布流水线。

---

## 七、常用命令速查

```bash
# 开发
make build              # 编译当前平台
make install            # 安装到系统
make clean              # 清理产物

# 打包
make build-all          # 所有平台（无需 git）
make release-local      # GoReleaser 本地打包（需 git）

# 测试
make test               # 运行测试
./gm --version          # 查看版本
```

---

## 附：Makefile 目标

运行 `make help` 查看所有可用目标：

```bash
$ make help
Gradmotion CLI Makefile

Usage:
  make build          - 编译当前平台
  make build-all      - 编译所有平台（无需 git）
  make release-local  - GoReleaser 本地打包（需 git）
  make clean          - 清理产物
  make install        - 安装到 /usr/local/bin（需 sudo）
  make test           - 运行测试

环境变量：
  VERSION=v0.1.0 make build-all  - 指定版本号
```
