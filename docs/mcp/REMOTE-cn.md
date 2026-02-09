# MCP 远程访问指南

本指南介绍如何通过 HTTP 运行 MCP 服务器，并从远程客户端连接到它。

## 服务器端设置

### 1. 构建 MCP 服务器二进制文件

```bash
make mcp
```

### 2. 配置 HTTP 传输模式

在 `config.yml` 中添加：

```yaml
mcp:
  transport: http
  listen_addr: ":8080"
  token: "your-secret-token"
```

### 3. 启动服务器

```bash
./mcp/timelog-mcp-server
```

服务器将监听 `listen_addr` 指定的地址，并暴露以下端点：

- `POST /` - MCP 协议端点
- `GET /health` - 健康检查

## 客户端配置

配置客户端指向 HTTP 端点，并在请求头中包含认证令牌。

通用 MCP HTTP 客户端配置示例：

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

## 安全注意事项

- 将 `token` 视为机密信息妥善保管
- 尽可能使用 VPN 或私有网络
- 如令牌泄露，立即更换新令牌

## 故障排除

- **连接被拒绝**：确认服务器正在运行且端口可访问
- **401 未授权**：检查服务器配置和客户端请求头中的 token 是否一致
- **无响应**：查看 MCP 日志，确认 `MCP_TRANSPORT=http` 已正确设置
