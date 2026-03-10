package auth

import (
	"strings"

	authstore "gradmotion-cli/internal/auth"
	"gradmotion-cli/internal/commands/shared"
	"gradmotion-cli/internal/config"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
	}

	cmd.AddCommand(
		newLoginCommand(),
		newLogoutCommand(),
		newStatusCommand(),
		newWhoamiCommand(),
	)

	return cmd
}

func newLoginCommand() *cobra.Command {
	var apiKey string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Save API key locally",
		RunE: func(_ *cobra.Command, _ []string) error {
			rt, err := shared.GetRuntime()
			if err != nil {
				return err
			}

			if strings.TrimSpace(apiKey) == "" {
				apiKey = rt.Profile.APIKey
			}
			if strings.TrimSpace(apiKey) == "" {
				return shared.EmitLocalError(
					"gm auth login",
					"INVALID_ARGUMENT",
					"api key is required",
					"use --api-key to pass your key",
				)
			}

			source := "keychain"
			store := authstore.NewStore()
			if err := store.Set(rt.ProfileName, apiKey); err != nil {
				source = "config"
			}

			if source == "config" {
				_ = rt.ConfigManager.UpdateCurrentProfile(func(p *config.Profile) error {
					p.APIKey = apiKey
					return nil
				})
				_ = rt.ConfigManager.Save()
			} else {
				_ = rt.ConfigManager.UpdateCurrentProfile(func(p *config.Profile) error {
					p.APIKey = ""
					return nil
				})
				_ = rt.ConfigManager.Save()
			}

			return shared.EmitLocalSuccess("gm auth login", map[string]any{
				"profile":  rt.ProfileName,
				"saved_to": source,
			})
		},
	}

	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key")
	return cmd
}

func newLogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove local API key",
		RunE: func(_ *cobra.Command, _ []string) error {
			rt, err := shared.GetRuntime()
			if err != nil {
				return err
			}
			store := authstore.NewStore()
			_ = store.Delete(rt.ProfileName)
			_ = rt.ConfigManager.UpdateCurrentProfile(func(p *config.Profile) error {
				p.APIKey = ""
				return nil
			})
			_ = rt.ConfigManager.Save()

			return shared.EmitLocalSuccess("gm auth logout", map[string]any{
				"profile": rt.ProfileName,
			})
		},
	}
}

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show local auth status",
		RunE: func(_ *cobra.Command, _ []string) error {
			rt, err := shared.GetRuntime()
			if err != nil {
				return err
			}

			keySource := "none"
			hasAPIKey := false
			if strings.TrimSpace(rt.Profile.APIKey) != "" {
				keySource = "config"
				hasAPIKey = true
			}

			if !hasAPIKey {
				store := authstore.NewStore()
				_, found, err := store.Get(rt.ProfileName)
				if err == nil && found {
					keySource = "keychain"
					hasAPIKey = true
				}
			}

			return shared.EmitLocalSuccess("gm auth status", map[string]any{
				"profile":     rt.ProfileName,
				"base_url":    rt.Profile.BaseURL,
				"has_api_key": hasAPIKey,
				"key_source":  keySource,
			})
		},
	}
}

func newWhoamiCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Get current user info from server",
		RunE: func(_ *cobra.Command, _ []string) error {
			return shared.CallAPI("gm auth whoami", "GET", "/user/me", nil, nil)
		},
	}
}
