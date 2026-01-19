package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/fredericalix/whooktown-cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  "Show, set, and manage CLI configuration and contexts",
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUseContextCmd)
	configCmd.AddCommand(configGetContextsCmd)
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Mask token for display
		displayCfg := &config.Config{
			CurrentContext: cfg.CurrentContext,
			Contexts:       make(map[string]*config.Context),
		}

		for name, ctx := range cfg.Contexts {
			masked := ""
			if ctx.Token != "" {
				if len(ctx.Token) > 12 {
					masked = ctx.Token[:8] + "..." + ctx.Token[len(ctx.Token)-4:]
				} else {
					masked = "***"
				}
			}
			displayCfg.Contexts[name] = &config.Context{
				Name:          ctx.Name,
				Token:         masked,
				Environment:   ctx.Environment,
				DefaultLayout: ctx.DefaultLayout,
				AuthURL:       ctx.AuthURL,
				UIURL:         ctx.UIURL,
				SensorURL:     ctx.SensorURL,
				WorkflowURL:   ctx.WorkflowURL,
			}
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, displayCfg)
		}

		data, _ := yaml.Marshal(displayCfg)
		fmt.Printf("Config file: %s\n\n%s", cfgFile, string(data))
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value for the current context.

Available keys:
  environment    - PROD or DEV
  default_layout - default layout ID for commands
  auth_url       - custom auth service URL
  ui_url         - custom UI service URL
  sensor_url     - custom sensor service URL
  workflow_url   - custom workflow service URL`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]
		ctx := cfg.CurrentCtx()

		switch key {
		case "environment":
			if value != "PROD" && value != "DEV" {
				return fmt.Errorf("environment must be PROD or DEV")
			}
			ctx.Environment = value
		case "default_layout":
			ctx.DefaultLayout = value
		case "auth_url":
			ctx.AuthURL = value
		case "ui_url":
			ctx.UIURL = value
		case "sensor_url":
			ctx.SensorURL = value
		case "workflow_url":
			ctx.WorkflowURL = value
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}

		if err := cfg.Save(cfgFile); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Set %s = %s", key, value))
		return nil
	},
}

var configUseContextCmd = &cobra.Command{
	Use:   "use-context <name>",
	Short: "Switch to a different context",
	Long: `Switch to a different named context (environment).

If the context doesn't exist, it will be created.

Examples:
  wt config use-context dev     # Switch to dev environment
  wt config use-context prod    # Switch to prod environment`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Create context if it doesn't exist
		if _, ok := cfg.Contexts[name]; !ok {
			env := "PROD"
			if name == "dev" || name == "development" {
				env = "DEV"
			}
			cfg.Contexts[name] = &config.Context{
				Name:        name,
				Environment: env,
			}
		}

		cfg.CurrentContext = name

		if err := cfg.Save(cfgFile); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Switched to context: %s", name))
		return nil
	},
}

var configGetContextsCmd = &cobra.Command{
	Use:   "get-contexts",
	Short: "List all contexts",
	RunE: func(cmd *cobra.Command, args []string) error {
		if jsonOutput {
			return formatter.Format(os.Stdout, cfg.Contexts)
		}

		headers := []string{"NAME", "ENVIRONMENT", "CURRENT", "LOGGED IN"}
		rows := make([][]string, 0, len(cfg.Contexts))

		for name, ctx := range cfg.Contexts {
			current := ""
			if name == cfg.CurrentContext {
				current = "*"
			}
			loggedIn := "no"
			if ctx.Token != "" {
				loggedIn = "yes"
			}
			rows = append(rows, []string{name, ctx.Environment, current, loggedIn})
		}

		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}
