import { Command } from 'commander';
import chalk from 'chalk';
import { setToken, isLoggedIn, getConfigPath } from '../config/index.js';
import { createUnauthClient } from '../lib/client.js';
import { success, error, info, warn } from '../lib/output.js';

export const loginCommand = new Command('login')
  .description('Login with a sensor token')
  .argument('<token>', 'Sensor token to authenticate with')
  .option('--no-validate', 'Skip token validation')
  .action(async (token: string, options: { validate: boolean }) => {
    // Check if already logged in
    if (isLoggedIn()) {
      warn('Already logged in. Use "wt logout" first to switch accounts.');
    }

    if (options.validate) {
      info('Validating token...');

      try {
        const client = createUnauthClient();
        // Set the token temporarily for validation
        client.setToken(token);
        const tokenInfo = await client.auth.checkToken(token);

        // Validate token type
        if (tokenInfo.type !== 'sensor') {
          error(`Invalid token type: ${tokenInfo.type}`);
          console.error(chalk.gray('This CLI only accepts sensor tokens.'));
          console.error(chalk.gray('Expected type: sensor'));
          process.exit(1);
        }

        // Validate roles
        const roles = tokenInfo.roles || {};
        if (roles.sensor !== 'rw') {
          error(`Invalid token roles: ${JSON.stringify(roles)}`);
          console.error(chalk.gray('Expected roles: {"sensor": "rw"}'));
          process.exit(1);
        }

        // Save token
        const accountId = tokenInfo.account_id || '';
        setToken(token, accountId);

        success('Logged in successfully!');
        console.log(chalk.gray(`  Account ID: ${accountId}`));
        console.log(chalk.gray(`  Token type: ${tokenInfo.type}`));
        console.log(chalk.gray(`  Config: ${getConfigPath()}`));
      } catch (err) {
        error('Token validation failed');
        if (err instanceof Error) {
          console.error(chalk.gray(err.message));
        }
        process.exit(1);
      }
    } else {
      // Skip validation, just save
      warn('Skipping token validation');
      setToken(token, '');
      success('Token saved');
      console.log(chalk.gray(`  Config: ${getConfigPath()}`));
    }
  });
