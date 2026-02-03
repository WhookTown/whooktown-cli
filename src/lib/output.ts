import chalk from 'chalk';

/**
 * Format a table for terminal output
 */
export function formatTable(headers: string[], rows: string[][]): string {
  if (rows.length === 0) {
    return chalk.gray('No data');
  }

  // Calculate column widths
  const widths = headers.map((h, i) =>
    Math.max(h.length, ...rows.map((r) => (r[i] || '').length))
  );

  // Format header
  const headerLine = headers.map((h, i) => h.padEnd(widths[i]!)).join('  ');
  const separator = widths.map((w) => '-'.repeat(w)).join('  ');

  // Format rows
  const rowLines = rows.map((row) =>
    row.map((cell, i) => (cell || '').padEnd(widths[i]!)).join('  ')
  );

  return [chalk.bold(headerLine), separator, ...rowLines].join('\n');
}

/**
 * Format status with color
 */
export function formatStatus(status: string): string {
  switch (status?.toLowerCase()) {
    case 'online':
      return chalk.green('online');
    case 'offline':
      return chalk.gray('offline');
    case 'warning':
      return chalk.yellow('warning');
    case 'critical':
      return chalk.red('critical');
    default:
      return status || chalk.gray('-');
  }
}

/**
 * Format activity with color
 */
export function formatActivity(activity: string): string {
  switch (activity?.toLowerCase()) {
    case 'slow':
      return chalk.blue('slow');
    case 'normal':
      return chalk.white('normal');
    case 'fast':
      return chalk.cyan('fast');
    default:
      return activity || chalk.gray('-');
  }
}

/**
 * Format boolean as enabled/disabled
 */
export function formatEnabled(enabled: boolean): string {
  return enabled ? chalk.green('enabled') : chalk.gray('disabled');
}

/**
 * Format speed
 */
export function formatSpeed(speed: string): string {
  switch (speed?.toLowerCase()) {
    case 'slow':
      return chalk.blue('slow');
    case 'normal':
      return chalk.white('normal');
    case 'fast':
      return chalk.cyan('fast');
    default:
      return speed || chalk.gray('-');
  }
}

/**
 * Format density as percentage with visual bar
 */
export function formatDensity(density: number): string {
  const percent = Math.round(density);
  const barLength = 10;
  const filled = Math.round((percent / 100) * barLength);
  const empty = barLength - filled;
  const bar = chalk.green('█'.repeat(filled)) + chalk.gray('░'.repeat(empty));
  return `${bar} ${percent}%`;
}

/**
 * Truncate string with ellipsis
 */
export function truncate(str: string, maxLength: number): string {
  if (!str) return '';
  if (str.length <= maxLength) return str;
  return str.slice(0, maxLength - 3) + '...';
}

/**
 * Format JSON output
 */
export function formatJSON(data: unknown): string {
  return JSON.stringify(data, null, 2);
}

/**
 * Print success message
 */
export function success(message: string): void {
  console.log(chalk.green('✓') + ' ' + message);
}

/**
 * Print error message
 */
export function error(message: string): void {
  console.error(chalk.red('✗') + ' ' + message);
}

/**
 * Print info message
 */
export function info(message: string): void {
  console.log(chalk.blue('ℹ') + ' ' + message);
}

/**
 * Print warning message
 */
export function warn(message: string): void {
  console.log(chalk.yellow('⚠') + ' ' + message);
}
