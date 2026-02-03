import Conf from 'conf';

/**
 * CLI configuration schema
 */
export interface CLIConfig {
  token: string;
  tokenType: 'sensor';
  accountId: string;
  environment: 'PROD' | 'DEV';
}

const config = new Conf<CLIConfig>({
  projectName: 'whooktown',
  defaults: {
    token: '',
    tokenType: 'sensor',
    accountId: '',
    environment: 'PROD',
  },
});

/**
 * Get the stored token
 */
export function getToken(): string | undefined {
  const token = config.get('token');
  return token || undefined;
}

/**
 * Save token after successful login
 */
export function setToken(token: string, accountId: string): void {
  config.set('token', token);
  config.set('tokenType', 'sensor');
  config.set('accountId', accountId);
}

/**
 * Clear token on logout
 */
export function clearToken(): void {
  config.delete('token');
  config.delete('accountId');
}

/**
 * Check if user is logged in
 */
export function isLoggedIn(): boolean {
  return Boolean(config.get('token'));
}

/**
 * Get account ID
 */
export function getAccountId(): string | undefined {
  const accountId = config.get('accountId');
  return accountId || undefined;
}

/**
 * Get the current environment
 */
export function getEnvironment(): 'PROD' | 'DEV' {
  // Check env var first
  if (process.env.WHOOKTOWN_ENV === 'DEV') {
    return 'DEV';
  }
  return config.get('environment') || 'PROD';
}

/**
 * Get config file path (for debugging)
 */
export function getConfigPath(): string {
  return config.path;
}
