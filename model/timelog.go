package model

import (
	"errors"
	"time"

	"github.com/blacksheepaul/timelog/core/logger"
	"gorm.io/gorm"
)

// TimeLog 表示一条时间日志
// 字段与 timelogs 表结构保持一致
// id, user_id, start_time, end_time, tag, remark

type TimeLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"column:user_id" json:"user_id"`
	StartTime time.Time      `gorm:"column:start_time" json:"start_time"`
	EndTime   *time.Time     `gorm:"column:end_time" json:"end_time"`
	TagID     uint           `gorm:"column:tag_id" json:"tag_id"`
	Tag       Tag            `gorm:"foreignKey:TagID" json:"tag"`
	TaskID    *uint          `gorm:"column:task_id" json:"task_id,omitempty"`
	Task      *Task          `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Remark    string         `gorm:"column:remark" json:"remarks"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TimeLog) TableName() string {
	return "timelogs"
}

// BeforeCreate 在创建记录前执行，确保created_at不是零值并验证user_id
func (tl *TimeLog) BeforeCreate(tx *gorm.DB) error {
	// 验证UserID不能为0
	if tl.UserID == 0 {
		log := logger.GetLogger()
		log.Errorw("TimeLog creation failed: UserID cannot be 0",
			"user_id", tl.UserID,
			"start_time", tl.StartTime,
			"tag_id", tl.TagID)
		return errors.New("user_id cannot be 0")
	}

	// Bug tracking: 确保CreatedAt不是零值
	if tl.CreatedAt.IsZero() {
		log := logger.GetLogger()
		log.Errorw("TimeLog creation failed: CreatedAt cannot be zero",
			"user_id", tl.UserID,
			"start_time", tl.StartTime,
			"tag_id", tl.TagID)
		return errors.New("created_at cannot be zero")
	}

	return nil
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
