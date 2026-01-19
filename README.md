# whooktown CLI

`wt` is a command-line interface for [whooktown](https://whook.town), a platform that visualizes IT infrastructure as a 3D virtual city.

## Installation

### From Source

```bash
cd cli/whooktown-cli
go build -o wt ./cmd/wt

# Optional: move to PATH
sudo mv wt /usr/local/bin/
```

### Shell Completion

```bash
# Bash
wt completion bash > /etc/bash_completion.d/wt

# Zsh
wt completion zsh > "${fpath[1]}/_wt"

# Fish
wt completion fish > ~/.config/fish/completions/wt.fish
```

## Authentication

### Email Login (Interactive)

```bash
wt login
```

1. Enter your email address
2. Check your email and click the validation link
3. CLI automatically detects validation and saves token

### Token Authentication

For automation and scripts, use the `--token` flag:

```bash
wt --token <your-token> workflow list
```

Or set it in your config file (see Configuration below).

### Logout

```bash
wt logout
```

## Commands

### Global Flags

| Flag | Description |
|------|-------------|
| `--config <path>` | Config file path (default: `~/.whooktown/config.yaml`) |
| `--token <token>` | Authentication token (overrides config) |
| `--env <env>` | Environment override: `dev` or `prod` |
| `--json` | Output in JSON format |
| `--debug` | Enable debug output |

### Account

```bash
# Show current account info
wt whoami

# Login
wt login

# Logout
wt logout
```

### Token Management

```bash
# List all tokens
wt token list

# Create a new token
wt token create --name "My Script" --type sensor

# Revoke a token
wt token revoke <token>
```

Token types: `user`, `sensor`, `viewer`

### Layouts

```bash
# List layouts and quota usage
wt layout list
wt layout quota

# Show layout details
wt layout show <layout_id>

# Create layout from JSON file
wt layout create -f layout.json

# Create layout with name only
wt layout create --name "My City"

# Update layout
wt layout update <layout_id> -f layout.json

# Delete layout
wt layout delete <layout_id>

# Archived layouts
wt layout archive list
wt layout restore <layout_id>
```

### Workflows

```bash
# List workflows
wt workflow list

# Show workflow details
wt workflow show <workflow_id>

# Create workflow from JSON file
wt workflow create -f workflow.json

# Delete workflow
wt workflow delete <workflow_id>

# Enable/disable workflow
wt workflow enable <workflow_id>
wt workflow disable <workflow_id>

# Export workflow as JSON
wt workflow export <workflow_id> > workflow.json

# List available operations
wt workflow operations
```

### Sensors

Send sensor data to update building states:

```bash
# Basic status update
wt sensor send --id <building_id> --status online --activity fast

# From JSON file
wt sensor send -f sensor-data.json

# With extra fields (for DataCenter, Arcade, etc.)
wt sensor send --id <id> --status online \
  --extra cpuUsage=75 \
  --extra temperature=42
```

Status values: `online`, `offline`, `warning`, `critical`

Activity values: `slow`, `normal`, `fast`

### Camera

```bash
# Presets
wt camera preset list <layout_id>
wt camera preset create --layout <id> --name "Overview" \
  --position "10,5,10" --rotation "0,45,0"
wt camera preset delete <preset_id>
wt camera preset set-default <preset_id>

# Paths
wt camera path list <layout_id>
wt camera path create --layout <id> --name "Tour" --loop
wt camera path delete <path_id>

# Commands
wt camera command --layout <id> --mode orbit
wt camera command --layout <id> --preset <preset_id>
wt camera command --layout <id> --path <path_id>
```

Camera modes: `orbit`, `fps`, `flyover`

### Traffic

```bash
# Get traffic state
wt traffic get

# Set traffic state
wt traffic set --layout <id> --density 50 --speed normal --enabled
wt traffic set --layout <id> --density 0 --enabled=false
```

Speed values: `slow`, `normal`, `fast`

### Audio

```bash
# Get audio state
wt audio get

# Set audio state
wt audio set --layout <id> --mood tension
wt audio set --layout <id> --music-volume 80 --sfx-volume 60
wt audio set --layout <id> --enabled=false
```

Mood values: `calm`, `active`, `tension`, `critical`, `epic`

### Configuration

```bash
# Show current config
wt config show

# Set config value
wt config set environment DEV
wt config set default_layout <layout_id>

# Multi-environment support
wt config use-context dev
wt config use-context prod
wt config get-contexts
```

## Configuration File

The config file is stored at `~/.whooktown/config.yaml`:

```yaml
current_context: default
contexts:
  default:
    name: default
    token: "eyJ..."
    environment: PROD
    default_layout: ""
  dev:
    name: dev
    token: "eyJ..."
    environment: DEV
    # Optional: custom URLs for self-hosted deployments
    auth_url: "https://auth.my-server.com"
    ui_url: "https://api.my-server.com"
```

## Environment Override

Use `--env` to temporarily switch environments without modifying config:

```bash
# Use dev environment for this command only
wt --env dev workflow list

# Use prod environment
wt --env prod layout quota
```

## JSON Output

All commands support `--json` for scripting:

```bash
# Get workflows as JSON
wt workflow list --json

# Use with jq
wt workflow list --json | jq '.[].name'

# Get quota info
wt layout quota --json | jq '.layouts.used'
```

## Examples

### Deploy a new layout

```bash
# Create layout
wt layout create -f my-city.json

# Enable traffic
wt traffic set --layout <id> --density 50 --enabled

# Start audio
wt audio set --layout <id> --mood calm --enabled
```

### Monitor a service

```bash
# Send status from a script
wt sensor send --id <building_id> --status online --activity normal

# Send critical alert
wt sensor send --id <building_id> --status critical
```

### Automate with workflows

```bash
# Export existing workflow
wt workflow export <id> > backup.json

# Modify and re-import
wt workflow create -f backup.json

# Enable it
wt workflow enable <new_id>
```

### CI/CD Integration

```bash
#!/bin/bash
# Update building status based on deployment result

BUILDING_ID="your-building-uuid"
TOKEN="your-sensor-token"

if deploy_succeeded; then
  wt --token $TOKEN sensor send --id $BUILDING_ID --status online
else
  wt --token $TOKEN sensor send --id $BUILDING_ID --status critical
fi
```

## Troubleshooting

### "not logged in" error

```bash
wt login
# or use --token flag
```

### "invalid environment" error

Use `dev` or `prod`:
```bash
wt --env dev ...  # correct
wt --env DEV ...  # also correct
```

### Connection errors

Check your environment:
```bash
wt config show
wt whoami --debug
```

## License

MIT License - see [LICENSE](LICENSE) file.
