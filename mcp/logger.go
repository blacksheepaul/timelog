package main

import (
	"os"
	"path/filepath"

	"github.com/blacksheepaul/timelog/core/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var mcpLogger *zap.SugaredLogger

// InitMCPLogger creates a file-only logger for MCP debugging
// This logger will never write to stdout to avoid breaking MCP protocol
func InitMCPLogger(cfg *config.Config) *zap.SugaredLogger {
	// Return no-op logger if MCP logging is disabled
	if cfg == nil || !cfg.MCP.Enabled {
		mcpLogger = zap.NewNop().Sugar()
		return mcpLogger
	}

	// Allow environment variable override
	if os.Getenv("MCP_DEBUG") == "false" {
		mcpLogger = zap.NewNop().Sugar()
		return mcpLogger
	}

	// Parse log level (default to debug for troubleshooting)
	level, err := zap.ParseAtomicLevel(cfg.MCP.Level)
	if err != nil {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	// Ensure log directory exists
	logPath := cfg.MCP.Path
	if logPath == "" {
		logPath = "logs/mcp.log"
	}

	logDir := filepath.Dir(logPath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			// If we can't create directory, return no-op logger
			mcpLogger = zap.NewNop().Sugar()
			return mcpLogger
		}
	}

	// Configure encoder for detailed debugging
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "timestamp",
		MessageKey:     "message",
		CallerKey:      "caller",
		NameKey:        "logger",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// File writer with rotation (reuse main app rotation settings)
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    cfg.Log.Rotation.MaxSize,
		MaxBackups: cfg.Log.Rotation.MaxBackups,
		MaxAge:     cfg.Log.Rotation.MaxAge,
	})

	// Create core - FILE OUTPUT ONLY (never stdout)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriter,
		level,
	)

	// Create logger with caller info for debugging
	mcpLogger = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	).Sugar()

	return mcpLogger
}

// GetMCPLogger returns the MCP logger instance
func GetMCPLogger() *zap.SugaredLogger {
	if mcpLogger == nil {
		return zap.NewNop().Sugar()
	}
	return mcpLogger
}

// Helper functions for structured MCP logging

// LogMCPToolCall logs when an MCP tool is called
func LogMCPToolCall(toolName string, params map[string]interface{}) {
	if mcpLogger == nil {
		return
	}
	mcpLogger.Infow("MCP tool called",
		"tool", toolName,
		"params", params,
	)
}

// LogMCPToolResult logs the result of an MCP tool call
func LogMCPToolResult(toolName string, success bool, resultCount int, err error) {
	if mcpLogger == nil {
		return
	}

	fields := []interface{}{
		"tool", toolName,
		"success", success,
		"result_count", resultCount,
	}

	if err != nil {
		fields = append(fields, "error", err.Error())
		mcpLogger.Errorw("MCP tool completed with error", fields...)
	} else {
		mcpLogger.Infow("MCP tool completed successfully", fields...)
	}
}

// LogMCPError logs MCP-specific errors for debugging
func LogMCPError(operation string, err error, context map[string]interface{}) {
	if mcpLogger == nil {
		return
	}

	fields := []interface{}{
		"operation", operation,
		"error", err.Error(),
	}

	for k, v := range context {
		fields = append(fields, k, v)
	}

	mcpLogger.Errorw("MCP error occurred", fields...)
}

// LogMCPDebug logs debug information for troubleshooting
func LogMCPDebug(message string, fields map[string]interface{}) {
	if mcpLogger == nil {
		return
	}

	logFields := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		logFields = append(logFields, k, v)
	}

	mcpLogger.Debugw(message, logFields...)
}
