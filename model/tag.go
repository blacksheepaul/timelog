package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag 表示一个标签
type Tag struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"column:name;not null;unique" json:"name"`
	Color       string         `gorm:"column:color" json:"color"`
	Description string         `gorm:"column:description" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Tag) TableName() string {
	return "tags"
}

// --- CRUD ---

// CreateTag 新增一个标签
func CreateTag(db *gorm.DB, tag *Tag) error {
	return db.Create(tag).Error
}

// GetTagByID 根据ID获取标签
func GetTagByID(db *gorm.DB, id uint) (*Tag, error) {
	var tag Tag
	err := db.First(&tag, id).Error
	return &tag, err
}

// GetTagByName 根据名称获取标签
func GetTagByName(db *gorm.DB, name string) (*Tag, error) {
	var tag Tag
	err := db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

// ListTags 查询标签列表
func ListTags(db *gorm.DB, conds ...interface{}) ([]Tag, error) {
	var tags []Tag
	err := db.Find(&tags, conds...).Error
	return tags, err
}

// UpdateTag 更新标签
func UpdateTag(db *gorm.DB, tag *Tag) error {
	return db.Save(tag).Error
}

// DeleteTag 删除标签
func DeleteTag(db *gorm.DB, id uint) error {
	return db.Delete(&Tag{}, id).Error
}
