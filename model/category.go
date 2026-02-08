package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Category 表示一个分类（支持3层深度：0=根, 1=子, 2=孙）
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"column:name;not null" json:"name"`
	Color       string         `gorm:"column:color" json:"color"`
	Description string         `gorm:"column:description" json:"description"`
	ParentID    *uint          `gorm:"column:parent_id" json:"parent_id,omitempty"`
	Level       int            `gorm:"column:level;default:0" json:"level"`
	SortOrder   int            `gorm:"column:sort_order;default:0" json:"sort_order"`
	Path        string         `gorm:"column:path;default:'/'" json:"path"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}

const MaxCategoryLevel = 2 // 最大层级：0, 1, 2

// ValidateLevel 验证分类层级是否合法
func (c *Category) ValidateLevel() error {
	if c.Level < 0 || c.Level > MaxCategoryLevel {
		return fmt.Errorf("category level must be between 0 and %d", MaxCategoryLevel)
	}
	return nil
}

// GetFullPath 获取完整路径（包含自身）
func (c *Category) GetFullPath() string {
	if c.Path == "/" {
		return fmt.Sprintf("/%s", c.Name)
	}
	return fmt.Sprintf("%s/%s", c.Path, c.Name)
}

// --- CRUD ---

// CreateCategory 创建分类（自动计算level和path）
func CreateCategory(db *gorm.DB, category *Category) error {
	// 如果有父分类，计算level和path
	if category.ParentID != nil && *category.ParentID > 0 {
		parent, err := GetCategoryByID(db, *category.ParentID)
		if err != nil {
			return fmt.Errorf("parent category not found: %w", err)
		}
		category.Level = parent.Level + 1
		if category.Level > MaxCategoryLevel {
			return fmt.Errorf("cannot create category: exceeds max level %d", MaxCategoryLevel)
		}
		category.Path = parent.GetFullPath()
	} else {
		category.Level = 0
		category.Path = "/"
		category.ParentID = nil
	}

	return db.Create(category).Error
}

// GetCategoryByID 根据ID获取分类
func GetCategoryByID(db *gorm.DB, id uint) (*Category, error) {
	var category Category
	err := db.First(&category, id).Error
	return &category, err
}

// GetCategoryByName 根据名称和父ID获取分类（用于检查重复）
func GetCategoryByName(db *gorm.DB, name string, parentID *uint) (*Category, error) {
	var category Category
	query := db.Where("name = ?", name)
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	err := query.First(&category).Error
	return &category, err
}

// ListCategories 查询分类列表（支持筛选条件）
func ListCategories(db *gorm.DB, conds ...interface{}) ([]Category, error) {
	var categories []Category
	err := db.Order("level ASC, sort_order DESC, name ASC").Find(&categories, conds...).Error
	return categories, err
}

// ListCategoriesByLevel 按层级查询分类
func ListCategoriesByLevel(db *gorm.DB, level int) ([]Category, error) {
	var categories []Category
	err := db.Where("level = ?", level).Order("sort_order DESC, name ASC").Find(&categories).Error
	return categories, err
}

// GetCategoriesByParentID 获取指定父分类下的子分类
func GetCategoriesByParentID(db *gorm.DB, parentID *uint) ([]Category, error) {
	var categories []Category
	query := db
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Order("sort_order DESC, name ASC").Find(&categories).Error
	return categories, err
}

// UpdateCategory 更新分类（不允许更改层级结构）
func UpdateCategory(db *gorm.DB, category *Category) error {
	// 获取原分类信息，防止更改层级
	existing, err := GetCategoryByID(db, category.ID)
	if err != nil {
		return err
	}

	// 保留层级相关信息，只允许更新基本属性
	category.Level = existing.Level
	category.ParentID = existing.ParentID
	category.Path = existing.Path

	return db.Save(category).Error
}


// MoveCategory 移动分类到新的父分类下
func MoveCategory(db *gorm.DB, categoryID uint, newParentID *uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		category, err := GetCategoryByID(tx, categoryID)
		if err != nil {
			return err
		}

		// 计算新的level
		newLevel := 0
		newPath := "/"

		if newParentID != nil && *newParentID > 0 {
			// 检查是否尝试移动到自己下面
			if *newParentID == categoryID {
				return fmt.Errorf("cannot move category to itself")
			}

			// 检查新父分类是否是当前分类的子分类（防止循环）
			var childCount int64
			if err := tx.Model(&Category{}).Where("id = ? AND (parent_id = ? OR path LIKE ?)", *newParentID, categoryID, fmt.Sprintf("%%%d%%", categoryID)).Count(&childCount).Error; err != nil {
				return err
			}
			if childCount > 0 {
				return fmt.Errorf("cannot move category to its own child")
			}

			parent, err := GetCategoryByID(tx, *newParentID)
			if err != nil {
				return fmt.Errorf("parent category not found: %w", err)
			}

			newLevel = parent.Level + 1
			if newLevel > MaxCategoryLevel {
				return fmt.Errorf("cannot move category: would exceed max level %d", MaxCategoryLevel)
			}
			newPath = parent.GetFullPath()
		}

		oldPath := category.GetFullPath()

		// 更新当前分类
		category.ParentID = newParentID
		category.Level = newLevel
		category.Path = newPath

		if err := tx.Save(category).Error; err != nil {
			return err
		}

		// 更新所有子分类的path
		newFullPath := category.GetFullPath()
		return tx.Model(&Category{}).Where("path LIKE ?", oldPath+"%").Update("path", gorm.Expr("REPLACE(path, ?, ?)", oldPath, newFullPath)).Error
	})
}

// GetCategoryTree 获取分类树
func GetCategoryTree(db *gorm.DB) ([]CategoryNode, error) {
	categories, err := ListCategories(db)
	if err != nil {
		return nil, err
	}

	return buildCategoryTree(categories), nil
}

// CategoryNode 分类树节点
type CategoryNode struct {
	Category Category       `json:"category"`
	Children []CategoryNode `json:"children,omitempty"`
}

// buildCategoryTree 构建分类树
func buildCategoryTree(categories []Category) []CategoryNode {
	nodeMap := make(map[uint]*CategoryNode)
	var roots []CategoryNode

	// 第一遍：创建所有节点
	for i := range categories {
		node := &CategoryNode{
			Category: categories[i],
			Children: []CategoryNode{},
		}
		nodeMap[categories[i].ID] = node
	}

	// 第二遍：构建父子关系
	for i := range categories {
		node := nodeMap[categories[i].ID]
		if categories[i].ParentID == nil || *categories[i].ParentID == 0 {
			roots = append(roots, *node)
		} else {
			if parent, ok := nodeMap[*categories[i].ParentID]; ok {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	return roots
}
