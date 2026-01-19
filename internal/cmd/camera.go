package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
)

var cameraCmd = &cobra.Command{
	Use:   "camera",
	Short: "Camera controls",
	Long:  "Manage camera presets, paths, and send camera commands",
}

func init() {
	cameraCmd.AddCommand(cameraPresetCmd)
	cameraCmd.AddCommand(cameraPathCmd)
	cameraCmd.AddCommand(cameraCommandCmd)
}

// Camera preset commands
var cameraPresetCmd = &cobra.Command{
	Use:   "preset",
	Short: "Manage camera presets",
}

func init() {
	cameraPresetCmd.AddCommand(cameraPresetListCmd)
	cameraPresetCmd.AddCommand(cameraPresetCreateCmd)
	cameraPresetCmd.AddCommand(cameraPresetDeleteCmd)
}

var cameraPresetListCmd = &cobra.Command{
	Use:   "list <layout_id>",
	Short: "List camera presets for a layout",
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

		presets, err := c.Camera.ListPresets(ctx, layoutID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, presets)
		}

		if len(presets) == 0 {
			fmt.Println("No camera presets found")
			return nil
		}

		headers := []string{"ID", "NAME", "DEFAULT", "POSITION"}
		rows := make([][]string, len(presets))
		for i, p := range presets {
			isDefault := ""
			if p.IsDefault {
				isDefault = "*"
			}
			pos := fmt.Sprintf("%.1f, %.1f, %.1f", p.PositionX, p.PositionY, p.PositionZ)
			rows[i] = []string{p.ID.String()[:8] + "...", p.Name, isDefault, pos}
		}

		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}

var (
	presetName     string
	presetLayout   string
	presetPosition string
	presetRotation string
)

var cameraPresetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a camera preset",
	Long: `Create a camera preset for a layout.

Example:
  wt camera preset create --layout <id> --name "Overview" --position "10,5,10"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		layoutID, err := uuid.FromString(presetLayout)
		if err != nil {
			return fmt.Errorf("invalid layout ID: %w", err)
		}

		// Parse position
		pos, err := parseVector3(presetPosition)
		if err != nil {
			return fmt.Errorf("invalid position: %w", err)
		}

		req := &whooktown.CreatePresetRequest{
			LayoutID:  layoutID,
			Name:      presetName,
			PositionX: pos.X,
			PositionY: pos.Y,
			PositionZ: pos.Z,
		}

		// Parse rotation (optional)
		if presetRotation != "" {
			rot, err := parseVector3(presetRotation)
			if err != nil {
				return fmt.Errorf("invalid rotation: %w", err)
			}
			req.RotationX = rot.X
			req.RotationY = rot.Y
			req.RotationZ = rot.Z
		}

		preset, err := c.Camera.CreatePreset(ctx, req)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, preset)
		}

		formatter.Success(fmt.Sprintf("Preset created: %s (%s)", preset.Name, preset.ID))
		return nil
	},
}

func init() {
	cameraPresetCreateCmd.Flags().StringVar(&presetLayout, "layout", "", "layout ID (required)")
	cameraPresetCreateCmd.Flags().StringVar(&presetName, "name", "", "preset name (required)")
	cameraPresetCreateCmd.Flags().StringVar(&presetPosition, "position", "", "position as x,y,z (required)")
	cameraPresetCreateCmd.Flags().StringVar(&presetRotation, "rotation", "", "rotation as x,y,z (optional)")
	cameraPresetCreateCmd.MarkFlagRequired("layout")
	cameraPresetCreateCmd.MarkFlagRequired("name")
	cameraPresetCreateCmd.MarkFlagRequired("position")
}

var cameraPresetDeleteCmd = &cobra.Command{
	Use:   "delete <preset_id>",
	Short: "Delete a camera preset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		presetID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid preset ID: %w", err)
		}

		if err := c.Camera.DeletePreset(ctx, presetID); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Preset deleted: %s", args[0]))
		return nil
	},
}

// Camera path commands
var cameraPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Manage camera paths",
}

func init() {
	cameraPathCmd.AddCommand(cameraPathListCmd)
	cameraPathCmd.AddCommand(cameraPathCreateCmd)
	cameraPathCmd.AddCommand(cameraPathDeleteCmd)
}

var cameraPathListCmd = &cobra.Command{
	Use:   "list <layout_id>",
	Short: "List camera paths for a layout",
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

		paths, err := c.Camera.ListPaths(ctx, layoutID)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, paths)
		}

		if len(paths) == 0 {
			fmt.Println("No camera paths found")
			return nil
		}

		headers := []string{"ID", "NAME", "CHECKPOINTS", "LOOP"}
		rows := make([][]string, len(paths))
		for i, p := range paths {
			loop := "no"
			if p.Loop {
				loop = "yes"
			}
			rows[i] = []string{
				p.ID.String()[:8] + "...",
				p.Name,
				fmt.Sprintf("%d", len(p.Checkpoints)),
				loop,
			}
		}

		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}

var (
	pathName   string
	pathLayout string
	pathLoop   bool
)

var cameraPathCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a camera path",
	Long: `Create a camera path for a layout.

Example:
  wt camera path create --layout <id> --name "Tour" --loop`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		layoutID, err := uuid.FromString(pathLayout)
		if err != nil {
			return fmt.Errorf("invalid layout ID: %w", err)
		}

		path, err := c.Camera.CreatePath(ctx, &whooktown.CreatePathRequest{
			LayoutID: layoutID,
			Name:     pathName,
			Loop:     pathLoop,
		})
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, path)
		}

		formatter.Success(fmt.Sprintf("Path created: %s (%s)", path.Name, path.ID))
		return nil
	},
}

func init() {
	cameraPathCreateCmd.Flags().StringVar(&pathLayout, "layout", "", "layout ID (required)")
	cameraPathCreateCmd.Flags().StringVar(&pathName, "name", "", "path name (required)")
	cameraPathCreateCmd.Flags().BoolVar(&pathLoop, "loop", false, "loop the path")
	cameraPathCreateCmd.MarkFlagRequired("layout")
	cameraPathCreateCmd.MarkFlagRequired("name")
}

var cameraPathDeleteCmd = &cobra.Command{
	Use:   "delete <path_id>",
	Short: "Delete a camera path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		pathID, err := uuid.FromString(args[0])
		if err != nil {
			return fmt.Errorf("invalid path ID: %w", err)
		}

		if err := c.Camera.DeletePath(ctx, pathID); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Path deleted: %s", args[0]))
		return nil
	},
}

// Camera command
var (
	cmdLayout   string
	cmdMode     string
	cmdPreset   string
	cmdPath     string
	cmdPosition string
	cmdAction   string
)

var cameraCommandCmd = &cobra.Command{
	Use:   "command",
	Short: "Send a camera command",
	Long: `Send a camera command to a layout.

Examples:
  wt camera command --layout <id> --mode orbit
  wt camera command --layout <id> --preset <preset_id>
  wt camera command --layout <id> --path <path_id> --action play
  wt camera command --layout <id> --position "10,5,10"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		camCmd := &whooktown.CameraCommand{
			LayoutID: cmdLayout,
		}

		if cmdMode != "" {
			camCmd.Command = "mode"
			camCmd.Mode = cmdMode
		}
		if cmdPreset != "" {
			camCmd.Command = "preset"
			camCmd.PresetID = cmdPreset
		}
		if cmdPath != "" {
			camCmd.Command = "path"
			camCmd.PathID = cmdPath
			if cmdAction != "" {
				camCmd.Action = cmdAction
			} else {
				camCmd.Action = "play"
			}
		}
		if cmdPosition != "" {
			pos, err := parseVector3(cmdPosition)
			if err != nil {
				return fmt.Errorf("invalid position: %w", err)
			}
			camCmd.Command = "position"
			camCmd.Position = pos
		}

		if err := c.Camera.SendCommand(ctx, camCmd); err != nil {
			return err
		}

		formatter.Success("Camera command sent")
		return nil
	},
}

func init() {
	cameraCommandCmd.Flags().StringVar(&cmdLayout, "layout", "", "layout ID (required)")
	cameraCommandCmd.Flags().StringVar(&cmdMode, "mode", "", "camera mode (orbit, fps, flyover)")
	cameraCommandCmd.Flags().StringVar(&cmdPreset, "preset", "", "preset ID")
	cameraCommandCmd.Flags().StringVar(&cmdPath, "path", "", "path ID")
	cameraCommandCmd.Flags().StringVar(&cmdAction, "action", "", "path action (play, pause, stop)")
	cameraCommandCmd.Flags().StringVar(&cmdPosition, "position", "", "position as x,y,z")
	cameraCommandCmd.MarkFlagRequired("layout")
}

// parseVector3 parses a "x,y,z" string into a Vector3
func parseVector3(s string) (*whooktown.Vector3, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("expected x,y,z format")
	}

	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid x: %w", err)
	}
	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid y: %w", err)
	}
	z, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid z: %w", err)
	}

	return &whooktown.Vector3{X: x, Y: y, Z: z}, nil
}
