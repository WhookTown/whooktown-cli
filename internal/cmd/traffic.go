package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
)

var trafficCmd = &cobra.Command{
	Use:   "traffic",
	Short: "Traffic controls",
	Long:  "Get and set traffic settings for layouts",
}

func init() {
	trafficCmd.AddCommand(trafficGetCmd)
	trafficCmd.AddCommand(trafficSetCmd)
}

var trafficGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get traffic state",
	Long: `Get traffic state for all layouts.

Examples:
  wt traffic get`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		states, err := c.Traffic.GetStates(ctx)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, states)
		}

		if len(states) == 0 {
			fmt.Println("No traffic states found")
			return nil
		}

		headers := []string{"LAYOUT", "DENSITY", "SPEED", "ENABLED", "LABELS"}
		rows := make([][]string, 0, len(states))
		for _, state := range states {
			enabled := "no"
			if state.Enabled {
				enabled = "yes"
			}
			labels := "hidden"
			if state.LabelsVisible {
				labels = "visible"
			}
			layoutID := state.LayoutID
			if len(layoutID) > 12 {
				layoutID = layoutID[:8] + "..."
			}
			rows = append(rows, []string{
				layoutID,
				fmt.Sprintf("%d%%", state.Density),
				state.Speed,
				enabled,
				labels,
			})
		}

		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}

var (
	trafficLayout  string
	trafficDensity int
	trafficSpeed   string
	trafficEnabled bool
)

var trafficSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set traffic state",
	Long: `Set traffic state for a layout.

Examples:
  wt traffic set --layout <id> --density 50
  wt traffic set --layout <id> --speed fast --enabled
  wt traffic set --layout <id> --density 0 --enabled=false`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		trafficCmd := &whooktown.TrafficCommand{
			LayoutID: trafficLayout,
		}

		if cmd.Flags().Changed("density") {
			trafficCmd.Density = trafficDensity
		}
		if cmd.Flags().Changed("speed") {
			trafficCmd.Speed = trafficSpeed
		}
		if cmd.Flags().Changed("enabled") {
			trafficCmd.Enabled = &trafficEnabled
		}

		if err := c.Traffic.SendCommand(ctx, trafficCmd); err != nil {
			return err
		}

		formatter.Success("Traffic settings updated")
		return nil
	},
}

func init() {
	trafficSetCmd.Flags().StringVar(&trafficLayout, "layout", "", "layout ID (required)")
	trafficSetCmd.Flags().IntVar(&trafficDensity, "density", 50, "traffic density (0-100)")
	trafficSetCmd.Flags().StringVar(&trafficSpeed, "speed", "normal", "traffic speed (slow, normal, fast)")
	trafficSetCmd.Flags().BoolVar(&trafficEnabled, "enabled", true, "enable traffic")
	trafficSetCmd.MarkFlagRequired("layout")
}
