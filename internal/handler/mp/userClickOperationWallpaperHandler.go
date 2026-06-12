package mp

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"

	"wallpaper/internal/logic/mp"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"
)

func UserClickOperationWallpaperHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserClickOperationWallpaperRst
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := mp.NewUserClickOperationWallpaperLogic(r.Context(), svcCtx)
		resp, err := l.UserClickOperationWallpaper(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
