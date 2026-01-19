package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
	"github.com/fredericalix/whooktown-cli/internal/client"
)

var layoutCmd = &cobra.Command{
	Use:   "layout",
	Short: "Manage city layouts",
	Long:  "Create, update, delete, and list city layouts",
}

func init() {
	layoutCmd.AddCommand(layoutListCmd)
	layoutCmd.AddCommand(layoutShowCmd)
	layoutCmd.AddCommand(layoutCreateCmd)
	layoutCmd.AddCommand(layoutUpdateCmd)
	layoutCmd.AddCommand(layoutDeleteCmd)
	layoutCmd.AddCommand(layoutQuotaCmd)
	layoutCmd.AddCommand(layoutArchiveCmd)
}

var layoutListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all layouts",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		// Get quota (which includes layout count info)
		quota, err := c.UI.GetQuota(ctx)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, quota)
		}

		fmt.Printf("Plan: %s (Status: %s)\n", quota.Plan, quota.Status)
		fmt.Printf("Layouts: %d/%d\n", quota.Layouts.Used, quota.Layouts.Max)
		fmt.Printf("Archived: %d\n", quota.Layouts.Archived)
		fmt.Printf("Assets per layout: %d max\n", quota.AssetsPerLayout.Max)

		return nil
	},
}

var layoutShowCmd = &cobra.Command{
	Use:   "show <layout_id>",
	Short: "Show layout details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Note: Need to add GetLayout to SDK or use direct API call
		return fmt.Errorf("not yet implemented - requires SDK extension for GetLayout")
	},
}

var (
	layoutFile string
	layoutName string
	gridWidth  int
	gridHeight int
)

var layoutCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new layout",
	Long: `Create a new layout from a JSON file or inline flags.

Example with file:
  wt layout create -f layout.json

Example with flags:
  wt layout create --name "My City" --grid-width 10 --grid-height 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		var layout whooktown.Layout

		if layoutFile != "" {
			// Read from file
			data, err := os.ReadFile(layoutFile)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			if err := json.Unmarshal(data, &layout); err != nil {
				return fmt.Errorf("invalid JSON: %w", err)
			}
		} else if layoutName != "" {
			// Create from flags
			layout = whooktown.Layout{
				Name: layoutName,
				Grid: whooktown.Grid{Width: gridWidth, Height: gridHeight},
			}
		} else {
			return fmt.Errorf("either --file or --name is required")
		}

		result, err := c.UI.CreateLayout(ctx, &layout)
		if err != nil {
			return fmt.Errorf("failed to create layout: %w", err)
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, result)
		}

		formatter.Success(fmt.Sprintf("Layout created: %s", result.LayoutID))
		return nil
	},
}

func init() {
	layoutCreateCmd.Flags().StringVarP(&layoutFile, "file", "f", "", "JSON file with layout definition")
	layoutCreateCmd.Flags().StringVar(&layoutName, "name", "", "layout name")
	layoutCreateCmd.Flags().IntVar(&gridWidth, "grid-width", 10, "grid width")
	layoutCreateCmd.Flags().IntVar(&gridHeight, "grid-height", 10, "grid height")
}

var layoutUpdateCmd = &cobra.Command{
	Use:   "update <layout_id>",
	Short: "Update a layout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		layoutID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid layout ID: %w", err)
		}

		if layoutFile == "" {
			return fmt.Errorf("--file is required")
		}

		// Read from file
		data, err := os.ReadFile(layoutFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		var layout whooktown.Layout
		if err := json.Unmarshal(data, &layout); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		// Set the layout ID from args
		layout.ID = layoutID

		result, err := c.UI.CreateLayout(ctx, &layout) // CreateLayout does upsert
		if err != nil {
			return fmt.Errorf("failed to update layout: %w", err)
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, result)
		}

		formatter.Success(fmt.Sprintf("Layout updated: %s", result.LayoutID))
		return nil
	},
}

func init() {
	layoutUpdateCmd.Flags().StringVarP(&layoutFile, "file", "f", "", "JSON file with layout definition")
}

var layoutDeleteCmd = &cobra.Command{
	Use:   "delete <layout_id>",
	Short: "Delete a layout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		layoutID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid layout ID: %w", err)
		}

		if err := c.UI.DeleteLayout(ctx, layoutID); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Layout deleted: %s", layoutID))
		return nil
	},
}

var layoutQuotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Show quota usage",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		quota, err := c.UI.GetQuota(ctx)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, quota)
		}

		headers := []string{"METRIC", "USED", "MAX"}
		rows := [][]string{
			{"Layouts", fmt.Sprintf("%d", quota.Layouts.Used), fmt.Sprintf("%d", quota.Layouts.Max)},
			{"Archived", fmt.Sprintf("%d", quota.Layouts.Archived), "-"},
			{"Assets/Layout", "-", fmt.Sprintf("%d", quota.AssetsPerLayout.Max)},
		}

		fmt.Printf("Plan: %s (Status: %s)\n\n", quota.Plan, quota.Status)
		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}

// Layout archive subcommand
var layoutArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Manage archived layouts",
}

func init() {
	layoutArchiveCmd.AddCommand(layoutArchiveListCmd)
	layoutArchiveCmd.AddCommand(layoutRestoreCmd)
}

var layoutArchiveListCmd = &cobra.Command{
	Use:   "list",
	Short: "List archived layouts",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		layouts, err := c.UI.GetArchivedLayouts(ctx)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, layouts)
		}

		if len(layouts) == 0 {
			fmt.Println("No archived layouts found")
			return nil
		}

		headers := []string{"ID", "NAME"}
		rows := make([][]string, len(layouts))
		for i, l := range layouts {
			rows[i] = []string{l.LayoutID.String(), fmt.Sprintf("(archived: %s)", l.ArchiveReason)}
		}

		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}

var layoutRestoreCmd = &cobra.Command{
	Use:   "restore <layout_id>",
	Short: "Restore an archived layout",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		layoutID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid layout ID: %w", err)
		}

		if err := c.UI.RestoreLayout(ctx, layoutID); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Layout restored: %s", layoutID))
		return nil
	},
}

// Helper to get authenticated client
func getClient() (*whooktown.Client, error) {
	return client.New(client.Options{
		Token:  tokenFlag,
		Config: cfg,
	})
}
