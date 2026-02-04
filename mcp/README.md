# TimeLog MCP Server

MCP server that gives LLM access to your TimeLog data for analysis and insights.

## What You Can Do

- Ask LLM about your work patterns and productivity trends
- Analyze time allocation across different activities
- Review task completion rates and efficiency
- Monitor current active time tracking sessions
- Get productivity insights for any time period
- Anything else based on your timelog, use your imagination to the fullest

## Setup

### 1. Build the MCP Server

```bash
# From the project root directory
make mcp
```

### 2. Configure your LLM client

Example for Claude Code as below.

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

### 3. (optional) Enable Debug Logging

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

Ask from your LLM client:

- "Show me my recent time logs"
- "Analyze my productivity this week"
- "What tasks are currently pending?"

## What You Can Ask

- "Show me my recent time logs"
- "Analyze my productivity this week"
- "What tasks are pending?"
- "How am I spending time across different activities?"
- "What's my task completion rate this month?"
- "Are there any active time logs running?"
- "What is the best I can do under the existing constraints today"
- ...

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
3. **Test manually and watch your logs!**

After my personal testing, the support for Claude Desktop, Claude Code, VSCode and Cherry Studio is very good, but opencat and jan need some prompt engineer.

### Common Issues:

- **Config not found**: Make sure you have defined environment variable `TIMELOG_CONFIG_PATH`
- **Server not found**: Use absolute paths in Claude Code config
- **Database errors**: Ensure your TimeLog app works normally
- **Permission errors**: Make sure the binary is executable (`chmod +x timelog-mcp-server`)
