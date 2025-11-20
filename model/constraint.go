package model

import (
	"time"

	"gorm.io/gorm"
)

// Constraint 约束模型
// 字段设计：
// - id: 主键
// - description: 约束描述 (必填)
// - end_reason: 结束理由 (可选)
// - punishment_quote: 惩罚语录 (必填)
// - start_date: 开始日期 (必填)
// - end_date: 结束日期 (可选)
// - is_active: 是否激活 (默认true)
// - created_at, updated_at, deleted_at: ORM基础字段

type Constraint struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Description     string         `gorm:"column:description;not null;type:text" json:"description"`
	EndReason       string         `gorm:"column:end_reason;type:text" json:"end_reason"`
	PunishmentQuote string         `gorm:"column:punishment_quote;not null;type:text" json:"punishment_quote"`
	StartDate       time.Time      `gorm:"column:start_date;not null" json:"start_date"`
	EndDate         *time.Time     `gorm:"column:end_date" json:"end_date"`
	IsActive        bool           `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Constraint) TableName() string {
	return "constraints"
}

// --- CRUD 操作 ---

// CreateConstraint 创建约束
func CreateConstraint(db *gorm.DB, constraint *Constraint) error {
	return db.Create(constraint).Error
}

// GetConstraintByID 根据ID获取约束
func GetConstraintByID(db *gorm.DB, id uint) (*Constraint, error) {
	var constraint Constraint
	err := db.First(&constraint, id).Error
	if err != nil {
		return nil, err
	}
	return &constraint, nil
}

// GetAllConstraints 获取所有约束
func GetAllConstraints(db *gorm.DB) ([]Constraint, error) {
	var constraints []Constraint
	err := db.Find(&constraints).Error
	return constraints, err
}

// GetActiveConstraints 获取活跃的约束
func GetActiveConstraints(db *gorm.DB) ([]Constraint, error) {
	var constraints []Constraint
	err := db.Where("is_active = ?", true).Find(&constraints).Error
	return constraints, err
}

// GetConstraintsByDateRange 根据日期范围获取约束
func GetConstraintsByDateRange(db *gorm.DB, startDate, endDate time.Time) ([]Constraint, error) {
	var constraints []Constraint
	err := db.Where("start_date >= ? AND start_date <= ?", startDate, endDate).Find(&constraints).Error
	return constraints, err
}

// UpdateConstraint 更新约束
func UpdateConstraint(db *gorm.DB, constraint *Constraint) error {
	return db.Save(constraint).Error
}

// DeleteConstraint 删除约束 (软删除)
func DeleteConstraint(db *gorm.DB, id uint) error {
	return db.Delete(&Constraint{}, id).Error
}

// MarkConstraintAsCompleted 标记约束为完成
func MarkConstraintAsCompleted(db *gorm.DB, constraintID uint, endReason string) error {
	now := time.Now()
	return db.Model(&Constraint{}).Where("id = ?", constraintID).Updates(map[string]interface{}{
		"is_active":  false,
		"end_date":   &now,
		"end_reason": endReason,
	}).Error
}

// MarkConstraintAsActive 重新激活约束
func MarkConstraintAsActive(db *gorm.DB, constraintID uint) error {
	return db.Model(&Constraint{}).Where("id = ?", constraintID).Updates(map[string]interface{}{
		"is_active":  true,
		"end_date":   nil,
		"end_reason": "",
	}).Error
}
