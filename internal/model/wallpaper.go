package model

import (
	"gorm.io/gorm"
)

type Wallpaper struct {
	gorm.Model

	Title         string `gorm:"column:title;type:varchar(255);NOT NULL" json:"title"`
	Url           string `gorm:"column:url;type:varchar(500);NOT NULL" json:"url"`
	Md5           string `gorm:"column:md5;type:varchar(255);NOT NULL" json:"md5"`
	LikeCount     uint   `gorm:"column:like_count;type:int unsigned;default:0" json:"like_count"`
	DownloadCount uint   `gorm:"column:download_count;type:int unsigned;default:0" json:"download_count"`
	CollectCount  uint   `gorm:"column:collect_count;type:int unsigned;default:0" json:"collect_count"`
	ViewCount     uint   `gorm:"column:view_count;type:int unsigned;default:0" json:"view_count"`
	ClassifyId    int    `gorm:"column:classify_id;type:int(11);NOT NULL" json:"classify_id"`
}

func (m *Wallpaper) TableName() string {
	return "wallpaper"
}
