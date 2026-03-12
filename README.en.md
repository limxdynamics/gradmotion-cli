English | [中文](./README.md)

# Gradmotion CLI

Cross-platform CLI for managing training tasks and projects on [Gradmotion](https://spaces.gradmotion.com). It outputs structured JSON by default and is designed for AI Agent integration.

## Features

- **Cross-platform binary** - macOS / Linux / Windows (amd64 & arm64)
- **Secure authentication** - API key is stored in system keychain first (macOS Keychain / Windows Credential Manager / Linux Secret Service), with fallback to local config
- **Multi-profile config** - switch quickly between dev/staging/prod environments
- **Full task lifecycle** - create, edit, run, stop, delete, logs, hyperparameters, and batch operations
- **Agent-friendly output** - JSON by default + JSONL request logs for automation
- **Observability** - auto-injected Trace ID, `--debug`, and `--log-file`
- **Resilient requests** - built-in retry with exponential backoff, configurable timeout and concurrency

## Installation

### npm (recommended)

```bash
npm install -g @limxdynatic-gradmotion/gradmotion-cli
```

> Requires Node.js >= 16. The installer auto-selects the correct binary for your OS/architecture.

### Download from GitHub Releases

Download the package for your platform from [Releases](https://github.com/limxdynamics/gradmotion-cli/releases), extract it, then put `gm` into your `PATH`:

```bash
# macOS / Linux
sudo install -m 0755 gm /usr/local/bin/gm
```

### Build from source

```bash
git clone https://github.com/limxdynamics/gradmotion-cli.git
cd gradmotion-cli
make build
make install
```

### Verify

```bash
gm --version
```

## Quick Start

```bash
# 1) Set service endpoint
gm config set base_url "https://spaces.gradmotion.com/prod-api"

# 2) Login
gm auth login --api-key "gm_sk_your_key"

# 3) Verify
gm auth status
gm auth whoami

# 4) Start using
gm project list
gm task list
```

> Create API keys in Gradmotion: avatar menu (bottom-left) -> **API Key Management**. The full key is shown only once at creation time.

## Command Overview

```text
gm
├── auth        # login / logout / status / whoami
├── config      # set / get / profile
├── project     # list / create / info / edit / delete
└── task        # create / edit / run / stop / logs / batch / ...
```

Use `gm --help` or `gm <command> --help` for details.

## Common Examples

### Task Management

```bash
gm task list --page 1 --limit 50
gm task info --task-id "task_xxx"
gm task create --file ./create.json
gm task run --task-id "task_xxx"
gm task logs --task-id "task_xxx" --follow --interval 2s --timeout 5m
gm task stop --task-id "task_xxx"
```

### Project Management

```bash
gm project list --page 1 --limit 50
gm project create --file ./project-create.json
gm project info --project-id "proj_xxx"
```

### Profile Switching

```bash
gm config profile set dev --base-url "https://dev.gradmotion.com/prod-api" --timeout 30s
gm config profile use dev
gm --profile dev task list
```

### Output Control

```bash
gm task list --human
gm task list --quiet
gm task list --debug
gm task logs --task-id "task_xxx" --raw --no-request-log
```

## Configuration

### Priority

```text
CLI flags  >  Environment variables  >  Config file
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `GM_BASE_URL` | Service endpoint |
| `GM_API_KEY` | API key (temporary override, not persisted) |
| `GM_TIMEOUT` | Request timeout (for example `30s`) |
| `GM_RETRY` | Retry count |
| `GM_CONCURRENCY` | Concurrency limit |
| `GM_PROFILE` | Temporary profile |

### Config File

- **macOS / Linux**: `~/.config/gradmotion/config.yaml`
- **Windows**: `%APPDATA%\gradmotion\config.yaml`

```yaml
profiles:
  default:
    base_url: https://spaces.gradmotion.com/prod-api
    timeout: 30s
    retry: 3
    concurrency: 4
current: default
```

## Cursor Agent Integration

Gradmotion CLI provides a Cursor Agent Skill: [gm-cli SKILL](./npm/skills/gm-cli/SKILL.md).

```bash
mkdir -p ~/.cursor/skills/gm-cli
cp npm/skills/gm-cli/SKILL.md ~/.cursor/skills/gm-cli/SKILL.md
```

Try these prompts in Cursor Agent:

- "List all tasks"
- "Show details for task task_xxx"
- "Create a training task"
- "Follow logs for task task_xxx"

See also [Getting Started](./docs/GETTING-STARTED.md).

## Development

### Prerequisites

- Go 1.23+
- Git
- [GoReleaser](https://goreleaser.com) (for release)

### Build and Test

```bash
make build
make build-all
make test
make clean
```

### Local Release Packaging

```bash
make release-local
```

### Project Structure

```text
gradmotion-cli/
├── cmd/gradmotion/          # Program entrypoint
├── internal/
│   ├── commands/            # Subcommands (auth / config / project / task)
│   ├── config/              # Config loading and merge
│   ├── auth/                # Keychain management
│   ├── client/              # HTTP client (retry, timeout, Trace ID)
│   ├── output/              # JSON / human-readable output
│   └── log/                 # JSONL structured logs
├── npm/                     # npm package
│   ├── bin/                 # npm entry script
│   ├── scripts/             # postinstall platform adapter
│   └── skills/gm-cli/       # Cursor Agent Skill
├── scripts/                 # Build scripts
├── docs/                    # Documentation
├── .goreleaser.yaml         # GoReleaser config
├── .github/workflows/       # GitHub Actions CI/CD
├── Makefile
└── go.mod
```

## Release

Push a semantic version tag, and [GitHub Actions](./.github/workflows/release.yml) will:

1. Build/package all target platforms with GoReleaser
2. Create GitHub Release and upload artifacts
3. Publish npm package `@limxdynatic-gradmotion/gradmotion-cli`

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

See [Build Guide](./docs/BUILD.md) for details.

## Documents

| Doc | Description |
|-----|-------------|
| [Getting Started](./docs/GETTING-STARTED.md) | End-to-end workflow from API key to Cursor Agent integration |
| [Build & Release](./docs/BUILD.md) | Local build, cross-platform packaging, and GoReleaser guide |
| [Tech Spec](./docs/Gradmotion-CLI-SPEC.md) | Architecture, API mapping, and output standards |
| [Agent Skill](./npm/skills/gm-cli/SKILL.md) | Cursor Agent operation and parameter validation rules |

## License

[Apache License 2.0](./LICENSE)
