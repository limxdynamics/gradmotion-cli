package project

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gradmotion-cli/internal/commands/shared"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project commands",
	}

	cmd.AddCommand(
		newListCommand(),
		newCreateCommand(),
		newEditCommand(),
		newDeleteCommand(),
		newInfoCommand(),
	)
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
		Short: "List projects",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm project list", "INVALID_ARGUMENT", err.Error(), "")
			}
			if body == nil {
				body = map[string]any{}
			}
			body["pageNum"] = page
			body["pageSize"] = limit
			body["page_num"] = page
			body["page_size"] = limit
			return shared.CallAPI("gm project list", "POST", "/project/list", body, nil)
		},
	}
	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().IntVar(&limit, "limit", 50, "page size")
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newCreateCommand() *cobra.Command {
	var bodyFlags rawBodyFlags
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create project",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm project create", "INVALID_ARGUMENT", err.Error(), "use --data or --file with JSON")
			}
			return shared.CallAPI("gm project create", "POST", "/project/create", body, nil)
		},
	}
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newEditCommand() *cobra.Command {
	var bodyFlags rawBodyFlags
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit project (e.g. rename)",
		RunE: func(_ *cobra.Command, _ []string) error {
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm project edit", "INVALID_ARGUMENT", err.Error(), "use --data or --file with JSON")
			}
			if body == nil {
				return shared.EmitLocalError("gm project edit", "INVALID_ARGUMENT", "request body is required (project_id and project_name)", "use --data or --file with JSON")
			}
			return shared.CallAPI("gm project edit", "POST", "/project/edit", body, nil)
		},
	}
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func newDeleteCommand() *cobra.Command {
	var (
		projectID string
		bodyFlags rawBodyFlags
	)
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete project",
		RunE: func(_ *cobra.Command, _ []string) error {
			if !confirmProjectDelete("gm project delete") {
				return nil
			}
			body, err := parseRawBody(bodyFlags)
			if err != nil {
				return shared.EmitLocalError("gm project delete", "INVALID_ARGUMENT", err.Error(), "use --data or --file with JSON, or --project-id")
			}
			if body == nil {
				body = map[string]any{}
			}
			if strings.TrimSpace(projectID) != "" {
				body["projectId"] = strings.TrimSpace(projectID)
				body["project_id"] = strings.TrimSpace(projectID)
			}
			if body["projectId"] == nil && body["project_id"] == nil {
				return shared.EmitLocalError("gm project delete", "INVALID_ARGUMENT", "project-id is required", "use --project-id or --file/--data with projectId")
			}
			return shared.CallAPI("gm project delete", "POST", "/project/del", body, nil)
		},
	}
	cmd.Flags().StringVar(&projectID, "project-id", "", "project id to delete")
	addRawBodyFlags(cmd, &bodyFlags)
	return cmd
}

func confirmProjectDelete(command string) bool {
	rt, err := shared.GetRuntime()
	if err != nil {
		return false
	}
	if rt.ForceYes {
		return true
	}
	fmt.Print("Delete project (and its tasks)? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	if line == "y" || line == "yes" {
		return true
	}
	_ = shared.EmitLocalSuccess(command, map[string]any{"skipped": true})
	return false
}

func newInfoCommand() *cobra.Command {
	var projectID string
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get project info",
		RunE: func(_ *cobra.Command, _ []string) error {
			projectID = strings.TrimSpace(projectID)
			if projectID == "" {
				return shared.EmitLocalError("gm project info", "INVALID_ARGUMENT", "project-id is required", "")
			}
			endpoint := "/project/info/" + url.PathEscape(projectID)
			return shared.CallAPI("gm project info", "GET", endpoint, nil, nil)
		},
	}
	cmd.Flags().StringVar(&projectID, "project-id", "", "project id")
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
