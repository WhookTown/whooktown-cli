/**
 * Format a table for terminal output
 */
export declare function formatTable(headers: string[], rows: string[][]): string;
/**
 * Format status with color
 */
export declare function formatStatus(status: string): string;
/**
 * Format activity with color
 */
export declare function formatActivity(activity: string): string;
/**
 * Format boolean as enabled/disabled
 */
export declare function formatEnabled(enabled: boolean): string;
/**
 * Format speed
 */
export declare function formatSpeed(speed: string): string;
/**
 * Format density as percentage with visual bar
 */
export declare function formatDensity(density: number): string;
/**
 * Truncate string with ellipsis
 */
export declare function truncate(str: string, maxLength: number): string;
/**
 * Format JSON output
 */
export declare function formatJSON(data: unknown): string;
/**
 * Print success message
 */
export declare function success(message: string): void;
/**
 * Print error message
 */
export declare function error(message: string): void;
/**
 * Print info message
 */
export declare function info(message: string): void;
/**
 * Print warning message
 */
export declare function warn(message: string): void;
//# sourceMappingURL=output.d.ts.map