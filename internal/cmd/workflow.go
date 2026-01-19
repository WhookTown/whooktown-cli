package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
)

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage workflows",
	Long:  "Create, list, delete, enable, and disable workflows",
}

func init() {
	workflowCmd.AddCommand(workflowListCmd)
	workflowCmd.AddCommand(workflowShowCmd)
	workflowCmd.AddCommand(workflowCreateCmd)
	workflowCmd.AddCommand(workflowDeleteCmd)
	workflowCmd.AddCommand(workflowEnableCmd)
	workflowCmd.AddCommand(workflowDisableCmd)
	workflowCmd.AddCommand(workflowExportCmd)
	workflowCmd.AddCommand(workflowOperationsCmd)
}

var workflowListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workflows",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		workflows, err := c.Workflow.List(ctx)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, workflows)
		}

		if len(workflows) == 0 {
			fmt.Println("No workflows found")
			return nil
		}

		headers := []string{"ID", "NAME", "ENABLED", "CREATED"}
		rows := make([][]string, len(workflows))
		for i, w := range workflows {
			enabled := "no"
			if w.Enabled {
				enabled = "yes"
			}
			rows[i] = []string{
				w.ID.String()[:8] + "...",
				w.Name,
				enabled,
				w.CreatedAt.Format("2006-01-02"),
			}
		}

		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}

var workflowShowCmd = &cobra.Command{
	Use:   "show <workflow_id>",
	Short: "Show workflow details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		wfID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid workflow ID: %w", err)
		}

		workflows, err := c.Workflow.List(ctx)
		if err != nil {
			return err
		}

		// Find the workflow
		for _, w := range workflows {
			if w.ID == wfID {
				return formatter.Format(os.Stdout, w)
			}
		}

		return fmt.Errorf("workflow not found: %s", args[0])
	},
}

var workflowFile string

var workflowCreateCmd = &cobra.Command{
	Use:   "create -f <workflow.json>",
	Short: "Create a workflow from JSON file",
	Long: `Create a workflow from a JSON file.

The JSON file should contain:
{
  "name": "My Workflow",
  "graph": {
    "node1": { "operator": "input", "name": "sensor-id" },
    "node2": { "operator": "output", "name": "building-id", "inputs": ["node1"] }
  },
  "enabled": true
}`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		if workflowFile == "" {
			return fmt.Errorf("--file is required")
		}

		data, err := os.ReadFile(workflowFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		var req whooktown.CreateWorkflowRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		workflow, err := c.Workflow.Create(ctx, &req)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, workflow)
		}

		formatter.Success(fmt.Sprintf("Workflow created: %s (%s)", workflow.Name, workflow.ID))
		return nil
	},
}

func init() {
	workflowCreateCmd.Flags().StringVarP(&workflowFile, "file", "f", "", "workflow JSON file (required)")
	workflowCreateCmd.MarkFlagRequired("file")
}

var workflowDeleteCmd = &cobra.Command{
	Use:   "delete <workflow_id>",
	Short: "Delete a workflow",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		wfID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid workflow ID: %w", err)
		}

		if err := c.Workflow.Delete(ctx, wfID); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Workflow deleted: %s", args[0]))
		return nil
	},
}

var workflowEnableCmd = &cobra.Command{
	Use:   "enable <workflow_id>",
	Short: "Enable a workflow",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		wfID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid workflow ID: %w", err)
		}

		if err := c.Workflow.Enable(ctx, wfID); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Workflow enabled: %s", args[0]))
		return nil
	},
}

var workflowDisableCmd = &cobra.Command{
	Use:   "disable <workflow_id>",
	Short: "Disable a workflow",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		wfID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid workflow ID: %w", err)
		}

		if err := c.Workflow.Disable(ctx, wfID); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Workflow disabled: %s", args[0]))
		return nil
	},
}

var workflowExportCmd = &cobra.Command{
	Use:   "export <workflow_id>",
	Short: "Export a workflow as JSON",
	Long: `Export a workflow as JSON to stdout.

Example:
  wt workflow export <id> > workflow.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		wfID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid workflow ID: %w", err)
		}

		workflows, err := c.Workflow.List(ctx)
		if err != nil {
			return err
		}

		// Find the workflow
		for _, w := range workflows {
			if w.ID == wfID {
				// Export as CreateWorkflowRequest format for re-import
				export := map[string]interface{}{
					"name":    w.Name,
					"graph":   json.RawMessage(w.Graph),
					"enabled": w.Enabled,
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(export)
			}
		}

		return fmt.Errorf("workflow not found: %s", args[0])
	},
}

var workflowOperationsCmd = &cobra.Command{
	Use:   "operations",
	Short: "List available workflow operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		ops, err := c.Workflow.GetOperations(ctx)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, ops)
		}

		headers := []string{"OPERATION", "DESCRIPTION", "OUTPUT TYPE"}
		rows := make([][]string, 0, len(ops))
		for name, op := range ops {
			desc := op.Description
			if len(desc) > 40 {
				desc = desc[:37] + "..."
			}
			rows = append(rows, []string{name, desc, op.OutputType})
		}

		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}
