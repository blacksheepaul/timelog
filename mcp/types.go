package main

import (
	"gorm.io/gorm"
)

// TimelogMCPServer is the main server struct
type TimelogMCPServer struct {
	db *gorm.DB
}

// Tool parameter structs
type RecentTimeLogsParams struct {
	Limit int `json:"limit,omitempty" jsonschema:"Maximum number of time logs to return (default: 10)"`
}

type DateRangeParams struct {
	StartDate string `json:"start_date" jsonschema:"Start date in YYYY-MM-DD format,required"`
	EndDate   string `json:"end_date" jsonschema:"End date in YYYY-MM-DD format,required"`
}

type TaskStatusParams struct {
	Status string `json:"status" jsonschema:"Task status filter (completed/pending/all),required"`
}

type StatsParams struct {
	Days int `json:"days,omitempty" jsonschema:"Number of days to analyze (default: 7)"`
}

type TaskAnalysisParams struct {
	Days int `json:"days,omitempty" jsonschema:"Number of days to analyze (default: 30)"`
}

type CurrentActivityParams struct {
	// No parameters needed
}