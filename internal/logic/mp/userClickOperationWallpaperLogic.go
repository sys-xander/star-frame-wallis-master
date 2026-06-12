package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type UserClickOperationWallpaperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserClickOperationWallpaperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClickOperationWallpaperLogic {
	return &UserClickOperationWallpaperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserClickOperationWallpaperLogic) UserClickOperationWallpaper(req *types.UserClickOperationWallpaperRst) (resp *types.UserClickOperationWallpaperRsp, err error) {

	var wallpaperListRsp []types.WallpaperListRsp
	// 获取用户收藏/下载/点赞的壁纸信息
	switch req.Operation {
	case "like":
		wallpaperListRsp, err = l.handleLikeOperation(req.Page, req.PageSize)
	case "collect":
		wallpaperListRsp, err = l.handleCollectOperation(req.Page, req.PageSize)
	case "download":
		wallpaperListRsp, err = l.handleDownloadOperation(req.Page, req.PageSize)
	}

	if err != nil {
		return nil, err
	}

	return &types.UserClickOperationWallpaperRsp{
		WallpaperList: wallpaperListRsp,
	}, nil
}

func (l *UserClickOperationWallpaperLogic) handleLikeOperation(page, pageSize int) ([]types.WallpaperListRsp, error) {
	userID := cast.ToUint64(l.ctx.Value("user_id"))

	var userLikes []model.UserLike
	err := l.svcCtx.Orm.Model(&model.UserLike{}).Where("user_id = ?", userID).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&userLikes).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户收藏/下载/点赞的壁纸信息失败")
	}

	wallpaperIDs := lo.Map(userLikes, func(item model.UserLike, _ int) uint64 {
		return item.WallpaperID
	})

	var wallpaperList []model.Wallpaper
	err = l.svcCtx.Orm.Model(&model.Wallpaper{}).Where("id IN (?)", wallpaperIDs).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&wallpaperList).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户收藏/下载/点赞的壁纸信息失败")
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
	return wallpaperListRsp, nil
}

func (l *UserClickOperationWallpaperLogic) handleCollectOperation(page, pageSize int) ([]types.WallpaperListRsp, error) {
	userID := cast.ToUint64(l.ctx.Value("user_id"))
	var userCollects []model.UserCollect
	err := l.svcCtx.Orm.Model(&model.UserCollect{}).Where("user_id = ?", userID).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&userCollects).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户收藏/下载/点赞的壁纸信息失败")
	}

	wallpaperIDs := lo.Map(userCollects, func(item model.UserCollect, _ int) uint64 {
		return item.WallpaperID
	})

	var wallpaperList []model.Wallpaper
	err = l.svcCtx.Orm.Model(&model.Wallpaper{}).Where("id IN (?)", wallpaperIDs).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&wallpaperList).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户收藏/下载/点赞的壁纸信息失败")
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
	return wallpaperListRsp, nil
}

func (l *UserClickOperationWallpaperLogic) handleDownloadOperation(page, pageSize int) ([]types.WallpaperListRsp, error) {
	userID := cast.ToUint64(l.ctx.Value("user_id"))
	var userDownloads []model.UserDownload
	err := l.svcCtx.Orm.Model(&model.UserDownload{}).Where("user_id = ?", userID).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&userDownloads).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户收藏/下载/点赞的壁纸信息失败")
	}

	wallpaperIDs := lo.Map(userDownloads, func(item model.UserDownload, _ int) uint64 {
		return item.WallpaperID
	})

	var wallpaperList []model.Wallpaper
	err = l.svcCtx.Orm.Model(&model.Wallpaper{}).Where("id IN (?)", wallpaperIDs).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&wallpaperList).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户收藏/下载/点赞的壁纸信息失败")
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
	return wallpaperListRsp, nil
}
