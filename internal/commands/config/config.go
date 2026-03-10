package config

import (
	"fmt"
	"sort"
	"strings"

	"gradmotion-cli/internal/commands/shared"
	"gradmotion-cli/internal/config"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration commands",
	}

	cmd.AddCommand(
		newSetCommand(),
		newGetCommand(),
		newProfileCommand(),
	)

	return cmd
}

func newSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set value in current profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			rt, err := shared.GetRuntime()
			if err != nil {
				return err
			}
			key := strings.ToLower(strings.TrimSpace(args[0]))
			value := strings.TrimSpace(args[1])

			err = rt.ConfigManager.UpdateCurrentProfile(func(p *config.Profile) error {
				switch key {
				case "base_url":
					p.BaseURL = value
				case "api_key":
					p.APIKey = value
				case "timeout":
					p.Timeout = value
				case "retry":
					var n int
					_, err := fmt.Sscanf(value, "%d", &n)
					if err != nil || n <= 0 {
						return fmt.Errorf("retry must be positive integer")
					}
					p.Retry = n
				case "concurrency":
					var n int
					_, err := fmt.Sscanf(value, "%d", &n)
					if err != nil || n <= 0 {
						return fmt.Errorf("concurrency must be positive integer")
					}
					p.Concurrency = n
				default:
					return fmt.Errorf("unsupported key %q", key)
				}
				return nil
			})
			if err != nil {
				return shared.EmitLocalError("gm config set", "INVALID_ARGUMENT", err.Error(), "")
			}

			if err := rt.ConfigManager.Save(); err != nil {
				return shared.EmitLocalError("gm config set", "CONFIG_SAVE_FAILED", err.Error(), "")
			}

			return shared.EmitLocalSuccess("gm config set", map[string]any{
				"profile": rt.ProfileName,
				"key":     key,
				"value":   value,
			})
		},
	}
}

func newGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get value from current profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			rt, err := shared.GetRuntime()
			if err != nil {
				return err
			}
			key := strings.ToLower(strings.TrimSpace(args[0]))
			p := rt.Profile

			var value any
			switch key {
			case "base_url":
				value = p.BaseURL
			case "api_key":
				value = p.APIKey
			case "timeout":
				value = p.Timeout
			case "retry":
				value = p.Retry
			case "concurrency":
				value = p.Concurrency
			default:
				return shared.EmitLocalError("gm config get", "INVALID_ARGUMENT", "unsupported key", key)
			}

			return shared.EmitLocalSuccess("gm config get", map[string]any{
				"profile": rt.ProfileName,
				"key":     key,
				"value":   value,
			})
		},
	}
}

func newProfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Profile management",
	}
	cmd.AddCommand(
		newProfileListCommand(),
		newProfileUseCommand(),
		newProfileSetCommand(),
	)
	return cmd
}

func newProfileListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List profiles",
		RunE: func(_ *cobra.Command, _ []string) error {
			rt, err := shared.GetRuntime()
			if err != nil {
				return err
			}
			profiles := rt.ConfigManager.ListProfiles()
			names := make([]string, 0, len(profiles))
			for name := range profiles {
				names = append(names, name)
			}
			sort.Strings(names)

			items := make([]map[string]any, 0, len(names))
			for _, name := range names {
				p := profiles[name]
				items = append(items, map[string]any{
					"name":        name,
					"current":     name == rt.ConfigManager.CurrentProfileName(),
					"base_url":    p.BaseURL,
					"timeout":     p.Timeout,
					"retry":       p.Retry,
					"concurrency": p.Concurrency,
				})
			}

			return shared.EmitLocalSuccess("gm config profile list", map[string]any{
				"profiles": items,
			})
		},
	}
}

func newProfileUseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Switch current profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			rt, err := shared.GetRuntime()
			if err != nil {
				return err
			}

			name := strings.TrimSpace(args[0])
			if err := rt.ConfigManager.SetCurrentProfileName(name); err != nil {
				return shared.EmitLocalError("gm config profile use", "INVALID_ARGUMENT", err.Error(), "")
			}
			if err := rt.ConfigManager.Save(); err != nil {
				return shared.EmitLocalError("gm config profile use", "CONFIG_SAVE_FAILED", err.Error(), "")
			}
			return shared.EmitLocalSuccess("gm config profile use", map[string]any{
				"current": name,
			})
		},
	}
}

func newProfileSetCommand() *cobra.Command {
	var (
		baseURL     string
		apiKey      string
		timeout     string
		retry       int
		concurrency int
	)

	cmd := &cobra.Command{
		Use:   "set <name>",
		Short: "Create or update a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := shared.GetRuntime()
			if err != nil {
				return err
			}
			name := strings.TrimSpace(args[0])

			p, ok := rt.ConfigManager.GetProfile(name)
			if !ok {
				p = config.Profile{
					Timeout:     "30s",
					Retry:       3,
					Concurrency: 4,
				}
			}
			if cmd.Flags().Changed("base-url") {
				p.BaseURL = strings.TrimSpace(baseURL)
			}
			if cmd.Flags().Changed("api-key") {
				p.APIKey = strings.TrimSpace(apiKey)
			}
			if cmd.Flags().Changed("timeout") {
				p.Timeout = strings.TrimSpace(timeout)
			}
			if cmd.Flags().Changed("retry") {
				p.Retry = retry
			}
			if cmd.Flags().Changed("concurrency") {
				p.Concurrency = concurrency
			}

			rt.ConfigManager.UpsertProfile(name, p)
			if err := rt.ConfigManager.Save(); err != nil {
				return shared.EmitLocalError("gm config profile set", "CONFIG_SAVE_FAILED", err.Error(), "")
			}

			return shared.EmitLocalSuccess("gm config profile set", map[string]any{
				"profile": name,
			})
		},
	}

	cmd.Flags().StringVar(&baseURL, "base-url", "", "base url")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "api key")
	cmd.Flags().StringVar(&timeout, "timeout", "", "timeout")
	cmd.Flags().IntVar(&retry, "retry", 3, "retry")
	cmd.Flags().IntVar(&concurrency, "concurrency", 4, "concurrency")
	return cmd
}
