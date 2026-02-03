import { Command } from 'commander';
import { createClient } from '../lib/client.js';
import { handleError } from '../lib/errors.js';
import {
  formatTable,
  formatStatus,
  formatActivity,
  formatJSON,
  success,
  error,
  warn,
} from '../lib/output.js';

interface LayoutData {
  name?: string;
  buildings?: Array<{ id: string; name?: string }>;
}

interface SensorInfo {
  layoutName: string;
  buildingName: string;
}

const sendCommand = new Command('send')
  .description('Send sensor data')
  .argument('<id>', 'Sensor ID (UUID)')
  .option('-n, --name <name>', 'Sensor name (for display)')
  .option('-s, --status <status>', 'Status: online, offline, warning, critical')
  .option('-a, --activity <activity>', 'Activity: slow, normal, fast')
  .option('-j, --json <json>', 'Additional JSON fields (e.g., \'{"cpuUsage": 75}\')')
  .option('-q, --quiet', 'Suppress output')
  .action(async (id: string, options: {
    name?: string;
    status?: string;
    activity?: string;
    json?: string;
    quiet?: boolean;
  }) => {
    try {
      const client = createClient();

      // Build sensor data
      const data: Record<string, unknown> = { id };

      if (options.name) {
        data.name = options.name;
      }

      if (options.status) {
        const validStatuses = ['online', 'offline', 'warning', 'critical'];
        if (!validStatuses.includes(options.status.toLowerCase())) {
          error(`Invalid status: ${options.status}`);
          console.error(`Valid values: ${validStatuses.join(', ')}`);
          process.exit(1);
        }
        data.status = options.status.toLowerCase();
      }

      if (options.activity) {
        const validActivities = ['slow', 'normal', 'fast'];
        if (!validActivities.includes(options.activity.toLowerCase())) {
          error(`Invalid activity: ${options.activity}`);
          console.error(`Valid values: ${validActivities.join(', ')}`);
          process.exit(1);
        }
        data.activity = options.activity.toLowerCase();
      }

      // Parse additional JSON
      if (options.json) {
        try {
          const extra = JSON.parse(options.json);
          Object.assign(data, extra);
        } catch {
          error('Invalid JSON in --json option');
          process.exit(1);
        }
      }

      await client.sensors.sendRaw(data);

      if (!options.quiet) {
        success(`Sent data for sensor: ${id}`);
      }
    } catch (err) {
      handleError(err);
    }
  });

const listCommand = new Command('list')
  .description('List sensor states')
  .option('-f, --format <format>', 'Output format: table, json', 'table')
  .action(async (options: { format: string }) => {
    try {
      const client = createClient();
      const sensors = await client.sensors.getSensors();

      // Try to fetch layouts to map sensor IDs to layout/building names
      let sensorInfoMap = new Map<string, SensorInfo>();
      try {
        const layouts = await client.ui.getLayouts();
        for (const layout of layouts) {
          const layoutData = layout.data as LayoutData;
          const layoutName = layoutData?.name || layout.layout_id;
          if (layoutData?.buildings) {
            for (const building of layoutData.buildings) {
              sensorInfoMap.set(building.id, {
                layoutName,
                buildingName: building.name || '-',
              });
            }
          }
        }
      } catch {
        // Silently ignore if we can't fetch layouts (e.g., token doesn't have ui:r)
        warn('Could not fetch layouts (token may lack ui:r permission)');
      }

      if (options.format === 'json') {
        // Add layout info to JSON output
        const enriched = sensors.map((s) => {
          const info = sensorInfoMap.get(s.id);
          return {
            ...s,
            buildingName: info?.buildingName || null,
            layoutName: info?.layoutName || null,
          };
        });
        console.log(formatJSON(enriched));
        return;
      }

      // Table format
      const headers = ['ID', 'Name', 'Layout', 'Status', 'Activity', 'Updated'];
      const rows = sensors.map((s) => {
        const info = sensorInfoMap.get(s.id);
        return [
          s.id || '-',
          info?.buildingName || '-',
          info?.layoutName || '-',
          formatStatus(s.data?.status as string || ''),
          formatActivity(s.data?.activity as string || ''),
          s.received_at ? new Date(s.received_at).toLocaleString() : '-',
        ];
      });

      console.log(formatTable(headers, rows));
      console.log(`\n${sensors.length} sensor(s)`);
    } catch (err) {
      handleError(err);
    }
  });

export const sensorCommand = new Command('sensor')
  .description('Manage sensors')
  .addCommand(sendCommand)
  .addCommand(listCommand);
