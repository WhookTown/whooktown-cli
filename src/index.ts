#!/usr/bin/env node
import { program } from 'commander';
import { loginCommand } from './commands/login.js';
import { logoutCommand } from './commands/logout.js';
import { sensorCommand } from './commands/sensor.js';
import { trafficCommand } from './commands/traffic.js';
import { layoutCommand } from './commands/layout.js';
import { tuiCommand } from './commands/tui.js';
import { popupCommand } from './commands/popup.js';

program
  .name('wt')
  .description('Whooktown CLI - Control your 3D IT city')
  .version('1.0.0');

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
