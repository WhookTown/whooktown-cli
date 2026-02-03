import { WhooktownClient, Environment } from '@whooktown/sdk';
import { getToken, getEnvironment } from '../config/index.js';
import chalk from 'chalk';
/**
 * Create an authenticated SDK client
 * Exits with error if not logged in
 */
export function createClient() {
    const token = getToken();
    if (!token) {
        console.error(chalk.red('Error: Not logged in. Run: wt login <token>'));
        process.exit(1);
    }
    const env = getEnvironment();
    // Build custom URLs from environment variables
    const customUrls = {};
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
export function createUnauthClient() {
    const env = getEnvironment();
    const customUrls = {};
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
//# sourceMappingURL=client.js.map