package model

import (
	"gorm.io/gorm"
)

type MpUser struct {
	gorm.Model

	Name       string `gorm:"column:name;NOT NULL"`
	Openid     string `gorm:"column:openid;NOT NULL"`
	Unionid    string `gorm:"column:unionid;NOT NULL"`
	SessionKey string `gorm:"column:session_key;NOT NULL"`
	AvatarUrl  string `gorm:"column:avatar_url"` // 用户头像
	Phone      string `gorm:"column:phone;NOT NULL"`
}

func (MpUser) TableName() string {
	return "mp_user"
}
