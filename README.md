# Whooktown CLI

Command-line interface for Whooktown - 3D IT Infrastructure Visualization.

## Installation

```bash
# Build
npm install
npm run build

# Link globally
npm link
```

## Quick Start

```bash
# Login with sensor token
wt login <your-sensor-token>

# List sensors (assets)
wt sensor list

# Send sensor data
wt sensor send my-sensor-id --status online --activity fast

# List camera states
wt camera list

# Set camera mode
wt camera set my-layout --mode flyover --speed 3.0

# List traffic states
wt traffic list

# Update traffic
wt traffic set my-layout --density 75 --speed fast --enabled

# Launch interactive TUI
wt tui
```

## Commands

| Command | Description |
|---------|-------------|
| `wt login <token>` | Login with sensor token |
| `wt logout` | Clear saved token |
| `wt sensor send <id>` | Send sensor data |
| `wt sensor list` | List sensor states |
| `wt camera set <layoutId>` | Set camera mode |
| `wt camera list` | List camera states |
| `wt traffic set <layoutId>` | Update traffic settings |
| `wt traffic list` | List traffic states |
| `wt tui` | Launch interactive dashboard |

## Token Requirements

This CLI only accepts **sensor tokens** (type: `sensor`, roles: `{"sensor": "rw"}`).

To get a sensor token:
1. Go to your Whooktown dashboard
2. Navigate to Settings > API Tokens
3. Create a new token with type "Sensor"

## Configuration

Configuration is stored in `~/.config/whooktown/config.json`.

## Documentation

See [cli-manual.md](./cli-manual.md) for complete command reference.
