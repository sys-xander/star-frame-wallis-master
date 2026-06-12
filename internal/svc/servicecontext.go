package svc

import (
	"wallpaper/internal/config"
	"wallpaper/internal/pkg/mp"
	"wallpaper/internal/pkg/qiniu"

	"gorm.io/driver/mysql"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	Orm         *gorm.DB
	Cache       *ristretto.Cache[string, any]
	MiniProgram *mp.MiniProgramService
	QiniuYun    *qiniu.QiniuYun
}

func NewServiceContext(c config.Config) *ServiceContext {

	db, err := gorm.Open(mysql.Open(c.Gorm.Dsn()), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		logx.Errorw("[NewServiceContext] failed to connect database", logx.Field("error", err), logx.Field("dsn", c.Gorm.Dsn()))
		panic(err)
	}

	cache, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	logx.Must(err)

	return &ServiceContext{
		Config:      c,
		Orm:         db,
		Cache:       cache,
		MiniProgram: mp.NewMiniProgram(&c),
		QiniuYun:    qiniu.NewQiniuYun(c.Qiniu),
	}
}
