package model

import (
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
	"gorm.io/gorm"
)

// --- CRUD 操作 ---

// CreateConstraint 创建约束
func CreateConstraint(db *gorm.DB, constraint *gen.Constraint) error {
	return db.Create(constraint).Error
}

// GetConstraintByID 根据ID获取约束
func GetConstraintByID(db *gorm.DB, id int32) (*gen.Constraint, error) {
	var constraint gen.Constraint
	err := db.First(&constraint, id).Error
	if err != nil {
		return nil, err
	}
	return &constraint, nil
}

// GetAllConstraints 获取所有约束
func GetAllConstraints(db *gorm.DB) ([]gen.Constraint, error) {
	var constraints []gen.Constraint
	err := db.Find(&constraints).Error
	return constraints, err
}

// GetActiveConstraints 获取活跃的约束
func GetActiveConstraints(db *gorm.DB) ([]gen.Constraint, error) {
	var constraints []gen.Constraint
	err := db.Where("is_active = ?", true).Find(&constraints).Error
	return constraints, err
}

// GetConstraintsByDateRange 根据日期范围获取约束
func GetConstraintsByDateRange(db *gorm.DB, startDate, endDate time.Time) ([]gen.Constraint, error) {
	var constraints []gen.Constraint
	err := db.Where("start_date >= ? AND start_date <= ?", startDate, endDate).Find(&constraints).Error
	return constraints, err
}

// UpdateConstraint 更新约束
func UpdateConstraint(db *gorm.DB, constraint *gen.Constraint) error {
	return db.Save(constraint).Error
}

// DeleteConstraint 删除约束 (软删除)
func DeleteConstraint(db *gorm.DB, id int32) error {
	return db.Delete(&gen.Constraint{}, id).Error
}

// MarkConstraintAsCompleted 标记约束为完成
func MarkConstraintAsCompleted(db *gorm.DB, constraintID int32, endReason string) error {
	now := time.Now()
	falseValue := false
	return db.Model(&gen.Constraint{}).Where("id = ?", constraintID).Updates(map[string]interface{}{
		"is_active":  &falseValue,
		"end_date":   &now,
		"end_reason": &endReason,
	}).Error
}

// MarkConstraintAsActive 重新激活约束
func MarkConstraintAsActive(db *gorm.DB, constraintID int32) error {
	trueValue := true
	emptyReason := ""
	return db.Model(&gen.Constraint{}).Where("id = ?", constraintID).Updates(map[string]interface{}{
		"is_active":  &trueValue,
		"end_date":   nil,
		"end_reason": &emptyReason,
	}).Error
}
