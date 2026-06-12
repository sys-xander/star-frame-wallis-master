package system

import (
	"wallpaper/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// init 运行在读取配置，启动项目之前
func init() {
	LoadEnv()
}

// Start 运行在读取配置文件，初始化服务之后, 启动项目之前
func Start(server *rest.Server, ctx *svc.ServiceContext) {
	StartSwagger(server, ctx)
	logx.Must(InitClassification(ctx))
}
