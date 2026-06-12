package convert

import (
	"wallpaper/internal/model"
	"wallpaper/internal/types"
)

// WallpaperToRecommend 将单个 model.Wallpaper 转换为 types.RecommendWallpaper
// urlHandler: 用于处理 URL 的函数，如添加七牛云鉴权等
func WallpaperToRecommend(w model.Wallpaper, urlHandler func(string) string) types.RecommendWallpaper {
	url := w.Url
	if urlHandler != nil {
		url = urlHandler(url)
	}

	return types.RecommendWallpaper{
		ID:            w.ID,
		Title:         w.Title,
		Url:           url,
		LikeCount:     w.LikeCount,
		DownloadCount: w.DownloadCount,
		CollectCount:  w.CollectCount,
		ViewCount:     w.ViewCount,
	}
}

// WallpaperSliceToRecommend 将 []model.Wallpaper 转换为 []types.RecommendWallpaper
// urlHandler: 用于处理 URL 的函数，如添加七牛云鉴权等
func WallpaperSliceToRecommend(wallpapers []model.Wallpaper, urlHandler func(string) string) []types.RecommendWallpaper {
	if len(wallpapers) == 0 {
		return []types.RecommendWallpaper{}
	}

	result := make([]types.RecommendWallpaper, 0, len(wallpapers))
	for _, w := range wallpapers {
		result = append(result, WallpaperToRecommend(w, urlHandler))
	}
	return result
}
