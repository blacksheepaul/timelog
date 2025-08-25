package main

import (
	"os"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/model"
)

// NewTimelogMCPServer creates and initializes a new TimelogMCPServer instance
func NewTimelogMCPServer() *TimelogMCPServer {
	// Initialize configuration
	// Check for config path from environment variable first
	configPath := os.Getenv("TIMELOG_CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yml" // fallback to default relative path
	}
	
	cfg := config.GetConfig(configPath)
	
	// Disable ORM logging to prevent ANSI escape codes in MCP output
	cfg.Log.ORMLogLevel = 1 // Silent mode
	
	// Initialize database connection using existing DAO pattern
	// Pass nil for logger to avoid stdout output that interferes with MCP protocol
	model.InitDao(cfg, nil)
	dao := model.GetDao()

	return &TimelogMCPServer{
		db: dao.Db(),
	}
}