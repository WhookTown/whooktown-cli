/**
 * CLI configuration schema
 */
export interface CLIConfig {
    token: string;
    tokenType: 'sensor';
    accountId: string;
    environment: 'PROD' | 'DEV';
}
/**
 * Get the stored token
 */
export declare function getToken(): string | undefined;
/**
 * Save token after successful login
 */
export declare function setToken(token: string, accountId: string): void;
/**
 * Clear token on logout
 */
export declare function clearToken(): void;
/**
 * Check if user is logged in
 */
export declare function isLoggedIn(): boolean;
/**
 * Get account ID
 */
export declare function getAccountId(): string | undefined;
/**
 * Get the current environment
 */
export declare function getEnvironment(): 'PROD' | 'DEV';
/**
 * Get config file path (for debugging)
 */
export declare function getConfigPath(): string;
//# sourceMappingURL=index.d.ts.map