package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type SearchWallpaperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchWallpaperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchWallpaperLogic {
	return &SearchWallpaperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchWallpaperLogic) SearchWallpaper(req *types.SearchWallpaperRst) (resp *types.SearchWallpaperRsp, err error) {
	// 搜索壁纸
	var wallpaperList []model.Wallpaper
	err = l.svcCtx.Orm.Model(&model.Wallpaper{}).Where("title LIKE ?", "%"+req.Keyword+"%").
		Order("created_at DESC,id DESC").Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).Find(&wallpaperList).Error
	if err != nil {
		return nil, xerrors.New(901, "搜索壁纸失败")
	}

	var wallpaperListRsp []types.WallpaperListRsp
	for _, wallpaper := range wallpaperList {
		wallpaperListRsp = append(wallpaperListRsp, types.WallpaperListRsp{
			ID:            wallpaper.ID,
			Title:         wallpaper.Title,
			Url:           l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle(wallpaper.Url),
			LikeCount:     wallpaper.LikeCount,
			DownloadCount: wallpaper.DownloadCount,
			CollectCount:  wallpaper.CollectCount,
			ViewCount:     wallpaper.ViewCount,
		})
	}
	return &types.SearchWallpaperRsp{
		WallpaperList: wallpaperListRsp,
	}, nil
}
