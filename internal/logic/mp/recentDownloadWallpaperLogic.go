package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type RecentDownloadWallpaperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecentDownloadWallpaperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecentDownloadWallpaperLogic {
	return &RecentDownloadWallpaperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecentDownloadWallpaperLogic) RecentDownloadWallpaper(req *types.RecentDownloadWallpaperRst) (resp *types.RecentDownloadWallpaperRsp, err error) {
	userID := cast.ToUint64(l.ctx.Value("user_id"))
	
	var userDownloads []model.UserDownload
	err = l.svcCtx.Orm.Where("user_id = ?", userID).Order("created_at DESC").Limit(4).Find(&userDownloads).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户最近下载的壁纸失败")
	}

	var wallpaperIDs []uint64
	for _, userDownload := range userDownloads {
		wallpaperIDs = append(wallpaperIDs, userDownload.WallpaperID)
	}

	var wallpaperList []model.Wallpaper
	err = l.svcCtx.Orm.Where("id IN (?)", wallpaperIDs).Find(&wallpaperList).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户最近下载的壁纸失败")
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

	return &types.RecentDownloadWallpaperRsp{
		WallpaperList: wallpaperListRsp,
	}, nil
}
