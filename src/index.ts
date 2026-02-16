#!/usr/bin/env node
import { program } from 'commander';
import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import { loginCommand } from './commands/login.js';
import { logoutCommand } from './commands/logout.js';
import { sensorCommand } from './commands/sensor.js';
import { trafficCommand } from './commands/traffic.js';
import { layoutCommand } from './commands/layout.js';
import { tuiCommand } from './commands/tui.js';
import { popupCommand } from './commands/popup.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const pkg = JSON.parse(readFileSync(join(__dirname, '..', 'package.json'), 'utf-8'));

program
  .name('wt')
  .description('Whooktown CLI - Control your 3D IT city')
  .version(pkg.version);

// Add commands
program.addCommand(loginCommand);
program.addCommand(logoutCommand);
program.addCommand(sensorCommand);
program.addCommand(trafficCommand);
program.addCommand(layoutCommand);
program.addCommand(tuiCommand);
program.addCommand(popupCommand);

// Default: show help if no args
if (process.argv.length === 2) {
  program.outputHelp();
} else {
  program.parse();
}
