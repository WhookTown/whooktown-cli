package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
)

var audioCmd = &cobra.Command{
	Use:   "audio",
	Short: "Audio controls",
	Long:  "Get and set audio settings for layouts",
}

func init() {
	audioCmd.AddCommand(audioGetCmd)
	audioCmd.AddCommand(audioSetCmd)
}

var audioGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get audio state",
	Long: `Get audio state for all layouts.

Examples:
  wt audio get`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		states, err := c.Audio.GetStates(ctx)
		if err != nil {
			return err
		}

		if jsonOutput {
			return formatter.Format(os.Stdout, states)
		}

		if len(states) == 0 {
			fmt.Println("No audio states found")
			return nil
		}

		headers := []string{"LAYOUT", "ENABLED", "MOOD", "MUSIC VOL", "SFX VOL", "AUTO"}
		rows := make([][]string, 0, len(states))
		for _, state := range states {
			enabled := "no"
			if state.Enabled {
				enabled = "yes"
			}
			autoMood := "no"
			if state.AutoMood {
				autoMood = "yes"
			}
			layoutID := state.LayoutID
			if len(layoutID) > 12 {
				layoutID = layoutID[:8] + "..."
			}
			rows = append(rows, []string{
				layoutID,
				enabled,
				state.CurrentMood,
				fmt.Sprintf("%d%%", state.MusicVolume),
				fmt.Sprintf("%d%%", state.SfxVolume),
				autoMood,
			})
		}

		return formatter.FormatTable(os.Stdout, headers, rows)
	},
}

var (
	audioLayout      string
	audioMood        string
	audioMusicVolume int
	audioSfxVolume   int
	audioEnabled     bool
	audioAutoMood    bool
)

var audioSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set audio state",
	Long: `Set audio state for a layout.

Examples:
  wt audio set --layout <id> --mood calm
  wt audio set --layout <id> --music-volume 80 --sfx-volume 60
  wt audio set --layout <id> --enabled --auto-mood`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		audioCmd := &whooktown.AudioCommand{
			LayoutID: audioLayout,
		}

		// Determine command based on flags
		if cmd.Flags().Changed("mood") {
			audioCmd.Command = "mood"
			audioCmd.Mood = audioMood
		} else if cmd.Flags().Changed("music-volume") || cmd.Flags().Changed("sfx-volume") {
			audioCmd.Command = "volume"
			if cmd.Flags().Changed("music-volume") {
				audioCmd.MusicVolume = &audioMusicVolume
			}
			if cmd.Flags().Changed("sfx-volume") {
				audioCmd.SfxVolume = &audioSfxVolume
			}
		} else if cmd.Flags().Changed("enabled") || cmd.Flags().Changed("auto-mood") {
			audioCmd.Command = "toggle"
			if cmd.Flags().Changed("enabled") {
				audioCmd.Enabled = &audioEnabled
			}
			if cmd.Flags().Changed("auto-mood") {
				audioCmd.AutoMood = &audioAutoMood
			}
		} else {
			return fmt.Errorf("at least one setting flag is required (--mood, --music-volume, --sfx-volume, --enabled, --auto-mood)")
		}

		if err := c.Audio.SendCommand(ctx, audioCmd); err != nil {
			return err
		}

		formatter.Success("Audio settings updated")
		return nil
	},
}

func init() {
	audioSetCmd.Flags().StringVar(&audioLayout, "layout", "", "layout ID (required)")
	audioSetCmd.Flags().StringVar(&audioMood, "mood", "", "audio mood (calm, active, tension, critical, epic)")
	audioSetCmd.Flags().IntVar(&audioMusicVolume, "music-volume", 80, "music volume (0-100)")
	audioSetCmd.Flags().IntVar(&audioSfxVolume, "sfx-volume", 80, "SFX volume (0-100)")
	audioSetCmd.Flags().BoolVar(&audioEnabled, "enabled", true, "enable audio")
	audioSetCmd.Flags().BoolVar(&audioAutoMood, "auto-mood", true, "enable automatic mood")
	audioSetCmd.MarkFlagRequired("layout")
}
