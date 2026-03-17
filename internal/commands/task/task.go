package task

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gradmotion-cli/internal/commands/shared"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Task commands",
	}

	cmd.AddCommand(
		newCreateCommand(),
		newEditCommand(),
		newCopyCommand(),
		newListCommand(),
		newInfoCommand(),
		newModelCommand(),
		newRunCommand(),
		newStopCommand(),
		newDeleteCommand(),
		newLogsCommand(),
		newParamsCommand(),
		newResourceCommand(),
		newImageCommand(),
		newStorageCommand(),
		newDataCommand(),
		newHPCommand(),
		newEnvCommand(),
		newTagCommand(),
		newBatchCommand(),
	)
	return cmd
}

func newCreateCommand() *cobra.Command {
	var bodyFlags rawBodyFlags
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create task",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm task create", "INVALID_ARGUMENT", err.Error(), "use --data or --file with JSON")
			}
			return shared.CallAPI("gm task create", "POST", "/task/create", body, nil)
		},
	}
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newEditCommand() *cobra.Command {
	var bodyFlags rawBodyFlags
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit task",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm task edit", "INVALID_ARGUMENT", err.Error(), "use --data or --file with JSON")
			}
			return shared.CallAPI("gm task edit", "POST", "/task/edit", body, nil)
		},
	}
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newListCommand() *cobra.Command {
	var (
		page      int
		limit     int
		bodyFlags rawBodyFlags
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm task list", "INVALID_ARGUMENT", err.Error(), "")
			}
			if body == nil {
				body = map[string]any{}
			}
			body["pageNum"] = page
			body["pageSize"] = limit
			body["page_num"] = page
			body["page_size"] = limit
			return shared.CallAPI("gm task list", "POST", "/task/list", body, nil)
		},
	}
	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().IntVar(&limit, "limit", 50, "page size")
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newInfoCommand() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get task info",
		RunE: func(_ *cobra.Command, _ []string) error {
			if strings.TrimSpace(taskID) == "" {
				return shared.EmitLocalError("gm task info", "INVALID_ARGUMENT", "task-id is required", "")
			}
			endpoint := "/task/info/" + url.PathEscape(taskID)
			return shared.CallAPI("gm task info", "GET", endpoint, nil, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	return cmd
}

func newModelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "model",
		Short: "Task model commands",
	}
	cmd.AddCommand(newModelListCommand())
	return cmd
}

func newModelListCommand() *cobra.Command {
	var (
		taskID     string
		checkpoint string
		pageNum    int
		pageSize   int
		bodyFlags  rawBodyFlags
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List task model checkpoints",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm task model list", "INVALID_ARGUMENT", err.Error(), "use --data or --file with JSON")
			}
			if body == nil {
				body = map[string]any{}
			}
			if strings.TrimSpace(taskID) != "" {
				body["taskId"] = strings.TrimSpace(taskID)
				body["task_id"] = strings.TrimSpace(taskID)
			}
			if strings.TrimSpace(checkpoint) != "" {
				body["checkpoint"] = strings.TrimSpace(checkpoint)
			}
			body["pageNum"] = pageNum
			body["pageSize"] = pageSize
			body["page_num"] = pageNum
			body["page_size"] = pageSize
			if body["taskId"] == nil && body["task_id"] == nil {
				return shared.EmitLocalError("gm task model list", "INVALID_ARGUMENT", "task-id is required", "use --task-id or --file/--data with taskId")
			}
			return shared.CallAPI("gm task model list", "POST", "/task/model/info", body, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	cmd.Flags().StringVar(&checkpoint, "checkpoint", "", "checkpoint filter")
	cmd.Flags().IntVar(&pageNum, "page-num", 1, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "page size")
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newRunCommand() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run task",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := taskIDBody(taskID)
			if err != nil {
				return shared.EmitLocalError("gm task run", "INVALID_ARGUMENT", err.Error(), "")
			}
			return shared.CallAPI("gm task run", "POST", "/task/run", body, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	return cmd
}

func newStopCommand() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop task",
		RunE: func(_ *cobra.Command, _ []string) error {
			if !confirmDangerous("Stop task?", "gm task stop") {
				return nil
			}
			body, err := taskIDBody(taskID)
			if err != nil {
				return shared.EmitLocalError("gm task stop", "INVALID_ARGUMENT", err.Error(), "")
			}
			return shared.CallAPI("gm task stop", "POST", "/task/stop", body, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	return cmd
}

func newDeleteCommand() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete task",
		RunE: func(_ *cobra.Command, _ []string) error {
			if !confirmDangerous("Delete task?", "gm task delete") {
				return nil
			}
			body, err := taskIDBody(taskID)
			if err != nil {
				return shared.EmitLocalError("gm task delete", "INVALID_ARGUMENT", err.Error(), "")
			}
			return shared.CallAPI("gm task delete", "POST", "/task/del", body, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	return cmd
}

func newLogsCommand() *cobra.Command {
	var (
		taskID        string
		follow        bool
		interval      time.Duration
		timeout       time.Duration
		raw          bool
		noRequestLog bool
	)

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get task logs",
		RunE: func(_ *cobra.Command, _ []string) error {
			if strings.TrimSpace(taskID) == "" {
				return shared.EmitLocalError("gm task logs", "INVALID_ARGUMENT", "task-id is required", "")
			}
			body := map[string]any{
				"task_id": taskID,
			}

			if raw {
				stdout := os.Stdout
				if !follow {
					return shared.CallAPIRawOutput("gm task logs", "POST", "/task/console/log", body, nil, false, stdout, noRequestLog)
				}
				if interval <= 0 {
					interval = 2 * time.Second
				}
				deadline := time.Time{}
				if timeout > 0 {
					deadline = time.Now().Add(timeout)
				}
				for {
					if err := shared.CallAPIRawOutput("gm task logs", "POST", "/task/console/log", body, nil, false, stdout, noRequestLog); err != nil {
						return err
					}
					if !deadline.IsZero() && time.Now().After(deadline) {
						break
					}
					time.Sleep(interval)
				}
				return nil
			}

			if !follow {
				return shared.CallAPIWithNoRequestLog("gm task logs", "POST", "/task/console/log", body, nil, noRequestLog)
			}
			if interval <= 0 {
				interval = 2 * time.Second
			}
			deadline := time.Time{}
			if timeout > 0 {
				deadline = time.Now().Add(timeout)
			}
			for {
				if err := shared.CallAPIWithNoRequestLog("gm task logs", "POST", "/task/console/log", body, nil, noRequestLog); err != nil {
					return err
				}
				if !deadline.IsZero() && time.Now().After(deadline) {
					break
				}
				time.Sleep(interval)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	cmd.Flags().BoolVar(&follow, "follow", false, "follow log updates")
	cmd.Flags().DurationVar(&interval, "interval", 2*time.Second, "poll interval when --follow")
	cmd.Flags().DurationVar(&timeout, "timeout", 0, "max follow time, 0 means no limit")
	cmd.Flags().BoolVar(&raw, "raw", false, "output only log content from data to stdout, no JSON envelope")
	cmd.Flags().BoolVar(&noRequestLog, "no-request-log", false, "do not output request metadata to stderr for this command")
	return cmd
}

func newParamsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Task params commands",
	}
	cmd.AddCommand(
		newParamsSubmitCommand(),
		newParamsUpdateCommand(),
	)
	return cmd
}

func newParamsSubmitCommand() *cobra.Command {
	var (
		taskID    string
		bodyFlags rawBodyFlags
	)
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "Submit task params",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm task params submit", "INVALID_ARGUMENT", err.Error(), "")
			}
			if body == nil {
				body = map[string]any{}
			}
			if strings.TrimSpace(taskID) != "" {
				body["task_id"] = taskID
			}
			return shared.CallAPI("gm task params submit", "POST", "/task/hp/up", body, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newParamsUpdateCommand() *cobra.Command {
	var (
		taskID    string
		bodyFlags rawBodyFlags
	)
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update task params",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm task params update", "INVALID_ARGUMENT", err.Error(), "")
			}
			if body == nil {
				body = map[string]any{}
			}
			if strings.TrimSpace(taskID) != "" {
				body["task_id"] = taskID
			}
			return shared.CallAPI("gm task params update", "POST", "/task/hp/edit", body, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newBatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Batch task operations",
	}
	cmd.AddCommand(
		newBatchStopCommand(),
		newBatchDeleteCommand(),
	)
	return cmd
}

func newBatchStopCommand() *cobra.Command {
	var taskIDs []string
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Batch stop tasks",
		RunE: func(_ *cobra.Command, _ []string) error {
			if !confirmDangerous("Batch stop tasks?", "gm task batch stop") {
				return nil
			}
			ids, err := normalizeTaskIDs(taskIDs)
			if err != nil {
				return shared.EmitLocalError("gm task batch stop", "INVALID_ARGUMENT", err.Error(), "")
			}
			body := map[string]any{"task_ids": ids}
			return shared.CallAPI("gm task batch stop", "POST", "/task/batch/stop", body, nil)
		},
	}
	cmd.Flags().StringSliceVar(&taskIDs, "task-ids", nil, "comma-separated task ids")
	return cmd
}

func newBatchDeleteCommand() *cobra.Command {
	var taskIDs []string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Batch delete tasks",
		RunE: func(_ *cobra.Command, _ []string) error {
			if !confirmDangerous("Batch delete tasks?", "gm task batch delete") {
				return nil
			}
			ids, err := normalizeTaskIDs(taskIDs)
			if err != nil {
				return shared.EmitLocalError("gm task batch delete", "INVALID_ARGUMENT", err.Error(), "")
			}
			body := map[string]any{"task_ids": ids}
			return shared.CallAPI("gm task batch delete", "POST", "/task/batch/delete", body, nil)
		},
	}
	cmd.Flags().StringSliceVar(&taskIDs, "task-ids", nil, "comma-separated task ids")
	return cmd
}

type rawBodyFlags struct {
	data string
	file string
}

func addRawBodyFlags(cmd *cobra.Command, flags *rawBodyFlags) {
	cmd.Flags().StringVar(&flags.data, "data", "", "request body as JSON string")
	cmd.Flags().StringVar(&flags.file, "file", "", "path to JSON request body file")
}

func parseRawBody(flags rawBodyFlags) (map[string]any, error) {
	if strings.TrimSpace(flags.data) == "" && strings.TrimSpace(flags.file) == "" {
		return nil, nil
	}
	if strings.TrimSpace(flags.data) != "" && strings.TrimSpace(flags.file) != "" {
		return nil, errors.New("use either --data or --file, not both")
	}

	raw := ""
	if strings.TrimSpace(flags.data) != "" {
		raw = flags.data
	} else {
		b, err := os.ReadFile(flags.file)
		if err != nil {
			return nil, fmt.Errorf("read body file failed: %w", err)
		}
		raw = string(b)
	}

	result := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("invalid json body: %w", err)
	}
	return result, nil
}

func taskIDBody(taskID string) (map[string]any, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, errors.New("task-id is required")
	}
	return map[string]any{"task_id": taskID}, nil
}

func normalizeTaskIDs(values []string) ([]string, error) {
	if len(values) == 0 {
		return nil, errors.New("task-ids is required")
	}
	out := make([]string, 0, len(values))
	for _, v := range values {
		for _, part := range strings.Split(v, ",") {
			p := strings.TrimSpace(part)
			if p != "" {
				out = append(out, p)
			}
		}
	}
	if len(out) == 0 {
		return nil, errors.New("task-ids is required")
	}
	return out, nil
}

func confirmDangerous(message, command string) bool {
	rt, err := shared.GetRuntime()
	if err != nil {
		return false
	}
	if rt.ForceYes {
		return true
	}
	_, _ = fmt.Printf("%s [y/N]: ", message)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	if line == "y" || line == "yes" {
		return true
	}
	_ = shared.EmitLocalSuccess(command, map[string]any{
		"skipped": true,
	})
	return false
}

func newCopyCommand() *cobra.Command {
	var bodyFlags rawBodyFlags
	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy task",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm task copy", "INVALID_ARGUMENT", err.Error(), "use --data or --file with JSON")
			}
			return shared.CallAPI("gm task copy", "POST", "/task/copy", body, nil)
		},
	}
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newResourceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Task resource commands",
	}
	cmd.AddCommand(newResourceListCommand())
	return cmd
}

func newResourceListCommand() *cobra.Command {
	var (
		goodsBackCategory int
		pageNum           int
		pageSize          int
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List task resources",
		RunE: func(_ *cobra.Command, _ []string) error {
			if goodsBackCategory <= 0 {
				return shared.EmitLocalError("gm task resource list", "INVALID_ARGUMENT", "goods-back-category is required", "")
			}
			query := map[string]string{
				"goodsBackCategory": strconv.Itoa(goodsBackCategory),
				"pageNum":           strconv.Itoa(pageNum),
				"pageSize":          strconv.Itoa(pageSize),
			}
			return shared.CallAPI("gm task resource list", "GET", "/task/goods/list-by-category", nil, query)
		},
	}
	cmd.Flags().IntVar(&goodsBackCategory, "goods-back-category", 0, "goods back category (3=train, 4=dev)")
	cmd.Flags().IntVar(&pageNum, "page-num", 1, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "page size")
	return cmd
}

func newImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Task image commands",
	}
	cmd.AddCommand(
		newImageOfficialCommand(),
		newImagePersonalCommand(),
		newImageVersionsCommand(),
	)
	return cmd
}

func newImageOfficialCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "official",
		Short: "List official images",
		RunE: func(_ *cobra.Command, _ []string) error {
			return shared.CallAPI("gm task image official", "GET", "/images/official/list", nil, nil)
		},
	}
}

func newImagePersonalCommand() *cobra.Command {
	var (
		versionStatus string
		pageNum       int
		pageSize      int
	)
	cmd := &cobra.Command{
		Use:   "personal",
		Short: "List personal images",
		RunE: func(_ *cobra.Command, _ []string) error {
			query := map[string]string{
				"versionStatus": versionStatus,
				"pageNum":       strconv.Itoa(pageNum),
				"pageSize":      strconv.Itoa(pageSize),
			}
			return shared.CallAPI("gm task image personal", "GET", "/images/personal/list", nil, query)
		},
	}
	cmd.Flags().StringVar(&versionStatus, "version-status", "1", "version status")
	cmd.Flags().IntVar(&pageNum, "page-num", 1, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 50, "page size")
	return cmd
}

func newImageVersionsCommand() *cobra.Command {
	var imageID string
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "List image versions",
		RunE: func(_ *cobra.Command, _ []string) error {
			imageID = strings.TrimSpace(imageID)
			if imageID == "" {
				return shared.EmitLocalError("gm task image versions", "INVALID_ARGUMENT", "image-id is required", "")
			}
			query := map[string]string{"imageId": imageID}
			return shared.CallAPI("gm task image versions", "GET", "/task/getImageVersion", nil, query)
		},
	}
	cmd.Flags().StringVar(&imageID, "image-id", "", "image id")
	return cmd
}

func newStorageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "Task storage commands",
	}
	cmd.AddCommand(newStorageListCommand())
	return cmd
}

func newStorageListCommand() *cobra.Command {
	var folderPath string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List personal storage files",
		RunE: func(_ *cobra.Command, _ []string) error {
			query := map[string]string{"folderPath": strings.TrimSpace(folderPath)}
			return shared.CallAPIAbsolute("gm task storage list", "GET", "/gm/storage/list", nil, query)
		},
	}
	cmd.Flags().StringVar(&folderPath, "folder-path", "", "folder path")
	return cmd
}

func newDataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "Task chart data commands",
	}
	cmd.AddCommand(
		newDataKeysCommand(),
		newDataGetCommand(),
		newDataDownloadCommand(),
	)
	return cmd
}

func newDataKeysCommand() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Get task chart keys",
		RunE: func(_ *cobra.Command, _ []string) error {
			taskID = strings.TrimSpace(taskID)
			if taskID == "" {
				return shared.EmitLocalError("gm task data keys", "INVALID_ARGUMENT", "task-id is required", "")
			}
			endpoint := "/task/data/keys/" + url.PathEscape(taskID)
			return shared.CallAPI("gm task data keys", "GET", endpoint, nil, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	return cmd
}

func newDataGetCommand() *cobra.Command {
	var (
		taskID        string
		dataKey       string
		endTime       string
		session       string
		samplingMode  string
		maxDataPoints int
		bodyFlags     rawBodyFlags
	)
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get task chart data",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm task data get", "INVALID_ARGUMENT", err.Error(), "use --data or --file with JSON")
			}
			if body == nil {
				body = map[string]any{}
			}
			if strings.TrimSpace(taskID) != "" {
				body["task_id"] = strings.TrimSpace(taskID)
			}
			if strings.TrimSpace(dataKey) != "" {
				body["data_key"] = strings.TrimSpace(dataKey)
			}
			if strings.TrimSpace(endTime) != "" {
				body["end_time"] = strings.TrimSpace(endTime)
			}
			if strings.TrimSpace(session) != "" {
				body["session"] = strings.TrimSpace(session)
			}
			if strings.TrimSpace(samplingMode) != "" {
				body["sampling_mode"] = strings.TrimSpace(samplingMode)
			}
			if maxDataPoints > 0 {
				body["max_data_points"] = maxDataPoints
			}
			if strings.TrimSpace(taskID) == "" || strings.TrimSpace(dataKey) == "" {
				if _, ok := body["task_id"]; !ok {
					return shared.EmitLocalError("gm task data get", "INVALID_ARGUMENT", "task-id is required", "")
				}
				if _, ok := body["data_key"]; !ok {
					return shared.EmitLocalError("gm task data get", "INVALID_ARGUMENT", "data-key is required", "")
				}
			}
			return shared.CallAPI("gm task data get", "POST", "/task/data/info", body, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	cmd.Flags().StringVar(&dataKey, "data-key", "", "chart data key")
	cmd.Flags().StringVar(&endTime, "end-time", "", "end time")
	cmd.Flags().StringVar(&session, "session", "", "session id")
	cmd.Flags().StringVar(&samplingMode, "sampling-mode", "", "sampling mode: precise|accelerate")
	cmd.Flags().IntVar(&maxDataPoints, "max-data-points", 0, "max sampled data points")
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newDataDownloadCommand() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Get task chart download link",
		RunE: func(_ *cobra.Command, _ []string) error {
			taskID = strings.TrimSpace(taskID)
			if taskID == "" {
				return shared.EmitLocalError("gm task data download", "INVALID_ARGUMENT", "task-id is required", "")
			}
			endpoint := "/task/data/download/" + url.PathEscape(taskID)
			return shared.CallAPI("gm task data download", "GET", endpoint, nil, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	return cmd
}

func newHPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hp",
		Short: "Task hyperparameter commands",
	}
	cmd.AddCommand(newHPGetCommand())
	return cmd
}

func newHPGetCommand() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get task hyperparameters",
		RunE: func(_ *cobra.Command, _ []string) error {
			taskID = strings.TrimSpace(taskID)
			if taskID == "" {
				return shared.EmitLocalError("gm task hp get", "INVALID_ARGUMENT", "task-id is required", "")
			}
			endpoint := "/task/hp/info/" + url.PathEscape(taskID)
			return shared.CallAPI("gm task hp get", "GET", endpoint, nil, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	return cmd
}

func newEnvCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Task runtime environment commands",
	}
	cmd.AddCommand(newEnvGetCommand())
	return cmd
}

func newEnvGetCommand() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get task runtime environment",
		RunE: func(_ *cobra.Command, _ []string) error {
			taskID = strings.TrimSpace(taskID)
			if taskID == "" {
				return shared.EmitLocalError("gm task env get", "INVALID_ARGUMENT", "task-id is required", "")
			}
			endpoint := "/task/run/env/" + url.PathEscape(taskID)
			return shared.CallAPI("gm task env get", "GET", endpoint, nil, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	return cmd
}

func newTagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Task tag commands (update / get task tags, list user tags)",
	}
	cmd.AddCommand(
		newTagUpdateCommand(),
		newTagGetCommand(),
		newTagListCommand(),
	)
	return cmd
}

func newTagUpdateCommand() *cobra.Command {
	var (
		taskID  string
		tagList []string
		bodyFlags rawBodyFlags
	)
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update task tags",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm task tag update", "INVALID_ARGUMENT", err.Error(), "use --data or --file with JSON, or --task-id and --tags")
			}
			if body == nil {
				body = map[string]any{}
			}
			if strings.TrimSpace(taskID) != "" {
				body["taskId"] = strings.TrimSpace(taskID)
				body["task_id"] = strings.TrimSpace(taskID)
			}
			if len(tagList) > 0 {
				body["taskTag"] = tagList
				body["task_tag"] = tagList
			}
			if body["taskId"] == nil && body["task_id"] == nil {
				return shared.EmitLocalError("gm task tag update", "INVALID_ARGUMENT", "task-id is required", "use --task-id or --file/--data with taskId")
			}
			return shared.CallAPI("gm task tag update", "POST", "/task/updateTag", body, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	cmd.Flags().StringSliceVar(&tagList, "tags", nil, "tag list (comma-separated or repeat)")
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newTagGetCommand() *cobra.Command {
	var taskID string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get task tags",
		RunE: func(_ *cobra.Command, _ []string) error {
			taskID = strings.TrimSpace(taskID)
			if taskID == "" {
				return shared.EmitLocalError("gm task tag get", "INVALID_ARGUMENT", "task-id is required", "")
			}
			body := map[string]any{"taskId": taskID, "task_id": taskID}
			return shared.CallAPI("gm task tag get", "POST", "/task/getTaskTag", body, nil)
		},
	}
	cmd.Flags().StringVar(&taskID, "task-id", "", "task id")
	return cmd
}

func newTagListCommand() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List user historical tags",
		RunE: func(_ *cobra.Command, _ []string) error {
			if limit <= 0 {
				limit = 200
			}
			query := map[string]string{"limit": strconv.Itoa(limit)}
			return shared.CallAPI("gm task tag list", "POST", "/task/getUserTag", nil, query)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 200, "max number of tags to return")
	return cmd
}
