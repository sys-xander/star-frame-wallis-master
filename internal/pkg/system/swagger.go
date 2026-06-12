package system

import (
	"net/http"
	"wallpaper/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func StartSwagger(server *rest.Server, ctx *svc.ServiceContext) {
	if ctx.Config.Mode == "dev" {
		server.AddRoute(rest.Route{
			Method:  http.MethodGet,
			Path:    "/swagger.json",
			Handler: http.FileServer(http.Dir(".")).ServeHTTP,
		})
	}
}
