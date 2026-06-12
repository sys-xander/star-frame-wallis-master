package model

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID        uint           `gorm:"primarykey"`
	Title     string         `gorm:"column:title;type:varchar(255);NOT NULL" json:"title"`
	Content   string         `gorm:"column:content;type:text;NOT NULL" json:"content"`
	Type      uint8          `gorm:"column:type;type:tinyint unsigned;default:1" json:"type"`     // 1系统 2活动 3更新
	Status    uint8          `gorm:"column:status;type:tinyint unsigned;default:1" json:"status"` // 1正常 0禁用
	Sort      int            `gorm:"column:sort;type:int;default:0" json:"sort"`
	Cover     string         `gorm:"column:cover;type:varchar(500);default:''" json:"cover"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (Announcement) TableName() string {
	return "announcement"
}
