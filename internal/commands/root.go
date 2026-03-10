package commands

import (
	"fmt"
	"os"
	"strings"

	authcmd "gradmotion-cli/internal/commands/auth"
	configcmd "gradmotion-cli/internal/commands/config"
	projectcmd "gradmotion-cli/internal/commands/project"
	"gradmotion-cli/internal/commands/shared"
	taskcmd "gradmotion-cli/internal/commands/task"
	"gradmotion-cli/internal/config"
	clilog "gradmotion-cli/internal/log"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type rootOptions struct {
	profile     string
	baseURL     string
	apiKey      string
	timeout     string
	retry       int
	concurrency int

	human   bool
	quiet   bool
	debug   bool
	yes     bool
	json    bool
	logFile string
}

func Execute(v, c, d string) error {
	version = v
	commit = c
	date = d
	return newRootCommand().Execute()
}

func newRootCommand() *cobra.Command {
	opts := &rootOptions{}

	cmd := &cobra.Command{
		Use:     "gm",
		Short:   "Gradmotion CLI",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			manager, err := config.NewManager("")
			if err != nil {
				return err
			}
			if err := manager.Load(); err != nil {
				return err
			}

			overrides := config.Overrides{
				ProfileName:        opts.profile,
				BaseURL:            opts.baseURL,
				APIKey:             opts.apiKey,
				Timeout:            opts.timeout,
				Retry:              opts.retry,
				Concurrency:        opts.concurrency,
				HasRetry:           cmd.Flags().Changed("retry"),
				HasConcurrency:     cmd.Flags().Changed("concurrency"),
				HasExplicitProfile: cmd.Flags().Changed("profile"),
			}

			if !overrides.HasExplicitProfile {
				if envProfile := strings.TrimSpace(os.Getenv("GM_PROFILE")); envProfile != "" {
					overrides.ProfileName = envProfile
					overrides.HasExplicitProfile = true
				}
			}

			profileName, profile := manager.EffectiveProfile(overrides)
			timeout, err := shared.ParseTimeout(profile.Timeout)
			if err != nil {
				return err
			}

			logger, closer, err := clilog.New(opts.logFile)
			if err != nil {
				return fmt.Errorf("init logger failed: %w", err)
			}
			if closer != nil {
				cobra.OnFinalize(func() {
					_ = closer.Close()
				})
			}

			rt := &shared.Runtime{
				ConfigManager: manager,
				ProfileName:   profileName,
				Profile:       profile,
				Timeout:       timeout,
				Retry:         profile.Retry,
				Concurrency:   profile.Concurrency,
				Human:         opts.human,
				Quiet:         opts.quiet,
				Debug:         opts.debug,
				ForceYes:      opts.yes,
			}
			rt.SetLogger(logger)
			shared.SetRuntime(rt)
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&opts.profile, "profile", "", "profile name")
	cmd.PersistentFlags().StringVar(&opts.baseURL, "base-url", "", "API base url")
	cmd.PersistentFlags().StringVar(&opts.apiKey, "api-key", "", "API key")
	cmd.PersistentFlags().StringVar(&opts.timeout, "timeout", "", "request timeout, e.g. 30s")
	cmd.PersistentFlags().IntVar(&opts.retry, "retry", 3, "request retry count")
	cmd.PersistentFlags().IntVar(&opts.concurrency, "concurrency", 4, "concurrency setting")
	cmd.PersistentFlags().BoolVar(&opts.human, "human", false, "human readable output")
	cmd.PersistentFlags().BoolVar(&opts.quiet, "quiet", false, "output key fields only")
	cmd.PersistentFlags().BoolVar(&opts.debug, "debug", false, "enable debug logs")
	cmd.PersistentFlags().BoolVar(&opts.yes, "yes", false, "skip confirmation prompts")
	cmd.PersistentFlags().BoolVar(&opts.json, "json", false, "force json output")
	cmd.PersistentFlags().StringVar(&opts.logFile, "log-file", "", "write logs to file")

	cmd.SetVersionTemplate(fmt.Sprintf(
		"{\n  \"success\": true,\n  \"data\": {\n    \"version\": %q,\n    \"commit\": %q,\n    \"date\": %q\n  }\n}\n",
		version, commit, date,
	))

	cmd.AddCommand(
		authcmd.NewCommand(),
		configcmd.NewCommand(),
		projectcmd.NewCommand(),
		taskcmd.NewCommand(),
	)

	return cmd
}
