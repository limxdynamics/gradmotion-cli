# Gradmotion CLI

Cross-platform CLI for managing training tasks and projects on [Gradmotion](https://spaces.gradmotion.com). It outputs structured JSON by default and is designed for AI Agent integration.


## Features

- **Cross-platform binary** — macOS / Linux / Windows (amd64 & arm64)
- **Secure authentication** — API key stored in system keychain (macOS Keychain / Windows Credential Manager / Linux Secret Service)
- **Multi-profile config** — switch quickly between dev / staging / prod environments
- **Full task lifecycle** — create, edit, run, stop, delete, logs, hyperparameters, and batch operations
- **Agent-friendly output** — JSON by default + JSONL request logs for automation
- **Resilient requests** — built-in retry with exponential backoff, configurable timeout and concurrency

## Installation

```bash
npm install -g @jacky777/gm-cli
```

> Requires Node.js >= 16. The installer auto-selects the correct binary for your OS/architecture.

### Verify

```bash
gm --version
```

## Quick Start

```bash
# 1. Login (API key saved to system keychain)
gm auth login --api-key "gm_sk_your_key"

# 2. Verify connection
gm auth status        # local check
gm auth whoami        # server check

# 3. Start using
gm task list
gm project list
```

> Create API keys in Gradmotion: avatar menu (bottom-left) → **API Key Management**. The full key is shown only once at creation time.

## Command Overview

```text
gm
├── auth        login / logout / status / whoami
├── config      set / get / profile (list / set / use)
├── project     list / create / info / edit / delete
└── task        create / edit / copy / list / info / run / stop / delete
                logs / resource / image / storage / data / hp / env
                params / tag / batch / model
```

Use `gm --help` or `gm <command> --help` for details.

## Common Examples

```bash
# Task management
gm task list --page 1 --limit 50
gm task info --task-id "task_xxx"
gm task create --file ./create.json
gm task run --task-id "task_xxx"
gm task logs --task-id "task_xxx" --follow --interval 2s --timeout 5m
gm task stop --task-id "task_xxx"

# Project management
gm project list --page 1 --limit 50
gm project create --file ./project-create.json

# Profile switching
gm config profile set dev --base-url "https://dev.example.com/prod-api" --timeout 30s
gm config profile use dev

# Output control
gm task list --human          # human-readable table
gm task list --quiet          # key fields only
gm task list --debug          # debug logs
```

## Configuration

**Priority:** CLI flags > Environment variables > Config file

| Variable | Description |
|----------|-------------|
| `GM_BASE_URL` | Service endpoint |
| `GM_API_KEY` | API key (temporary, not persisted) |
| `GM_TIMEOUT` | Request timeout (e.g. `30s`) |
| `GM_RETRY` | Retry count |
| `GM_CONCURRENCY` | Concurrency limit |
| `GM_PROFILE` | Temporary profile |

Config file location:
- **macOS / Linux**: `~/.config/gradmotion/config.yaml`
- **Windows**: `%APPDATA%\gradmotion\config.yaml`


