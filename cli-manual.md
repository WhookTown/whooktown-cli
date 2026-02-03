# Whooktown CLI Manual

Complete command reference for the `wt` command-line tool.

## Authentication

### wt login

Login with a sensor token.

```bash
wt login <token> [options]

Arguments:
  token              Sensor token to authenticate with

Options:
  --no-validate      Skip token validation (not recommended)
```

The token must be:
- Type: `sensor`
- Roles: `{"sensor": "rw"}`

**Examples:**

```bash
# Login with validation
wt login eyJhbGciOiJIUzI1NiIs...

# Login without validation (development only)
wt login my-token --no-validate
```

### wt logout

Clear saved token.

```bash
wt logout
```

---

## Sensors

### wt sensor send

Send sensor data to update an asset's state.

```bash
wt sensor send <id> [options]

Arguments:
  id                 Sensor ID (asset identifier)

Options:
  -s, --status <status>     Status: online, offline, warning, critical
  -a, --activity <activity> Activity: slow, normal, fast
  -j, --json <json>         Additional JSON fields
  -q, --quiet               Suppress output
```

**Examples:**

```bash
# Set status and activity
wt sensor send my-server --status online --activity fast

# Send with additional data
wt sensor send datacenter-1 --status warning --json '{"cpuUsage": 85, "ramUsage": 70}'

# Update status only
wt sensor send web-server --status critical

# Quiet mode for scripting
wt sensor send my-sensor --status online -q
```

### wt sensor list

List all sensor states.

```bash
wt sensor list [options]

Options:
  -f, --format <format>     Output format: table, json (default: table)
```

**Examples:**

```bash
# Table output
wt sensor list

# JSON output for scripting
wt sensor list --format json
```

---

## Traffic

### wt traffic set

Update traffic settings for a layout.

```bash
wt traffic set <layoutId> [options]

Arguments:
  layoutId           Layout ID

Options:
  -d, --density <density>   Traffic density (0-100)
  -s, --speed <speed>       Speed: slow, normal, fast
  --enabled                 Enable traffic
  --disabled                Disable traffic
```

At least one option must be provided. Unspecified options retain their current values.

**Examples:**

```bash
# Set density
wt traffic set my-layout --density 75

# Enable traffic with settings
wt traffic set my-layout --density 50 --speed normal --enabled

# Disable traffic
wt traffic set my-layout --disabled

# Set speed only
wt traffic set my-layout --speed fast
```

### wt traffic list

List traffic states for all layouts.

```bash
wt traffic list [options]

Options:
  -f, --format <format>     Output format: table, json (default: table)
```

---

## Popup & Building Metadata

### wt popup labels

Toggle building labels visibility for a layout.

```bash
wt popup labels <layoutId> [options]

Arguments:
  layoutId           Layout ID

Options:
  --on               Enable labels
  --off              Disable labels
```

One of `--on` or `--off` must be provided.

**Examples:**

```bash
# Show labels
wt popup labels my-layout --on

# Hide labels
wt popup labels my-layout --off
```

### wt popup set

Set building metadata (description, tags, notes).

```bash
wt popup set <buildingId> [options]

Arguments:
  buildingId         Building ID (UUID)

Options:
  -d, --description <text>    Set description
  -t, --tags <tags>           Set comma-separated tags
  -n, --notes <text>          Set notes
  --clear-description         Clear description
  --clear-tags                Clear tags
  --clear-notes               Clear notes
```

At least one option must be provided.

**Examples:**

```bash
# Set description
wt popup set abc-123 --description "Production web server"

# Set tags (comma-separated)
wt popup set abc-123 --tags "web,production,critical"

# Set multiple fields
wt popup set abc-123 -d "Main database" -t "database,primary" -n "SLA critical"

# Clear a field
wt popup set abc-123 --clear-tags
```

### wt popup get

View building metadata.

```bash
wt popup get <buildingId> [options]

Arguments:
  buildingId         Building ID (UUID)

Options:
  -f, --format <format>       Output format: text, json (default: text)
```

**Examples:**

```bash
# Text output (default)
wt popup get abc-123

# JSON output
wt popup get abc-123 --format json
```

### wt popup list

List all buildings with their metadata.

```bash
wt popup list <layoutId> [options]

Arguments:
  layoutId           Layout ID

Options:
  -f, --format <format>       Output format: table, json (default: table)
  --tags <filter>             Filter by tags (comma-separated, any match)
```

**Examples:**

```bash
# List all buildings
wt popup list my-layout

# JSON output
wt popup list my-layout --format json

# Filter by tags
wt popup list my-layout --tags "production,critical"
```

---

## Interactive TUI

### wt tui

Launch the interactive terminal dashboard.

```bash
wt tui [options]

Options:
  -r, --refresh <ms>        Auto-refresh interval in milliseconds (default: 5000)
```

**Keyboard Shortcuts:**

| Key | Action |
|-----|--------|
| `1` | Switch to Sensors panel |
| `2` | Switch to Traffic panel |
| `r` | Manual refresh |
| `q` | Quit |
| `Ctrl+C` | Quit |

**Examples:**

```bash
# Default (5 second refresh)
wt tui

# Fast refresh (2 seconds)
wt tui --refresh 2000
```

---

## Output Formats

Commands that list data support two output formats:

### Table (default)

Human-readable table format with colors.

```bash
wt sensor list
# ID                    Status    Activity   Updated
# ─────────────────────────────────────────────────
# server-1              online    fast       10:30:45
# server-2              warning   slow       10:30:42
```

### JSON

Machine-readable JSON for scripting and piping.

```bash
wt sensor list --format json
# [{"sid":"server-1","data":{"status":"online","activity":"fast"},...}]
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (authentication, validation, network) |

---

## Configuration File

Location: `~/.config/whooktown/config.json`

The CLI stores:
- Authentication token
- Account ID
- Environment setting

To view config path:
```bash
wt login --help
# Shows config path on successful login
```
