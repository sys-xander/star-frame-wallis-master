package main

import (
	"flag"
	"fmt"

	"wallpaper/internal/config"
	"wallpaper/internal/handler"
	"wallpaper/internal/middleware/log"
	"wallpaper/internal/pkg/sls"
	"wallpaper/internal/pkg/system"
	"wallpaper/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/wallpaper.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	// 如果启用 SLS，替换 logx 的 Writer
	if c.Sls.Enabled {
		slsWriter, err := sls.NewWriter(c.Sls)
		if err != nil {
			logx.Errorw("init sls writer failed", logx.Field("error", err))
		} else if slsWriter != nil {
			logx.SetWriter(slsWriter)
			defer slsWriter.Close()
		}
	}

	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()

	server.Use(log.Middleware)

	ctx := svc.NewServiceContext(c)

	system.Start(server, ctx)

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
