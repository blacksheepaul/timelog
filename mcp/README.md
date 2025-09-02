# TimeLog MCP Server

MCP server that gives Claude Code access to your TimeLog data for productivity analysis and insights.

## What You Can Do

- Ask Claude Code about your work patterns and productivity trends
- Analyze time allocation across different activities
- Review task completion rates and efficiency
- Monitor current active time tracking sessions
- Get productivity insights for any time period

## Setup

### 1. Build the MCP Server

```bash
# From the project root directory
cd mcp && go build -o timelog-mcp-server *.go
```

### 2. Configure Claude Code

Add to your Claude Code configuration file:

- **macOS**: `~/Library/Application Support/claude-code/config.json`
- **Linux**: `~/.config/claude-code/config.json`
- **Windows**: `%APPDATA%\\claude-code\\config.json`

```json
{
  "mcp": {
    "servers": {
      "timelog": {
        "command": "/absolute/path/to/your/timelog/mcp/timelog-mcp-server",
        "args": [],
        "env": {
          "TIMELOG_CONFIG_PATH": "/absolute/path/to/your/config.yml"
        }
      }
    }
  }
}
```

### 3. Enable Debug Logging (for troubleshooting)

Add MCP logging configuration to your `config.yml`:

```yaml
log:
  level: info
  path: logs/app.log
  rotation:
    max_size: 100
    max_backups: 5
    max_age: 30
  # Enable MCP debug logging
  mcp:
    enabled: true
    level: debug
    path: logs/mcp.log
```

### 4. Test the Setup

Ask Claude Code:
- "Show me my recent time logs"
- "Analyze my productivity this week" 
- "What tasks are currently pending?"

## What You Can Ask Claude Code

- "Show me my recent time logs"
- "Analyze my productivity this week"
- "What tasks are pending?"
- "How am I spending time across different activities?"
- "What's my task completion rate this month?"
- "Are there any active time logs running?"

## Troubleshooting

### If MCP server doesn't work:

1. **Check the binary path** in your Claude Code config file
2. **Enable debug logging** in your `config.yml`:
   ```yaml
   log:
     mcp:
       enabled: true
       level: debug
       path: logs/mcp.log
   ```
3. **Test manually**:
   ```bash
   cd /path/to/timelog/mcp
   ./timelog-mcp-server
   ```
4. **Check the debug log** at `logs/mcp.log` for error details

### Common Issues:

- **Server not found**: Use absolute paths in Claude Code config
- **Database errors**: Ensure your TimeLog app works normally
- **Permission errors**: Make sure the binary is executable (`chmod +x timelog-mcp-server`)