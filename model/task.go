package model

import (
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
	"gorm.io/gorm"
)

// --- CRUD 操作 ---

// CreateTask 创建任务
func CreateTask(db *gorm.DB, task *gen.Task) error {
	return db.Create(task).Error
}

// GetTaskByID 根据ID获取任务
func GetTaskByID(db *gorm.DB, id int32) (*gen.Task, error) {
	var task gen.Task
	err := db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAllTasks 获取所有任务
// includeSuspended: 是否包含暂停的任务
// includeCompleted: 是否包含已完成的任务
func GetAllTasks(db *gorm.DB, includeSuspended bool, includeCompleted bool) ([]gen.Task, error) {
	var tasks []gen.Task
	query := db

	if !includeSuspended {
		query = query.Where("is_suspended = ?", false)
	}

	if !includeCompleted {
		query = query.Where("is_completed = ?", false)
	}

	err := query.Find(&tasks).Error
	return tasks, err
}

// GetTasksByDate 根据日期获取任务
// includeSuspended: 是否包含暂停的任务
// includeCompleted: 是否包含已完成的任务
func GetTasksByDate(db *gorm.DB, date time.Time, includeSuspended bool, includeCompleted bool) ([]gen.Task, error) {
	var tasks []gen.Task
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := db.Where("due_date >= ? AND due_date < ?", startOfDay, endOfDay)

	if !includeSuspended {
		query = query.Where("is_suspended = ?", false)
	}

	if !includeCompleted {
		query = query.Where("is_completed = ?", false)
	}

	err := query.Find(&tasks).Error
	return tasks, err
}

// GetTasksByDateRange 根据日期范围获取任务
func GetTasksByDateRange(db *gorm.DB, startDate, endDate time.Time) ([]gen.Task, error) {
	var tasks []gen.Task
	err := db.Where("due_date >= ? AND due_date <= ?", startDate, endDate).Find(&tasks).Error
	return tasks, err
}

// UpdateTask 更新任务
func UpdateTask(db *gorm.DB, task *gen.Task) error {
	return db.Save(task).Error
}

// DeleteTask 删除任务 (软删除)
func DeleteTask(db *gorm.DB, id int32) error {
	return db.Delete(&gen.Task{}, id).Error
}

// MarkTaskAsCompleted 标记任务为完成
func MarkTaskAsCompleted(db *gorm.DB, taskID int32) error {
	now := time.Now()
	trueValue := true
	return db.Model(&gen.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"is_completed": &trueValue,
		"completed_at": &now,
	}).Error
}

// MarkTaskAsIncomplete 标记任务为未完成
func MarkTaskAsIncomplete(db *gorm.DB, taskID int32) error {
	falseValue := false
	return db.Model(&gen.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"is_completed": &falseValue,
		"completed_at": nil,
	}).Error
}

// SuspendTask 暂停任务
func SuspendTask(db *gorm.DB, taskID int32) error {
	trueValue := true
	return db.Model(&gen.Task{}).Where("id = ?", taskID).Update("is_suspended", &trueValue).Error
}

// UnsuspendTask 取消暂停任务
func UnsuspendTask(db *gorm.DB, taskID int32) error {
	falseValue := false
	return db.Model(&gen.Task{}).Where("id = ?", taskID).Update("is_suspended", &falseValue).Error
}

// GetCompletedTasksInDateRange 获取指定日期范围内的已完成任务
func GetCompletedTasksInDateRange(db *gorm.DB, startDate, endDate time.Time) ([]gen.Task, error) {
	var tasks []gen.Task
	err := db.Where("is_completed = ? AND completed_at >= ? AND completed_at <= ?",
		true, startDate, endDate).Find(&tasks).Error
	return tasks, err
}

// GetTaskStats 获取任务统计信息
func GetTaskStats(db *gorm.DB, date time.Time) (map[string]interface{}, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var totalTasks, completedTasks int64

	// 获取当天的总任务数
	err := db.Model(&gen.Task{}).Where("due_date >= ? AND due_date < ?", startOfDay, endOfDay).Count(&totalTasks).Error
	if err != nil {
		return nil, err
	}

	// 获取当天已完成的任务数
	err = db.Model(&gen.Task{}).Where("due_date >= ? AND due_date < ? AND is_completed = ?",
		startOfDay, endOfDay, true).Count(&completedTasks).Error
	if err != nil {
		return nil, err
	}

	// 计算完成率
	completionRate := float64(0)
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}

	return map[string]interface{}{
		"total_tasks":     totalTasks,
		"completed_tasks": completedTasks,
		"completion_rate": completionRate,
	}, nil
}
