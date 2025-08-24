# TimeLog MCP Server

This directory contains the Model Context Protocol (MCP) server for the TimeLog application, enabling Claude Code to analyze your behavior patterns, productivity metrics, and time tracking data.

## Overview

The TimeLog MCP Server provides Claude Code with direct access to your TimeLog database through a set of specialized tools designed for behavior analysis and productivity insights. This allows you to ask Claude Code questions about your work patterns, task completion rates, time allocation, and overall productivity trends.

## Features

### Available Tools

1. **get_recent_timelogs** - Get your most recent time log entries
2. **get_timelogs_by_date_range** - Analyze time logs within specific date ranges
3. **get_tasks_by_status** - Review tasks by completion status
4. **get_productivity_stats** - Comprehensive productivity pattern analysis
5. **get_task_completion_analysis** - Task completion efficiency metrics
6. **get_current_activity** - Check currently active/running time logs

### Analytics Capabilities

- **Time Allocation Analysis**: See how you spend time across different tags/categories
- **Productivity Patterns**: Daily, weekly, and custom period productivity metrics
- **Task Completion Rates**: Track task completion efficiency by tag and time period
- **Trend Analysis**: Identify patterns in work habits and productivity
- **Current Status**: Monitor active time tracking sessions

## Setup Instructions

### 1. Build the MCP Server

```bash
# From the project root directory
make mcp-server

# Or manually
cd mcp && go build -o timelog-mcp-server server.go
```

### 2. Configure Claude Code

Add the server configuration to your Claude Code settings. The configuration file location depends on your system:

**macOS**: `~/Library/Application Support/claude-code/config.json`
**Linux**: `~/.config/claude-code/config.json`
**Windows**: `%APPDATA%\\claude-code\\config.json`

Add this configuration to your `config.json`:

```json
{
  "mcp": {
    "servers": {
      "timelog-analyzer": {
        "command": "/path/to/your/timelog/mcp/timelog-mcp-server",
        "args": [],
        "env": {}
      }
    }
  }
}
```

**Important**: Replace `/path/to/your/timelog/mcp/timelog-mcp-server` with the absolute path to your built MCP server binary.

### 3. Restart Claude Code

After adding the configuration, restart Claude Code for the changes to take effect.

### 4. Verify Connection

You can verify the MCP server is working by asking Claude Code questions like:
- "Show me my recent time logs"
- "Analyze my productivity over the last week"
- "What are my current active time tracking sessions?"

## Usage Examples

### Basic Time Log Queries

```
# Recent activity
"Show me my last 20 time log entries"

# Date range analysis  
"Analyze my time logs from 2024-01-01 to 2024-01-31"

# Current status
"What time logs are currently running?"
```

### Productivity Analysis

```
# Weekly productivity
"Give me productivity stats for the last 7 days"

# Monthly overview
"Analyze my productivity patterns over the last 30 days"

# Tag-based analysis
"How am I allocating time across different activities?"
```

### Task Management Insights

```
# Task status
"Show me all pending tasks"

# Completion analysis
"Analyze my task completion rates over the last month"

# Efficiency metrics
"How well am I meeting my task estimates?"
```

## Architecture

The MCP server is built in Go and leverages your existing TimeLog infrastructure:

- **Database**: Uses the same SQLite database and GORM models as your main application
- **Configuration**: Reads from the same `config.yml` file
- **Models**: Reuses existing `TimeLog`, `Task`, and `Tag` models
- **Business Logic**: Implements specialized analytics functions

### Dependencies

- **Go MCP SDK**: `github.com/modelcontextprotocol/go-sdk`
- **Existing TimeLog Models**: Reuses your current database models and configuration

## Tool Reference

### get_recent_timelogs

Get your most recent time log entries.

**Parameters:**
- `limit` (optional): Maximum number of entries to return (default: 10)

**Returns:** List of time logs with calculated durations, tag information, and associated tasks.

### get_timelogs_by_date_range

Analyze time logs within a specific date range.

**Parameters:**
- `start_date` (required): Start date in YYYY-MM-DD format
- `end_date` (required): End date in YYYY-MM-DD format

**Returns:** Time logs within the range plus total duration statistics.

### get_tasks_by_status

Review tasks filtered by completion status.

**Parameters:**
- `status` (required): One of "completed", "pending", or "all"

**Returns:** Filtered task list with tag information and completion details.

### get_productivity_stats

Comprehensive productivity pattern analysis.

**Parameters:**
- `days` (optional): Number of days to analyze (default: 7, max: 365)

**Returns:** Daily breakdown, tag allocation, averages, and productivity metrics.

### get_task_completion_analysis  

Task completion efficiency and pattern analysis.

**Parameters:**
- `days` (optional): Number of days to analyze (default: 30, max: 365)

**Returns:** Completion rates, overdue tasks, estimated vs. actual time, and tag-based completion patterns.

### get_current_activity

Check currently active/running time logs.

**Parameters:** None

**Returns:** List of time logs that haven't been ended yet, with current running durations.

## Troubleshooting

### Common Issues

1. **Server not found**: Ensure the path in your Claude Code config points to the built binary
2. **Database connection errors**: Verify your `config.yml` file is accessible and database exists
3. **Permission errors**: Ensure the MCP server binary has execute permissions

### Debug Steps

1. Test the server manually:
   ```bash
   cd /path/to/timelog/mcp
   ./timelog-mcp-server
   ```

2. Check Claude Code logs for MCP-related errors

3. Verify your TimeLog application is working normally

### Configuration Verification

You can use the provided `claude-code-config.json` as a template, but remember to update the path to match your system:

```bash
# Copy and modify the template
cp mcp/claude-code-config.json ~/.config/claude-code/config.json
# Edit the file to update the path to your binary
```

## Security Considerations

- The MCP server runs locally and only provides read-only access to your TimeLog data
- No network connections are made; all communication is through stdin/stdout
- Your data remains on your local machine and is not transmitted externally
- The server inherits the same file permissions as your TimeLog application

## Performance Notes

- The server is designed to handle typical personal time logging data volumes efficiently
- Database queries are optimized for common analysis patterns
- For very large datasets (thousands of time logs), consider using date range filters to improve response times

## Contributing

To extend the MCP server with additional tools:

1. Add new tool functions following the existing patterns in `server.go`
2. Register the tools with appropriate schemas in the `main()` function
3. Update this documentation with the new tool details
4. Test thoroughly with your actual TimeLog data

## Support

For issues specific to the MCP server, check:
1. Your TimeLog application is working correctly
2. The MCP server binary builds and runs successfully
3. Claude Code configuration is correct and the server is restarting properly

For general TimeLog application issues, refer to the main project documentation.