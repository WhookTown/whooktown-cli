package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
)

var sensorCmd = &cobra.Command{
	Use:   "sensor",
	Short: "Send sensor data",
	Long:  "Send sensor data to update building states",
}

func init() {
	sensorCmd.AddCommand(sensorSendCmd)
}

var (
	sensorID       string
	sensorStatus   string
	sensorActivity string
	sensorFile     string
	sensorExtra    []string
)

var sensorSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send sensor data to a building",
	Long: `Send sensor data to update a building's state.

Examples:
  # Set building status
  wt sensor send --id <building_id> --status online --activity fast

  # Send from JSON file
  wt sensor send -f sensor-data.json

  # Send with extra fields (for DataCenter, Arcade, etc.)
  wt sensor send --id <id> --status online --extra cpuUsage=75 --extra temperature=42`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		c, err := getClient()
		if err != nil {
			return err
		}

		var data *whooktown.SensorData

		if sensorFile != "" {
			// Read from file
			fileData, err := os.ReadFile(sensorFile)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			data = &whooktown.SensorData{}
			if err := json.Unmarshal(fileData, data); err != nil {
				return fmt.Errorf("invalid JSON: %w", err)
			}
		} else {
			// Build from flags
			if sensorID == "" {
				return fmt.Errorf("--id is required")
			}

			id, err := uuid.FromString(sensorID)
			if err != nil {
				return fmt.Errorf("invalid building ID: %w", err)
			}

			data = &whooktown.SensorData{
				ID: id,
			}

			if sensorStatus != "" {
				data.Status = whooktown.Status(sensorStatus)
			}
			if sensorActivity != "" {
				data.Activity = whooktown.Activity(sensorActivity)
			}

			// Handle extra fields
			if len(sensorExtra) > 0 {
				data.Extra = make(map[string]interface{})
				for _, kv := range sensorExtra {
					parts := strings.SplitN(kv, "=", 2)
					if len(parts) != 2 {
						return fmt.Errorf("invalid extra format: %s (expected key=value)", kv)
					}
					key, value := parts[0], parts[1]

					// Try to parse as number
					if i, err := strconv.Atoi(value); err == nil {
						data.Extra[key] = i
					} else if f, err := strconv.ParseFloat(value, 64); err == nil {
						data.Extra[key] = f
					} else if b, err := strconv.ParseBool(value); err == nil {
						data.Extra[key] = b
					} else {
						data.Extra[key] = value
					}
				}
			}
		}

		if err := c.Sensors.Send(ctx, data); err != nil {
			return err
		}

		formatter.Success(fmt.Sprintf("Sensor data sent to %s", data.ID))
		return nil
	},
}

func init() {
	sensorSendCmd.Flags().StringVar(&sensorID, "id", "", "building ID")
	sensorSendCmd.Flags().StringVar(&sensorStatus, "status", "", "status (online, offline, warning, critical)")
	sensorSendCmd.Flags().StringVar(&sensorActivity, "activity", "", "activity level (slow, normal, fast)")
	sensorSendCmd.Flags().StringVarP(&sensorFile, "file", "f", "", "JSON file with sensor data")
	sensorSendCmd.Flags().StringArrayVar(&sensorExtra, "extra", nil, "extra fields (key=value)")
}
