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

### 2. Choose Transport Mode

The MCP server supports two transport modes:

| Mode | Use Case | Configuration |
|------|----------|---------------|
| `stdio` (default) | Local LLM clients (Claude Code, Claude Desktop) | Set `transport: stdio` or leave unset |
| `http` (streamable) | Remote access, web clients, or HTTP-based integrations | Set `transport: http` and configure `listen_addr` |

#### Option A: stdio Mode (Default)

For local LLM clients like Claude Code or Claude Desktop.

Add to your Claude Code configuration file:

- **macOS**: `~/Library/Application Support/claude-code/config.json`
- **Linux**: `~/.config/claude-code/config.json`
- **Windows**: `%APPDATA%\\claude-code\config.json`

```json
{
  "mcp": {
    "servers": {
      "timelog": {
        "command": "/absolute/path/to/your/timelog/mcp/timelog-mcp-server",
        "args": [],
        "env": {
          "TIMELOG_CONFIG_PATH": "/absolute/path/to/your/config.yml",
          "MCP_TRANSPORT": "stdio"
        }
      }
    }
  }
}
```

Or configure in `config.yml`:

```yaml
mcp:
  transport: stdio
```

#### Option B: HTTP Mode (Streamable Transport)

For remote access or HTTP-based clients. The server runs as a standalone HTTP service.

**Configuration via environment variables:**

```bash
export TIMELOG_CONFIG_PATH="/absolute/path/to/your/config.yml"
export MCP_TRANSPORT="http"
export MCP_LISTEN_ADDR=":8080"      # Optional, defaults to :8080
export MCP_TOKEN="your-secret-token" # Optional, for authentication
./timelog-mcp-server
```

**Configuration via config.yml:**

```yaml
mcp:
  transport: http
  listen_addr: ":8080"      # Listen address for HTTP server
  token: "your-secret-token" # Optional Bearer token for authentication
```

**Running the HTTP server:**

```bash
# Start the server
./timelog-mcp-server

# Server will output:
# HTTP server listening on :8080
```

**HTTP Endpoints:**

- `POST /` - MCP protocol endpoint (JSON-RPC)
- `GET /health` - Health check endpoint

**Authentication:**

If `token` is configured, all requests to `/` must include:

```
Authorization: Bearer your-secret-token
```

**Connecting clients to HTTP server:**

For clients that support MCP HTTP transport (like Claude Desktop with HTTP configuration):

```json
{
  "mcp": {
    "servers": {
      "timelog": {
        "url": "http://localhost:8080",
        "headers": {
          "Authorization": "Bearer your-secret-token"
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
- **HTTP mode connection refused**: Check if `MCP_LISTEN_ADDR` is set correctly and port is not in use
- **HTTP 401 Unauthorized**: Verify your `MCP_TOKEN` matches between server and client
