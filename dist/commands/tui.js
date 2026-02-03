import { Command } from 'commander';
import React from 'react';
import { render } from 'ink';
import { App } from '../tui/App.js';
import { isLoggedIn } from '../config/index.js';
import { error } from '../lib/output.js';
export const tuiCommand = new Command('tui')
    .description('Launch interactive TUI dashboard')
    .option('-r, --refresh <ms>', 'Auto-refresh interval in milliseconds', '5000')
    .action((options) => {
    // Check login
    if (!isLoggedIn()) {
        error('Not logged in. Run: wt login <token>');
        process.exit(1);
    }
    const refreshInterval = parseInt(options.refresh, 10);
    if (isNaN(refreshInterval) || refreshInterval < 1000) {
        error('Refresh interval must be at least 1000ms');
        process.exit(1);
    }
    // Render Ink app
    render(React.createElement(App, { refreshInterval }));
});
//# sourceMappingURL=tui.js.map