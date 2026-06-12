package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type WallpaperListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWallpaperListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WallpaperListLogic {
	return &WallpaperListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WallpaperListLogic) WallpaperList(req *types.WallpaperListRst) (resp []types.WallpaperListRsp, err error) {

	var wallpapers []model.Wallpaper
	err = l.svcCtx.Orm.Where("classify_id = ?", req.ClassifyId).
		Order("created_at DESC, id DESC").Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).Find(&wallpapers).Error

	if err != nil {
		return nil, xerrors.New(901, "获取壁纸列表失败")
	}

	var wallpaperListRsp []types.WallpaperListRsp
	for _, wallpaper := range wallpapers {
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
	return wallpaperListRsp, nil
}
