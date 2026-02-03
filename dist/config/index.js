import Conf from 'conf';
const config = new Conf({
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
export function getToken() {
    const token = config.get('token');
    return token || undefined;
}
/**
 * Save token after successful login
 */
export function setToken(token, accountId) {
    config.set('token', token);
    config.set('tokenType', 'sensor');
    config.set('accountId', accountId);
}
/**
 * Clear token on logout
 */
export function clearToken() {
    config.delete('token');
    config.delete('accountId');
}
/**
 * Check if user is logged in
 */
export function isLoggedIn() {
    return Boolean(config.get('token'));
}
/**
 * Get account ID
 */
export function getAccountId() {
    const accountId = config.get('accountId');
    return accountId || undefined;
}
/**
 * Get the current environment
 */
export function getEnvironment() {
    // Check env var first
    if (process.env.WHOOKTOWN_ENV === 'DEV') {
        return 'DEV';
    }
    return config.get('environment') || 'PROD';
}
/**
 * Get config file path (for debugging)
 */
export function getConfigPath() {
    return config.path;
}
//# sourceMappingURL=index.js.map