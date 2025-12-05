package model

import (
	"time"

	"gorm.io/gorm"
)

// Model 基础结构模型
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"deleted_at;index" json:"-"`
}
