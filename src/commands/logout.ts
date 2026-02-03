import { Command } from 'commander';
import { clearToken, isLoggedIn } from '../config/index.js';
import { success, info } from '../lib/output.js';

export const logoutCommand = new Command('logout')
  .description('Clear saved token')
  .action(() => {
    if (!isLoggedIn()) {
      info('Not logged in');
      return;
    }

    clearToken();
    success('Logged out successfully');
  });
