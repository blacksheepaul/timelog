package main

import (
	"gorm.io/gorm"
)

// TimelogMCPServer is the main server struct
type TimelogMCPServer struct {
	db *gorm.DB
}

// Tool parameter structs
type DateRangeParams struct {
	StartDate string `json:"start_date" jsonschema:"Start date in YYYY-MM-DD format,required"`
	EndDate   string `json:"end_date" jsonschema:"End date in YYYY-MM-DD format,required"`
}

type TaskStatusParams struct {
	Status string `json:"status" jsonschema:"Task status filter (completed/pending/all),required"`
}

type CurrentActivityParams struct{
	// No parameters needed
}

type ConstraintParams struct {
	// No parameters needed
}
