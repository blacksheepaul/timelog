# MCP Remote Access Guide

This guide explains how to run the MCP server over HTTP and connect to it from a remote client.

## Server Setup

1. Build the MCP server binary:

```bash
make mcp
```

2. Configure MCP for HTTP transport:

```yaml
mcp:
  transport: http
  listen_addr: ":8080"
  token: "your-secret-token"
```

3. Start the server:

```bash
./mcp/timelog-mcp-server
```

The server listens on `listen_addr` and exposes:

- `POST /` - MCP protocol endpoint
- `GET /health` - health check

## Client Setup

Configure the client to point to the HTTP endpoint and include the token.

Example (generic MCP HTTP client config):

```json
{
  "mcp": {
    "servers": {
      "timelog": {
        "url": "http://your-host:8080",
        "headers": {
          "Authorization": "Bearer your-secret-token"
        }
      }
    }
  }
}
```

## Security Notes

- Treat `token` as a secret.
- Use VPN or a private network when possible.
- Rotate the token if it is exposed.

## Troubleshooting

- **Connection refused**: verify the server is running and the port is reachable.
- **401 Unauthorized**: verify the token in server config and client header.
- **No response**: check MCP logs and ensure `MCP_TRANSPORT=http`.
