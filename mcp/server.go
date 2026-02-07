package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var server *TimelogMCPServer

func main() {
	server = NewTimelogMCPServer()

	// Create MCP server with implementation
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "timelog",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_timelogs_by_date_range",
		Description: "Get time logs within a specific date range",
	}, GetTimeLogsByDateRange)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_tasks_by_status",
		Description: "Get tasks filtered by completion status",
	}, GetTasksByStatus)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_current_activity",
		Description: "Get currently active/running time logs",
	}, GetCurrentActivity)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_active_constraints",
		Description: "To know self discipline and external conditions",
	}, GetActiveConstraints)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_date_info",
		Description: "Get current date, time, today, yesterday, and this week's date range",
	}, GetDateInfo)

	transportMode := server.config.MCP.Transport
	switch transportMode {
	case "http":
		listenAddr := server.config.MCP.ListenAddr
		token := server.config.MCP.Token

		handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
			return mcpServer
		}, nil)

		wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
				return
			}

			if token != "" {
				authHeader := r.Header.Get("Authorization")
				if !strings.HasPrefix(authHeader, "Bearer ") {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				provided := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
				if provided != token {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			handler.ServeHTTP(w, r)
		})

		httpServer := &http.Server{
			Addr:    listenAddr,
			Handler: wrappedHandler,
		}

		if err := httpServer.ListenAndServe(); err != nil {
			LogMCPError("http_server", err, map[string]interface{}{"addr": listenAddr})
		}
	default:
		// Run MCP server - no logging to avoid stdout contamination
		ctx := context.Background()
		transport := &mcp.StdioTransport{}
		if err := mcpServer.Run(ctx, transport); err != nil {
			LogMCPError("stdio_server", err, nil)
		}
	}
}
