import chalk from 'chalk';
import { isUnauthorized, isNetworkError, isTimeout, isBadRequest, isNotFound } from '@whooktown/sdk';
/**
 * Handle SDK errors and exit with appropriate message
 */
export function handleError(err) {
    if (isUnauthorized(err)) {
        console.error(chalk.red('Authentication failed. Your token may be invalid or expired.'));
        console.error(chalk.gray('Run: wt login <token>'));
        process.exit(1);
    }
    if (isNetworkError(err)) {
        console.error(chalk.red('Network error. Check your internet connection.'));
        process.exit(1);
    }
    if (isTimeout(err)) {
        console.error(chalk.red('Request timed out. The server may be slow or unreachable.'));
        process.exit(1);
    }
    if (isBadRequest(err)) {
        const message = err instanceof Error ? err.message : 'Invalid request';
        console.error(chalk.red('Bad request:'), message);
        process.exit(1);
    }
    if (isNotFound(err)) {
        console.error(chalk.red('Not found. The resource does not exist.'));
        process.exit(1);
    }
    // Generic error
    const message = err instanceof Error ? err.message : String(err);
    console.error(chalk.red('Error:'), message);
    process.exit(1);
}
//# sourceMappingURL=errors.js.map