package model

import (
	"time"
)

type Classification struct {
	ID uint `gorm:"primarykey"`
	Name  string `gorm:"column:name;NOT NULL"`
	Url   string `gorm:"column:url;NOT NULL"`
	Count uint   `gorm:"column:count;type:int unsigned;default:0" json:"count"`
	Icon  string `gorm:"column:icon;type:varchar(255);NOT NULL"`
	Gradient string `gorm:"column:gradient;type:varchar(255);NOT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Classification) TableName() string {
	return "classification"
}
