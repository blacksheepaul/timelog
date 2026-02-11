package model

import (
	"time"

	"gorm.io/gorm"
)

type TempPassword struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PasswordHash string         `gorm:"column:password_hash;not null" json:"password_hash"`
	ExpiresAt    time.Time      `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TempPassword) TableName() string {
	return "temp_passwords"
}

func CreateTempPassword(db *gorm.DB, tempPassword *TempPassword) error {
	return db.Create(tempPassword).Error
}

func ListTempPasswords(db *gorm.DB) ([]TempPassword, error) {
	var passwords []TempPassword
	err := db.Order("created_at DESC").Find(&passwords).Error
	return passwords, err
}

func DeleteTempPassword(db *gorm.DB, id uint) error {
	return db.Delete(&TempPassword{}, id).Error
}

func DeleteExpiredTempPasswords(db *gorm.DB, now time.Time) error {
	return db.Where("expires_at <= ?", now).Delete(&TempPassword{}).Error
}

func GetTempPasswordByHash(db *gorm.DB, hash string, now time.Time) (*TempPassword, error) {
	var password TempPassword
	err := db.Where("password_hash = ? AND expires_at > ?", hash, now).First(&password).Error
	return &password, err
}
