package model

import (
	"time"

	"gorm.io/gorm"
)

// TimeLog 表示一条时间日志
// 字段与 timelogs 表结构保持一致
// id, user_id, start_time, end_time, tag, remark

type TimeLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"column:user_id" json:"user_id"`
	StartTime time.Time      `gorm:"column:start_time" json:"start_time"`
	EndTime   time.Time      `gorm:"column:end_time" json:"end_time"`
	Tag       string         `gorm:"column:tag" json:"tag"`
	Remark    string         `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
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
	err := db.First(&tl, id).Error
	return &tl, err
}

// ListTimeLogs 查询时间日志（可扩展条件）
func ListTimeLogs(db *gorm.DB, conds ...interface{}) ([]TimeLog, error) {
	var tls []TimeLog
	err := db.Find(&tls, conds...).Error
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
