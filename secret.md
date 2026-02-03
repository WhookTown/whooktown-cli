# Secret Configuration

Internal documentation for developers. Not for end-users.

## Environment Variables

These environment variables override default behavior:

| Variable | Description | Default |
|----------|-------------|---------|
| `WHOOKTOWN_ENV` | Environment: `PROD` or `DEV` | `PROD` |
| `WHOOKTOWN_AUTH_URL` | Override auth service URL | `https://auth.whook.town` |
| `WHOOKTOWN_SENSOR_URL` | Override sensor endpoint URL | `https://sensors.whook.town` |

## Usage

### Development Environment

```bash
# Use dev servers
export WHOOKTOWN_ENV=DEV
wt sensor list
```

This switches to:
- Auth: `https://auth.dev.whook.town`
- Sensors: `https://sensors.dev.whook.town`

### Custom URLs

```bash
# Local development
export WHOOKTOWN_AUTH_URL=http://localhost:8981
export WHOOKTOWN_SENSOR_URL=http://localhost:8081
wt login my-token
wt sensor list
```

### CI/CD

```bash
# In CI pipeline
WHOOKTOWN_AUTH_URL=http://auth:8981 \
WHOOKTOWN_SENSOR_URL=http://sensors:8081 \
wt sensor send ci-test --status online -q
```

## Production URLs

| Service | Production URL |
|---------|----------------|
| Auth | `https://auth.whook.town` |
| Sensors | `https://sensors.whook.town` |
| UI/API | `https://api.whook.town` |
| WebSocket | `https://ws.whook.town` |

## Development URLs

| Service | Development URL |
|---------|-----------------|
| Auth | `https://auth.dev.whook.town` |
| Sensors | `https://sensors.dev.whook.town` |
| UI/API | `https://api.dev.whook.town` |
| WebSocket | `https://ws.dev.whook.town` |

## Token Security

- Tokens are stored in `~/.config/whooktown/config.json`
- File permissions: user-only read/write (handled by `conf` package)
- Never commit tokens to version control
- Rotate tokens periodically
- Revoke tokens via the web dashboard when no longer needed
