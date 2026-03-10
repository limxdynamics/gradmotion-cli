# 训练任务模块 CLI 好用 MVP 接口清单

> 目标：在 `gm-cli` 中实现“可独立完成新建仿真训练 + 运行 + 观察 + 收尾”的好用 MVP。  
> 说明：本清单已将你强调的三点全部纳入必需项（算力资源、提醒设置、个人存储挂载路径）。

## 一、必须接口（建议 MVP 首批实现）

### 1) 项目域（任务归属前置）
- `POST /project/list`：项目列表（选择 `projectId`）
- `POST /project/create`：新建项目
- `GET /project/info/{projectId}`：项目校验/详情

### 2) 新建仿真训练（核心）
- `POST /task/create`：创建任务
- `POST /task/edit`：编辑任务草稿
- `POST /task/copy`：复制任务（提升创建效率）

### 3) 你指定必须包含的三点
- `GET /task/goods/list-by-category`：算力资源列表（用于选择 `goodsId`）
- 创建请求体字段 `runtimeReminderConfig`：提醒设置（随 `/task/create` 提交）
- `GET /gm/storage/list`：个人存储目录（用于 `personalDataPath` 数据挂载路径选择）

### 4) 镜像选择（创建前置）
- `GET /images/official/list`：官方镜像列表
- `GET /task/getImageVersion`：镜像版本列表
- `GET /images/personal/list`：个人镜像列表

### 5) 任务运行控制（闭环必需）
- `POST /task/run`：运行任务
- `POST /task/stop`：停止任务
- `POST /task/restart`：继续运行任务
- `POST /task/del`：删除任务

### 6) 任务查询与可观测（闭环必需）
- `POST /task/list`：任务列表/状态
- `GET /task/info/{taskId}`：任务详情
- `POST /task/console/log`：控制台日志
- `GET /task/run/env/{taskId}`：运行环境信息

### 7) 图表与结果（好用 MVP 必需）
- `GET /task/data/keys/{taskId}`：图表指标 key
- `POST /task/data/info`：按 key 查询图表数据
- `GET /task/data/download/{taskId}`：下载图表数据

### 8) 超参读取与编辑（好用 MVP 必需）
- `GET /task/hp/info/{taskId}`：读取超参
- `POST /task/hp/edit`：更新超参

## 二、增强接口（第二批）

### 1) 标签与笔记
- `POST /task/updateTag`
- `POST /task/getTaskNote`
- `POST /task/updateTaskNote`
- `GET /task/tag/list/{projectId}`

### 2) 批量操作
- `POST /task/batch/stop`
- `POST /task/batch/delete`

### 3) 任务对比与调优
- `POST /task/getCodeDiff`
- `GET /task-params-tuning/{taskId}`
- `POST /task-params-tuning/{taskId}`

## 三、建议的 CLI 命令能力映射

### 1) project 命令
- `gm project list`
- `gm project create`
- `gm project info`

### 2) task 命令（新增建议）
- `gm task copy`
- `gm task resource list`（映射 `/task/goods/list-by-category`）
- `gm task image official list`
- `gm task image personal list`
- `gm task image versions`
- `gm task storage list`（映射 `/gm/storage/list`）
- `gm task data keys`
- `gm task data get`
- `gm task data download`
- `gm task hp get`
- `gm task env get`

## 四、MVP 验收标准（最小）

- 能通过 CLI 完成：创建项目 -> 新建仿真训练任务 -> 运行 -> 看日志 -> 看状态 -> 停止/删除。
- 新建任务时可通过 CLI 选择：算力资源、镜像版本、提醒配置、个人存储挂载路径。
- 训练后可通过 CLI 获取：图表数据 key、图表数据、图表压缩包下载链接/文件。
