import { WhooktownClient, Environment } from '@whooktown/sdk';
import { getToken, getEnvironment } from '../config/index.js';
import chalk from 'chalk';

/**
 * Create an authenticated SDK client
 * Exits with error if not logged in
 */
export function createClient(): WhooktownClient {
  const token = getToken();
  if (!token) {
    console.error(chalk.red('Error: Not logged in. Run: wt login <token>'));
    process.exit(1);
  }

  const env = getEnvironment();

  // Build custom URLs from environment variables
  const customUrls: Partial<{
    auth: string;
    sensor: string;
    ui: string;
    workflow: string;
    backoffice: string;
    sse: string;
    subscription: string;
  }> = {};

  if (process.env.WHOOKTOWN_AUTH_URL) {
    customUrls.auth = process.env.WHOOKTOWN_AUTH_URL;
  }
  if (process.env.WHOOKTOWN_SENSOR_URL) {
    customUrls.sensor = process.env.WHOOKTOWN_SENSOR_URL;
  }

  return new WhooktownClient({
    token,
    environment: env === 'DEV' ? Environment.Development : Environment.Production,
    urls: Object.keys(customUrls).length > 0 ? customUrls : undefined,
  });
}

/**
 * Create an unauthenticated client for token validation
 */
export function createUnauthClient(): WhooktownClient {
  const env = getEnvironment();

  const customUrls: Partial<{
    auth: string;
    sensor: string;
  }> = {};

  if (process.env.WHOOKTOWN_AUTH_URL) {
    customUrls.auth = process.env.WHOOKTOWN_AUTH_URL;
  }
  if (process.env.WHOOKTOWN_SENSOR_URL) {
    customUrls.sensor = process.env.WHOOKTOWN_SENSOR_URL;
  }

  return new WhooktownClient({
    environment: env === 'DEV' ? Environment.Development : Environment.Production,
    urls: Object.keys(customUrls).length > 0 ? customUrls : undefined,
  });
}
