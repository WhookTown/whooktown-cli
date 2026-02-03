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

# List layouts
wt layout list

# List traffic states
wt traffic list

# Update traffic
wt traffic set my-layout --density 75 --speed fast --enabled

# Toggle building labels
wt popup labels my-layout --on

# View building metadata
wt popup get building-uuid

# Set building metadata
wt popup set building-uuid --description "Web server" --tags "web,prod"

# List buildings with metadata
wt popup list my-layout

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
| `wt layout list` | List layouts with building counts |
| `wt traffic set <layoutId>` | Update traffic settings |
| `wt traffic list` | List traffic states |
| `wt popup labels <layoutId>` | Toggle building labels |
| `wt popup set <buildingId>` | Set description/tags/notes |
| `wt popup get <buildingId>` | View building metadata |
| `wt popup list <layoutId>` | List buildings with metadata |
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
