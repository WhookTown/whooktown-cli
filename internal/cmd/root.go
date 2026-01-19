package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fredericalix/whooktown-cli/internal/config"
	"github.com/fredericalix/whooktown-cli/internal/output"
)

var (
	// Global flags
	cfgFile    string
	tokenFlag  string
	envFlag    string
	jsonOutput bool
	debugMode  bool

	// Shared state
	cfg       *config.Config
	formatter output.Formatter
)

var rootCmd = &cobra.Command{
	Use:   "wt",
	Short: "whooktown CLI - Manage your virtual IT city",
	Long: `wt is a command-line interface for whooktown,
a platform that visualizes IT infrastructure as a 3D virtual city.

Use wt to manage layouts, send sensor data, control camera and traffic,
and automate workflows.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for completion commands
		if cmd.Name() == "completion" || cmd.Name() == "__complete" {
			return nil
		}

		// Load configuration
		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return err
		}

		// Override environment from flag
		if envFlag != "" {
			switch envFlag {
			case "dev", "DEV", "development":
				cfg.CurrentCtx().Environment = "DEV"
			case "prod", "PROD", "production":
				cfg.CurrentCtx().Environment = "PROD"
			default:
				return fmt.Errorf("invalid environment: %s (use dev or prod)", envFlag)
			}
		}

		// Set up formatter
		if jsonOutput {
			formatter = output.New(output.FormatJSON)
		} else {
			formatter = output.New(output.FormatTable)
		}

		return nil
	},
	SilenceUsage: true, // Don't show usage on errors
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", config.DefaultConfigPath(),
		"config file path")
	rootCmd.PersistentFlags().StringVar(&tokenFlag, "token", "",
		"authentication token (overrides config)")
	rootCmd.PersistentFlags().StringVar(&envFlag, "env", "",
		"environment override (dev or prod)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false,
		"output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false,
		"enable debug output")

	// Add subcommands
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(layoutCmd)
	rootCmd.AddCommand(cameraCmd)
	rootCmd.AddCommand(trafficCmd)
	rootCmd.AddCommand(audioCmd)
	rootCmd.AddCommand(workflowCmd)
	rootCmd.AddCommand(sensorCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(completionCmd)
}
