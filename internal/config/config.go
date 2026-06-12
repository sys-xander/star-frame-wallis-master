package config

import (
	"fmt"

	"wallpaper/internal/pkg/sls"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth  Auth
	Wx    Wx
	Gorm  Gorm
	Qiniu Qiniu
	Sls   sls.Config
}

type Auth struct {
	AccessSecret    string
	AccessExpire    int64
	RefreshSecret   string
	RefreshExpire   int64
}

type Wx struct {
	AppID     string
	AppSecret string
}

type Gorm struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
	Params   string
}

func (g *Gorm) Dsn() string {
	return g.User + ":" + g.Password + "@tcp(" + g.Host + ":" + fmt.Sprint(g.Port) + ")/" + g.Database + "?" + g.Params
}

type Token struct {
	Expire int
	Secret string
	Iss    string
}

type Qiniu struct {
	Host      string
	AccessKey string
	SecretKey string
	Bucket    string
	Style     string
	Expire    int
}
