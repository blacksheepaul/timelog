package model

import (
	"time"

	"gorm.io/gorm"
)

// TimeLog 表示一条时间日志
// 字段与 timelogs 表结构保持一致
// id, user_id, start_time, end_time, tag, remark

type TimeLog struct {
	gorm.Model
	UserID    uint       `gorm:"column:user_id" json:"user_id"`
	StartTime time.Time  `gorm:"column:start_time" json:"start_time"`
	EndTime   *time.Time `gorm:"column:end_time" json:"end_time"`
	TagID     uint       `gorm:"column:tag_id" json:"tag_id"`
	Tag       Tag        `gorm:"foreignKey:TagID" json:"tag"`
	TaskID    *uint      `gorm:"column:task_id" json:"task_id,omitempty"`
	Task      *Task      `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Remark    string     `gorm:"column:remark" json:"remarks"`
}

func (TimeLog) TableName() string {
	return "timelogs"
}

// --- CRUD ---

// CreateTimeLog 新增一条时间日志
func CreateTimeLog(db *gorm.DB, tl *TimeLog) error {
	return db.Create(tl).Error
}

// GetTimeLogByID 根据ID获取时间日志
func GetTimeLogByID(db *gorm.DB, id uint) (*TimeLog, error) {
	var tl TimeLog
	err := db.Preload("Tag").Preload("Task").First(&tl, id).Error
	return &tl, err
}

// ListTimeLogs 查询时间日志（可扩展条件）
func ListTimeLogs(db *gorm.DB, conds ...interface{}) ([]TimeLog, error) {
	var tls []TimeLog
	err := db.Preload("Tag").Preload("Task").Find(&tls, conds...).Error
	return tls, err
}

// ListTimeLogsWithOptions 查询时间日志（支持排序和限制）
func ListTimeLogsWithOptions(db *gorm.DB, limit int, orderBy string, conds ...interface{}) ([]TimeLog, error) {
	var tls []TimeLog
	query := db.Preload("Tag").Preload("Task")

	if len(conds) > 0 {
		query = query.Where(conds[0], conds[1:]...)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&tls).Error
	return tls, err
}

// UpdateTimeLog 更新时间日志
func UpdateTimeLog(db *gorm.DB, tl *TimeLog) error {
	return db.Save(tl).Error
}

// DeleteTimeLog 删除时间日志
func DeleteTimeLog(db *gorm.DB, id uint) error {
	return db.Delete(&TimeLog{}, id).Error
}
