package model

import "time"

type UserCollect struct {
	ID          uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement" json:"id"`
	UserID      uint64    `gorm:"column:user_id;type:bigint unsigned;NOT NULL;index:idx_user_wallpaper,unique" json:"user_id"`
	WallpaperID uint64    `gorm:"column:wallpaper_id;type:bigint unsigned;NOT NULL;index:idx_user_wallpaper,unique" json:"wallpaper_id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (UserCollect) TableName() string {
	return "user_collect"
}
