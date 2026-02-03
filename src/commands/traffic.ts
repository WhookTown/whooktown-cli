import { Command } from 'commander';
import { createClient } from '../lib/client.js';
import { handleError } from '../lib/errors.js';
import {
  formatTable,
  formatJSON,
  formatEnabled,
  formatSpeed,
  formatDensity,
  truncate,
  success,
  error,
} from '../lib/output.js';
import chalk from 'chalk';

const setCommand = new Command('set')
  .description('Set traffic state for a layout')
  .argument('<layoutId>', 'Layout ID')
  .option('-d, --density <density>', 'Traffic density (0-100)')
  .option('-s, --speed <speed>', 'Traffic speed: slow, normal, fast')
  .option('--enabled', 'Enable traffic')
  .option('--disabled', 'Disable traffic')
  .action(async (layoutId: string, options: {
    density?: string;
    speed?: string;
    enabled?: boolean;
    disabled?: boolean;
  }) => {
    try {
      // Validate at least one option is provided
      if (!options.density && !options.speed && !options.enabled && !options.disabled) {
        error('At least one option is required: --density, --speed, --enabled, or --disabled');
        process.exit(1);
      }

      // Get current state first
      const client = createClient();
      const states = await client.sensors.getTrafficStates();
      const current = states.find((s) => s.layout_id === layoutId);

      // Default values from current state or sensible defaults
      let density = current?.density ?? 50;
      let speed = current?.speed ?? 'normal';
      let enabled = current?.enabled ?? true;

      // Apply options
      if (options.density !== undefined) {
        density = parseInt(options.density, 10);
        if (isNaN(density) || density < 0 || density > 100) {
          error('Density must be between 0 and 100');
          process.exit(1);
        }
      }

      if (options.speed) {
        const validSpeeds = ['slow', 'normal', 'fast'];
        if (!validSpeeds.includes(options.speed.toLowerCase())) {
          error(`Invalid speed: ${options.speed}`);
          console.error(`Valid values: ${validSpeeds.join(', ')}`);
          process.exit(1);
        }
        speed = options.speed.toLowerCase();
      }

      if (options.enabled) {
        enabled = true;
      } else if (options.disabled) {
        enabled = false;
      }

      await client.sensors.setTrafficState(layoutId, density, speed, enabled);

      success('Traffic updated');
      console.log(chalk.gray(`  Density: ${density}%`));
      console.log(chalk.gray(`  Speed: ${speed}`));
      console.log(chalk.gray(`  Enabled: ${enabled}`));
    } catch (err) {
      handleError(err);
    }
  });

const listCommand = new Command('list')
  .description('List traffic states')
  .option('-f, --format <format>', 'Output format: table, json', 'table')
  .action(async (options: { format: string }) => {
    try {
      const client = createClient();
      const states = await client.sensors.getTrafficStates();

      if (options.format === 'json') {
        console.log(formatJSON(states));
        return;
      }

      // Table format
      const headers = ['Layout ID', 'Density', 'Speed', 'Enabled'];
      const rows = states.map((s) => [
        truncate(s.layout_id || '', 20),
        formatDensity(s.density ?? 0),
        formatSpeed(s.speed || ''),
        formatEnabled(s.enabled ?? false),
      ]);

      console.log(formatTable(headers, rows));
      console.log(`\n${states.length} layout(s)`);
    } catch (err) {
      handleError(err);
    }
  });

export const trafficCommand = new Command('traffic')
  .description('Control traffic')
  .addCommand(setCommand)
  .addCommand(listCommand);
