package model

import (
	"time"

	"gorm.io/gorm"
)

// Task 任务模型
// 字段设计：
// - id: 主键
// - title: 任务标题 (必填)
// - description: 任务详情描述 (可选)
// - tag_id: 关联的标签ID (必填)
// - tag: 标签对象 (外键关联)
// - due_date: 任务预定日期 (必填)
// - estimated_minutes: 预估完成时间(分钟) (必填)
// - is_completed: 是否完成 (默认false)
// - completed_at: 完成时间 (完成时自动设置)
// - is_suspended: 是否暂停 (默认false)
// - created_at, updated_at, deleted_at: ORM基础字段

type Task struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Title            string         `gorm:"column:title;not null;size:200" json:"title"`
	Description      string         `gorm:"column:description;type:text" json:"description"`
	TagID            uint           `gorm:"column:tag_id;not null" json:"tag_id"`
	Tag              Tag            `gorm:"foreignKey:TagID" json:"tag"`
	DueDate          time.Time      `gorm:"column:due_date;not null" json:"due_date"`
	EstimatedMinutes int            `gorm:"column:estimated_minutes;not null" json:"estimated_minutes"`
	IsCompleted      bool           `gorm:"column:is_completed;default:false" json:"is_completed"`
	CompletedAt      *time.Time     `gorm:"column:completed_at" json:"completed_at"`
	IsSuspended      bool           `gorm:"column:is_suspended;default:false" json:"is_suspended"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Task) TableName() string {
	return "tasks"
}

// --- CRUD 操作 ---

// CreateTask 创建任务
func CreateTask(db *gorm.DB, task *Task) error {
	return db.Create(task).Error
}

// GetTaskByID 根据ID获取任务
func GetTaskByID(db *gorm.DB, id uint) (*Task, error) {
	var task Task
	err := db.Preload("Tag").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAllTasks 获取所有任务
// includeSuspended: 是否包含暂停的任务
// includeCompleted: 是否包含已完成的任务
func GetAllTasks(db *gorm.DB, includeSuspended bool, includeCompleted bool) ([]Task, error) {
	var tasks []Task
	query := db.Preload("Tag")

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
func GetTasksByDate(db *gorm.DB, date time.Time, includeSuspended bool, includeCompleted bool) ([]Task, error) {
	var tasks []Task
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := db.Preload("Tag").Where("due_date >= ? AND due_date < ?", startOfDay, endOfDay)

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
func GetTasksByDateRange(db *gorm.DB, startDate, endDate time.Time) ([]Task, error) {
	var tasks []Task
	err := db.Preload("Tag").Where("due_date >= ? AND due_date <= ?", startDate, endDate).Find(&tasks).Error
	return tasks, err
}

// UpdateTask 更新任务
func UpdateTask(db *gorm.DB, task *Task) error {
	return db.Save(task).Error
}

// DeleteTask 删除任务 (软删除)
func DeleteTask(db *gorm.DB, id uint) error {
	return db.Delete(&Task{}, id).Error
}

// MarkTaskAsCompleted 标记任务为完成
func MarkTaskAsCompleted(db *gorm.DB, taskID uint) error {
	now := time.Now()
	return db.Model(&Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"is_completed": true,
		"completed_at": &now,
	}).Error
}

// MarkTaskAsIncomplete 标记任务为未完成
func MarkTaskAsIncomplete(db *gorm.DB, taskID uint) error {
	return db.Model(&Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"is_completed": false,
		"completed_at": nil,
	}).Error
}

// SuspendTask 暂停任务
func SuspendTask(db *gorm.DB, taskID uint) error {
	return db.Model(&Task{}).Where("id = ?", taskID).Update("is_suspended", true).Error
}

// UnsuspendTask 取消暂停任务
func UnsuspendTask(db *gorm.DB, taskID uint) error {
	return db.Model(&Task{}).Where("id = ?", taskID).Update("is_suspended", false).Error
}

// GetCompletedTasksInDateRange 获取指定日期范围内的已完成任务
func GetCompletedTasksInDateRange(db *gorm.DB, startDate, endDate time.Time) ([]Task, error) {
	var tasks []Task
	err := db.Preload("Tag").Where("is_completed = ? AND completed_at >= ? AND completed_at <= ?",
		true, startDate, endDate).Find(&tasks).Error
	return tasks, err
}

// GetTaskStats 获取任务统计信息
func GetTaskStats(db *gorm.DB, date time.Time) (map[string]interface{}, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var totalTasks, completedTasks int64

	// 获取当天的总任务数
	err := db.Model(&Task{}).Where("due_date >= ? AND due_date < ?", startOfDay, endOfDay).Count(&totalTasks).Error
	if err != nil {
		return nil, err
	}

	// 获取当天已完成的任务数
	err = db.Model(&Task{}).Where("due_date >= ? AND due_date < ? AND is_completed = ?",
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
