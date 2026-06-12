package model

import "time"

type UserDownload struct {
	ID          uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`
	UserID      uint64    `gorm:"column:user_id;type:bigint unsigned;NOT NULL;index" json:"user_id"`
	WallpaperID uint64    `gorm:"column:wallpaper_id;type:bigint unsigned;NOT NULL;index" json:"wallpaper_id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

func (UserDownload) TableName() string {
	return "user_download"
}
