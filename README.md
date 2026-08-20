# English | [中文](README_cn.md)


# Flux CLI (`gm`)

[![npm](https://img.shields.io/npm/v/@limxdynamics/gm-cli)](https://www.npmjs.com/package/@limxdynamics/gm-cli)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](./LICENSE)
[![Node](https://img.shields.io/badge/node-%3E%3D16-brightgreen)](https://nodejs.org)
[![Platforms](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)]()

Cross-platform CLI for managing projects and training tasks on
Flux, LimX Dynamics' training
platform: create, run, monitor, pull logs and model checkpoints, and more.

The CLI outputs structured JSON by default and is built for scripting and
AI-agent automation — you can hand `gm` to a coding agent to drive training
tasks without writing a separate SDK integration.

> **About this repository**: this repo currently serves as the public
> homepage for the Flux CLI (it only contains `LICENSE` — there are no
> GitHub Releases or binaries here). The **only distribution channel is
> the npm package**
> [`@limxdynamics/gm-cli`](https://www.npmjs.com/package/@limxdynamics/gm-cli),
> which bundles Go-compiled binaries for 6 platform/arch combinations
> (macOS / Linux / Windows × amd64 / arm64). All commands and output shown
> below were verified by actually installing and running the latest
> published version (`0.1.14`).

## Features

- **Cross-platform single-binary** — macOS / Linux / Windows, amd64 /
  arm64. `npm install` just picks the right binary from `vendor/` and
  execs it directly (Node.js isn't on the execution path, only the
  installer).
- **Agent-friendly** — structured JSON output (`{success, data, meta,
  error}`) by default, plus a bundled
  [`skills/gm-cli/SKILL.md`](#ai-agent-integration) skill file (used with
  Cursor's Agent Skills; the format may work with other coding-agent skill
  systems too).
- **Full project/task lifecycle** — create, edit, run, stop, list, logs,
  hyperparameters, checkpoint/model listing, chart data, batch ops.
- **Multi-profile config** — switch `base_url`, timeout, and concurrency
  quickly between dev / staging / prod.
- **Secure credential storage** — `auth login` saves the API key to the
  system keychain (macOS Keychain / Windows Credential Manager / Linux
  Secret Service), falling back to the config file only if no keychain is
  available.
- **Resilient requests** — retry with exponential backoff, configurable
  timeout/concurrency; every write operation supports `--dry-run`.
- **Explicit exit codes** — distinct exit codes per error category, for
  scripting/CI.

## Installation

```bash
npm install -g @limxdynamics/gm-cli
```

> Requires Node.js >= 16 (used only for npm distribution and binary
> selection — it doesn't affect runtime performance). The installer picks
> `vendor/gm-<os>-<arch>` based on `process.platform` / `process.arch`.

Verify:

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

## Quick Start

```bash
# 1. Log in (API key goes to the system keychain / config file)
gm auth login --api-key "gm_sk_your_key"

# 2. Verify: local check, then a server-side check
gm auth status
gm auth whoami

# 3. Start using it
gm project list
gm task list --page 1 --limit 20
```

> Create API keys in Flux: avatar menu (bottom-left) → **API Key
> Management**. The full key is shown only once, at creation time — save
> it.

## Command Overview

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

Use `gm --help` or `gm <command> [subcommand] --help` for details at any
level — every subcommand ships its own `Examples` section.

### Global flags

| Flag | Description |
|---|---|
| `--api-key` | API key for this invocation only (not persisted, highest priority) |
| `--base-url` | Override the service endpoint for this invocation |
| `--profile` | Use a specific config profile |
| `--dry-run` | Print the request (method / endpoint / body) without sending it |
| `--yes` | Skip the interactive confirmation (**required** for dangerous operations in Agent/CI, non-interactive contexts) |
| `--human` | Human-readable output (default is JSON) |
| `--quiet` | Only print key fields |
| `--debug` | Print debug logs |
| `--timeout` / `--retry` / `--concurrency` | Request timeout, retry count, concurrency |
| `--log-file` | Write logs to a file |

## Common Examples

```bash
# —— Task management ——
gm task list --status "3,5"                     # 3=running, 5=finished
gm task info --task-id "task_xxx"
gm task create --file ./create.json              # recommended: describe the task in a JSON file
gm task create --file ./create.json --dry-run    # preview first, then drop --dry-run
gm task run --task-id "task_xxx"
gm task logs --task-id "task_xxx" --follow --interval 2s --timeout 5m
gm --yes task stop --task-id "task_xxx"
gm task model list --task-id "task_xxx"          # list checkpoints produced by training

# —— Project management ——
gm project list --page 1 --limit 50
gm project create --data '{"projectName":"my-project"}'

# —— Profile switching ——
gm config profile set dev --base-url "https://dev.example.com/prod-api" --timeout 30s
gm config profile use dev

# —— Output control ——
gm task list --human     # table
gm task list --quiet     # key fields only
```

Minimum fields for `task create`:

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

## Configuration

**Priority:** CLI flags > environment variables > config file

| Variable | Description |
|---|---|
| `GM_BASE_URL` | Service endpoint |
| `GM_API_KEY` | API key (temporary, not persisted) |
| `GM_TIMEOUT` | Request timeout, e.g. `30s` |
| `GM_RETRY` | Retry count |
| `GM_CONCURRENCY` | Concurrency limit |
| `GM_PROFILE` | Temporary profile |

Config file location:
- macOS / Linux: `~/.config/gradmotion/config.yaml`
- Windows: `%APPDATA%\gradmotion\config.yaml`

Example content (from an actual `gm auth login` run, with the system
keychain available — the API key itself stays out of the file):

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

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success |
| `1` | General error (network, server 5xx, unexpected) |
| `2` | Invalid argument (bad/missing flag, or a dangerous operation ran non-interactively without `--yes`) |
| `3` | Resource not found (404) |
| `4` | Permission denied (401 / 403) |
| `5` | Conflict / already exists (409) |
| `10` | Dry-run passed — safe to re-run without `--dry-run` |

## AI Agent Integration

The npm package bundles
[`skills/gm-cli/SKILL.md`](https://www.npmjs.com/package/@limxdynamics/gm-cli)
and `docs/GradmotionCLISkill-Guide.md`. The shipped guide documents
installing the skill for **Cursor's Agent Skills** (`~/.cursor/skills/gm-cli`);
the skill file's format may also work with other coding-agent skill
systems, though that isn't documented here. Key constraints baked into the
skill file:

- **Non-interactive detection**: when stdin isn't a terminal, dangerous
  operations (create/edit/run/stop/delete, etc.) don't hang waiting for
  confirmation — they fail immediately with exit code `2`. An agent
  running these **must** pass `--yes`.
- **Prefer `--dry-run` for write operations**: preview the request body
  and endpoint, then re-run without `--dry-run` once confirmed.
- Use relative paths for file arguments (e.g. `--file ./payload.json`),
  and clean up temporary JSON files after a successful call.
- All output is JSON by default, so an agent can parse it directly without
  matching human-readable text.

## License

[LICENSE](./LICENSE) at the root of this repository is **Apache License
2.0**.

---

<sub>This README was written from an actual download, unpack, and run of
`@limxdynamics/gm-cli@0.1.14` — not translated from another doc. If a
later version behaves differently, trust the live output of
`gm <command> --help` over this file.</sub>
