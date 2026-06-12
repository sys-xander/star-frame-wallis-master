package mp

import (
	"time"

	"golang.org/x/sync/singleflight"

	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"

	"wallpaper/internal/logic/mp"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"
)

var sf = new(singleflight.Group)

func RecommendWallpaperHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecommendWallpaperRst
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 优先从缓存读取
		cacheKey := "recommend_wallpaper"
		cacheValue, ok := svcCtx.Cache.Get(cacheKey)
		if ok {
			xhttp.JsonBaseResponseCtx(r.Context(), w, cacheValue.(*types.RecommendWallpaperRsp))
			return
		}

		l := mp.NewRecommendWallpaperLogic(r.Context(), svcCtx)

		val, err, _ := sf.Do(cacheKey, func() (interface{}, error) {
			resp, err := l.RecommendWallpaper(&req)
			if err != nil {
				return nil, err
			}
			svcCtx.Cache.SetWithTTL(cacheKey, resp, 0, time.Hour)
			return resp, nil
		})

		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, val)
		}
	}
}
