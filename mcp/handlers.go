package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var singaporeLocation *time.Location

func init() {
	var err error
	singaporeLocation, err = time.LoadLocation("Asia/Singapore")
	if err != nil {
		// Fallback to UTC+8 if timezone data is not available
		singaporeLocation = time.FixedZone("SGT", 8*60*60)
	}
}

// formatMCPResponse wraps the response in the standard format to prevent LLM hallucinations
func formatMCPResponse(summaryText string, data interface{}) (*mcp.CallToolResult, interface{}, error) {
	// Add summary to data
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("data must be a map[string]interface{}")
	}
	dataMap["_summary"] = summaryText

	jsonBytes, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	responseText := fmt.Sprintf("<STRICT_JSON>\n%s\n</STRICT_JSON>", string(jsonBytes))

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{
			Text: responseText,
		}},
	}, nil, nil
}

// Tool handlers with correct MCP signature
func GetRecentTimeLogs(ctx context.Context, req *mcp.CallToolRequest, args RecentTimeLogsParams) (*mcp.CallToolResult, interface{}, error) {
	limit := args.Limit
	if limit == 0 {
		limit = 10 // default limit
	}

	timeLogs, err := model.ListTimeLogsWithOptions(server.db, limit, "created_at DESC")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get recent time logs: %w", err)
	}

	// Calculate durations and format response
	var result []map[string]interface{}
	for _, tl := range timeLogs {
		duration := ""
		if tl.EndTime != nil {
			d := tl.EndTime.Sub(tl.StartTime)
			hours := int(d.Hours())
			minutes := int(d.Minutes()) % 60
			duration = fmt.Sprintf("%dh %dm", hours, minutes)
		} else {
			// Still ongoing
			d := time.Since(tl.StartTime)
			hours := int(d.Hours())
			minutes := int(d.Minutes()) % 60
			duration = fmt.Sprintf("%dh %dm (ongoing)", hours, minutes)
		}

	entry := map[string]interface{}{
		"id":         tl.ID,
		"start_time": tl.StartTime.In(singaporeLocation).Format("2006-01-02 15:04:05"),
		"end_time":   nil,
		"duration":   duration,
		"tag":        tl.Tag.Name,
		"tag_color":  tl.Tag.Color,
		"remarks":    tl.Remark,
		"created_at": tl.CreatedAt.In(singaporeLocation).Format("2006-01-02 15:04:05"),
	}

	if tl.EndTime != nil {
		entry["end_time"] = tl.EndTime.In(singaporeLocation).Format("2006-01-02 15:04:05")
	}

	if tl.Task != nil {
		entry["task"] = map[string]interface{}{
			"id":    tl.Task.ID,
			"title": tl.Task.Title,
		}
	}

	result = append(result, entry)
	}

	response := map[string]interface{}{
		"time_logs": result,
		"count":     len(result),
	}

	summaryText := fmt.Sprintf("Found %d recent time logs", len(result))
	return formatMCPResponse(summaryText, response)
}

	// 新增：获取当前日期、时间、今天、昨天、本周日期范围
	type DateInfoParams struct{}

	func GetDateInfo(ctx context.Context, req *mcp.CallToolRequest, args DateInfoParams) (*mcp.CallToolResult, interface{}, error) {
		now := time.Now().In(singaporeLocation)
		today := now.Format("2006-01-02")
		yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
		weekday := now.Weekday()
		// 以周一为一周的开始
		daysSinceMonday := (int(weekday) + 6) % 7
		monday := now.AddDate(0, 0, -daysSinceMonday)
		sunday := monday.AddDate(0, 0, 6)
		weekRange := []string{
			monday.Format("2006-01-02"),
			sunday.Format("2006-01-02"),
		}

		response := map[string]interface{}{
			"timezone":   "Asia/Singapore (SGT, UTC+8)",
			"now":        now.Format("2006-01-02 15:04:05"),
			"today":      today,
			"yesterday":  yesterday,
			"weekday":    weekday.String(),
			"week_range": weekRange,
		}
		summaryText := "当前日期和时间信息，包括今天、昨天和本周日期范围"
		return formatMCPResponse(summaryText, response)
	}

func GetTimeLogsByDateRange(ctx context.Context, req *mcp.CallToolRequest, args DateRangeParams) (*mcp.CallToolResult, interface{}, error) {
	startDateStr := args.StartDate
	endDateStr := args.EndDate

	// 解析为新加坡时区的日期
	startDate, err := time.ParseInLocation("2006-01-02", startDateStr, singaporeLocation)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid start_date format, use YYYY-MM-DD: %w", err)
	}

	endDate, err := time.ParseInLocation("2006-01-02", endDateStr, singaporeLocation)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid end_date format, use YYYY-MM-DD: %w", err)
	}

	// Set time to cover full day in local timezone
	// startDate 已经是当天 00:00:00 SGT
	// endDate 需要设置为当天 23:59:59 SGT
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	timeLogs, err := model.ListTimeLogsWithOptions(server.db, 0, "start_time ASC", "start_time >= ? AND start_time <= ?", startDate, endDate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get time logs by date range: %w", err)
	}

	var result []map[string]interface{}
	totalDuration := time.Duration(0)

	for _, tl := range timeLogs {
		duration := time.Duration(0)
		durationStr := ""

		if tl.EndTime != nil {
			duration = tl.EndTime.Sub(tl.StartTime)
			totalDuration += duration
			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60
			durationStr = fmt.Sprintf("%dh %dm", hours, minutes)
		} else {
			durationStr = "ongoing"
		}

	entry := map[string]interface{}{
		"id":         tl.ID,
		"start_time": tl.StartTime.In(singaporeLocation).Format("2006-01-02 15:04:05"),
		"end_time":   nil,
		"duration":   durationStr,
		"tag":        tl.Tag.Name,
		"tag_color":  tl.Tag.Color,
		"remarks":    tl.Remark,
	}

	if tl.EndTime != nil {
		entry["end_time"] = tl.EndTime.In(singaporeLocation).Format("2006-01-02 15:04:05")
	}

	if tl.Task != nil {
		entry["task"] = map[string]interface{}{
			"id":    tl.Task.ID,
			"title": tl.Task.Title,
		}
	}

	result = append(result, entry)
}

	totalHours := int(totalDuration.Hours())
	totalMinutes := int(totalDuration.Minutes()) % 60

	response := map[string]interface{}{
		"time_logs":      result,
		"count":          len(result),
		"date_range":     fmt.Sprintf("%s to %s", startDateStr, endDateStr[:10]),
		"total_duration": fmt.Sprintf("%dh %dm", totalHours, totalMinutes),
	}

	summaryText := fmt.Sprintf("Found %d time logs from %s to %s, total duration: %dh %dm", len(result), startDateStr, endDateStr[:10], totalHours, totalMinutes)
	return formatMCPResponse(summaryText, response)
}

func GetTasksByStatus(ctx context.Context, req *mcp.CallToolRequest, args TaskStatusParams) (*mcp.CallToolResult, interface{}, error) {
	statusStr := args.Status

	tasks, err := model.GetAllTasks(server.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var result []map[string]interface{}
	for _, task := range tasks {
		// Filter by status
		if statusStr == "completed" && !task.IsCompleted {
			continue
		}
		if statusStr == "pending" && task.IsCompleted {
			continue
		}

		entry := map[string]interface{}{
			"id":                task.ID,
			"title":             task.Title,
			"description":       task.Description,
			"tag":               task.Tag.Name,
			"tag_color":         task.Tag.Color,
			"due_date":          task.DueDate.In(singaporeLocation).Format("2006-01-02"),
			"estimated_minutes": task.EstimatedMinutes,
			"is_completed":      task.IsCompleted,
			"created_at":        task.CreatedAt.In(singaporeLocation).Format("2006-01-02 15:04:05"),
		}

		if task.CompletedAt != nil {
			entry["completed_at"] = task.CompletedAt.In(singaporeLocation).Format("2006-01-02 15:04:05")
		}

		result = append(result, entry)
	}

	response := map[string]interface{}{
		"tasks":  result,
		"count":  len(result),
		"status": statusStr,
	}

	summaryText := fmt.Sprintf("Found %d %s tasks", len(result), statusStr)
	return formatMCPResponse(summaryText, response)
}

func GetProductivityStats(ctx context.Context, req *mcp.CallToolRequest, args StatsParams) (*mcp.CallToolResult, interface{}, error) {
	days := args.Days
	if days == 0 {
		days = 7 // default to last 7 days
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Get time logs in the date range
	timeLogs, err := model.ListTimeLogsWithOptions(server.db, 0, "start_time ASC", "start_time >= ?", startDate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get time logs for productivity stats: %w", err)
	}

	// Calculate daily stats
	dailyStats := make(map[string]map[string]interface{})
	tagStats := make(map[string]time.Duration)
	totalDuration := time.Duration(0)

	for _, tl := range timeLogs {
		if tl.EndTime == nil {
			continue // Skip ongoing logs
		}

		duration := tl.EndTime.Sub(tl.StartTime)
		totalDuration += duration

		date := tl.StartTime.In(singaporeLocation).Format("2006-01-02")
		if _, exists := dailyStats[date]; !exists {
			dailyStats[date] = map[string]interface{}{
				"date":     date,
				"duration": time.Duration(0),
				"entries":  0,
			}
		}

		dailyStats[date]["duration"] = dailyStats[date]["duration"].(time.Duration) + duration
		dailyStats[date]["entries"] = dailyStats[date]["entries"].(int) + 1

		// Tag statistics
		tagStats[tl.Tag.Name] += duration
	}

	// Convert daily stats to slice
	var dailyArray []map[string]interface{}
	for _, stats := range dailyStats {
		duration := stats["duration"].(time.Duration)
		hours := duration.Hours()
		dailyArray = append(dailyArray, map[string]interface{}{
			"date":            stats["date"],
			"duration_hours":  math.Round(hours*100) / 100,
			"duration_string": fmt.Sprintf("%.1fh", hours),
			"entries":         stats["entries"],
		})
	}

	// Convert tag stats
	var tagArray []map[string]interface{}
	for tagName, duration := range tagStats {
		hours := duration.Hours()
		percentage := (duration.Seconds() / totalDuration.Seconds()) * 100
		tagArray = append(tagArray, map[string]interface{}{
			"tag":             tagName,
			"duration_hours":  math.Round(hours*100) / 100,
			"duration_string": fmt.Sprintf("%.1fh", hours),
			"percentage":      math.Round(percentage*100) / 100,
		})
	}

	avgDaily := totalDuration.Hours() / float64(days)

	response := map[string]interface{}{
		"period":          fmt.Sprintf("Last %d days", days),
		"total_hours":     math.Round(totalDuration.Hours()*100) / 100,
		"average_daily":   math.Round(avgDaily*100) / 100,
		"daily_breakdown": dailyArray,
		"tag_breakdown":   tagArray,
		"total_entries":   len(timeLogs),
	}

	summaryText := fmt.Sprintf("Productivity stats for last %d days: %.1f total hours, %.1f average daily hours, %d entries", days, totalDuration.Hours(), avgDaily, len(timeLogs))
	return formatMCPResponse(summaryText, response)
}

func GetTaskCompletionAnalysis(ctx context.Context, req *mcp.CallToolRequest, args TaskAnalysisParams) (*mcp.CallToolResult, interface{}, error) {
	days := args.Days
	if days == 0 {
		days = 30 // default to last 30 days
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Get all tasks in the date range
	tasks, err := model.GetTasksByDateRange(server.db, startDate, endDate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tasks for completion analysis: %w", err)
	}

	totalTasks := len(tasks)
	completedTasks := 0
	overdueTasks := 0
	totalEstimatedMinutes := 0
	completedEstimatedMinutes := 0

	tagStats := make(map[string]map[string]int)

	for _, task := range tasks {
		totalEstimatedMinutes += task.EstimatedMinutes

		if task.IsCompleted {
			completedTasks++
			completedEstimatedMinutes += task.EstimatedMinutes
		}

		if !task.IsCompleted && task.DueDate.Before(time.Now()) {
			overdueTasks++
		}

		// Track by tag
		if _, exists := tagStats[task.Tag.Name]; !exists {
			tagStats[task.Tag.Name] = map[string]int{
				"total":     0,
				"completed": 0,
			}
		}
		tagStats[task.Tag.Name]["total"]++
		if task.IsCompleted {
			tagStats[task.Tag.Name]["completed"]++
		}
	}

	completionRate := 0.0
	if totalTasks > 0 {
		completionRate = (float64(completedTasks) / float64(totalTasks)) * 100
	}

	// Convert tag stats
	var tagArray []map[string]interface{}
	for tagName, stats := range tagStats {
		rate := 0.0
		if stats["total"] > 0 {
			rate = (float64(stats["completed"]) / float64(stats["total"])) * 100
		}
		tagArray = append(tagArray, map[string]interface{}{
			"tag":             tagName,
			"total":           stats["total"],
			"completed":       stats["completed"],
			"completion_rate": math.Round(rate*100) / 100,
		})
	}

	response := map[string]interface{}{
		"period":                    fmt.Sprintf("Last %d days", days),
		"total_tasks":               totalTasks,
		"completed_tasks":           completedTasks,
		"overdue_tasks":             overdueTasks,
		"completion_rate":           math.Round(completionRate*100) / 100,
		"total_estimated_hours":     float64(totalEstimatedMinutes) / 60,
		"completed_estimated_hours": float64(completedEstimatedMinutes) / 60,
		"tag_breakdown":             tagArray,
	}

	summaryText := fmt.Sprintf("Task completion analysis for last %d days: %d total tasks, %d completed (%.1f%% rate), %d overdue", days, totalTasks, completedTasks, completionRate, overdueTasks)
	return formatMCPResponse(summaryText, response)
}

func GetCurrentActivity(ctx context.Context, req *mcp.CallToolRequest, args CurrentActivityParams) (*mcp.CallToolResult, interface{}, error) {
	timeLogs, err := model.ListTimeLogsWithOptions(server.db, 0, "start_time DESC", "end_time IS NULL")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current activity: %w", err)
	}

	var result []map[string]interface{}
	for _, tl := range timeLogs {
		duration := time.Since(tl.StartTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60

		entry := map[string]interface{}{
			"id":         tl.ID,
			"start_time": tl.StartTime.In(singaporeLocation).Format("2006-01-02 15:04:05"),
			"duration":   fmt.Sprintf("%dh %dm", hours, minutes),
			"tag":        tl.Tag.Name,
			"tag_color":  tl.Tag.Color,
			"remarks":    tl.Remark,
		}

		if tl.Task != nil {
			entry["task"] = map[string]interface{}{
				"id":    tl.Task.ID,
				"title": tl.Task.Title,
			}
		}

		result = append(result, entry)
	}

	response := map[string]interface{}{
		"active_logs": result,
		"count":       len(result),
	}

	summaryText := fmt.Sprintf("Found %d active time logs", len(result))
	return formatMCPResponse(summaryText, response)
}

func GetActiveConstraints(ctx context.Context, req *mcp.CallToolRequest, args ConstraintParams) (*mcp.CallToolResult, interface{}, error) {
	constraints, err := model.GetActiveConstraints(server.db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get active constraints: %w", err)
	}

	var result []map[string]interface{}
	for _, constraint := range constraints {
		entry := map[string]interface{}{
			"id":               constraint.ID,
			"description":      constraint.Description,
			"punishment_quote": constraint.PunishmentQuote,
			"start_date":       constraint.StartDate.In(singaporeLocation).Format("2006-01-02"),
			"is_active":        constraint.IsActive,
			"created_at":       constraint.CreatedAt.In(singaporeLocation).Format("2006-01-02 15:04:05"),
		}

		if constraint.EndDate != nil {
			entry["end_date"] = constraint.EndDate.In(singaporeLocation).Format("2006-01-02")
		}
		if constraint.EndReason != "" {
			entry["end_reason"] = constraint.EndReason
		}

		result = append(result, entry)
	}

	response := map[string]interface{}{
		"constraints": result,
		"count":       len(result),
	}

	summaryText := fmt.Sprintf("Found %d active constraints", len(result))
	return formatMCPResponse(summaryText, response)
}
