# TimeLog MCP Server Setup Guide

## Quick Setup

### 1. Build the Server
```bash
# From the timelog project root
make mcp-server
```

### 2. Add to Claude Code Configuration

Edit your Claude Code configuration file:
- **macOS**: `~/Library/Application Support/claude-code/config.json`
- **Linux**: `~/.config/claude-code/config.json`
- **Windows**: `%APPDATA%\\claude-code\\config.json`

Add this configuration:
```json
{
  "mcp": {
    "servers": {
      "timelog-analyzer": {
        "command": "/Users/n/Documents/Developer/timelog/mcp/timelog-mcp-server",
        "args": [],
        "env": {}
      }
    }
  }
}
```

**Important**: Update the `command` path to match your actual project location.

### 3. Restart Claude Code

Restart Claude Code to load the MCP server.

### 4. Test the Integration

Try these example queries with Claude Code:

```
"Show me my recent time logs"
"Analyze my productivity over the last week"
"What tasks are pending?"
"Give me task completion stats for the last month"
"Are there any active time logs running?"
```

## Available Tools

1. **get_recent_timelogs** - Recent time log entries
2. **get_timelogs_by_date_range** - Time logs within date ranges
3. **get_tasks_by_status** - Tasks filtered by status
4. **get_productivity_stats** - Productivity analysis
5. **get_task_completion_analysis** - Task efficiency metrics
6. **get_current_activity** - Active time tracking sessions

## Troubleshooting

**Server not found**: Verify the command path in your config points to the built binary.

**Permission denied**: Ensure the binary has execute permissions:
```bash
chmod +x /path/to/timelog/mcp/timelog-mcp-server
```

**Database errors**: Verify your `config.yml` exists and the TimeLog app works normally.

## Example Behavior Analysis Queries

Once integrated, you can ask Claude Code questions like:

- "How productive was I this week compared to last week?"
- "What activities am I spending the most time on?"
- "Show me my task completion patterns"
- "Am I meeting my time estimates for tasks?"
- "What's my average daily work time?"
- "Which tags show the best productivity?"

The MCP server provides Claude Code with direct access to your TimeLog data for comprehensive behavior analysis and productivity insights.