package mp

import (
	"wallpaper/internal/config"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	"github.com/zeromicro/go-zero/core/logx"
)

type MiniProgramService struct {
	*miniProgram.MiniProgram
}

func NewMiniProgram(c *config.Config) *MiniProgramService {
	app, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:     c.Wx.AppID,
		Secret:    c.Wx.AppSecret,
		HttpDebug: true,
		Debug:     false,
	})

	logx.Must(err)

	return &MiniProgramService{
		MiniProgram: app,
	}
}