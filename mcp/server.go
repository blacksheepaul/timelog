package main

import (
	"context"

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

	// Run MCP server - no logging to avoid stdout contamination
	ctx := context.Background()
	transport := &mcp.StdioTransport{}

	// Run server (any error will exit the process)
	mcpServer.Run(ctx, transport)
}
