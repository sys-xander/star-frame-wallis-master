package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/pkg/cache"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
	"gorm.io/gorm"
)

type WallpaperDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWallpaperDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WallpaperDetailLogic {
	return &WallpaperDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WallpaperDetailLogic) WallpaperDetail(req *types.WallpaperDetailRst) (resp *types.WallpaperDetailRsp, err error) {
	var wallpaper model.Wallpaper
	err = l.svcCtx.Orm.Where("id = ?", req.ID).First(&wallpaper).Error
	if err != nil {
		return nil, xerrors.New(901, "获取壁纸详情失败")
	}

	var classification string
	classification, err = cache.GetClassificationById(l.svcCtx.Cache, wallpaper.ClassifyId)
	if err != nil {
		return nil, xerrors.New(901, "获取壁纸分类失败")
	}

	// 随机获取这个分类下的9张壁纸
	var recommendList []model.Wallpaper
	err = l.svcCtx.Orm.Where("classify_id = ?", wallpaper.ClassifyId).Order("RAND()").Limit(9).Find(&recommendList).Error
	if err != nil {
		return nil, xerrors.New(901, "获取推荐壁纸失败")
	}

	var recommendListRsp []types.RecommendWallpaper
	for _, wallpaper := range recommendList {
		recommendListRsp = append(recommendListRsp, types.RecommendWallpaper{
			ID:            wallpaper.ID,
			Title:         wallpaper.Title,
			Url:           l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle(wallpaper.Url),
			LikeCount:     wallpaper.LikeCount,
			DownloadCount: wallpaper.DownloadCount,
			CollectCount:  wallpaper.CollectCount,
			ViewCount:     wallpaper.ViewCount,
		})
	}

	// 浏览量+1 (使用数据库原子操作，避免并发问题)
	l.svcCtx.Orm.Model(&model.Wallpaper{}).Where("id = ?", wallpaper.ID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	return &types.WallpaperDetailRsp{
		ID:            wallpaper.ID,
		Title:         wallpaper.Title,
		Url:           l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle(wallpaper.Url),
		CleanUrl:      l.svcCtx.QiniuYun.AuthUrlNoStyle(wallpaper.Url),
		LikeCount:     wallpaper.LikeCount,
		ViewCount:     wallpaper.ViewCount,
		DownloadCount: wallpaper.DownloadCount,
		CollectCount:  wallpaper.CollectCount,
		CreatedAt:     wallpaper.CreatedAt.Format("2006-01-02 15:04:05"),
		Classify:      classification,
		RecommendList: recommendListRsp,
	}, nil
}
