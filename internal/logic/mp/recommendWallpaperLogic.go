package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/pkg/convert"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type RecommendWallpaperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendWallpaperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendWallpaperLogic {
	return &RecommendWallpaperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendWallpaperLogic) RecommendWallpaper(req *types.RecommendWallpaperRst) (resp *types.RecommendWallpaperRsp, err error) {
	// 获取最新的9条壁纸
	var newList []model.Wallpaper
	err = l.svcCtx.Orm.Order("created_at DESC").Limit(9).Find(&newList).Error
	if err != nil {
		return nil, xerrors.New(901, "获取最新壁纸失败")
	}

	// 获取最热的9条壁纸
	var hotList []model.Wallpaper
	err = l.svcCtx.Orm.Order("like_count DESC").Limit(9).Find(&hotList).Error
	if err != nil {
		return nil, xerrors.New(901, "获取最热壁纸失败")
	}

	// 获取精选的9条壁纸
	var selectedList []model.Wallpaper
	err = l.svcCtx.Orm.Order("collect_count DESC").Limit(9).Find(&selectedList).Error
	if err != nil {
		return nil, xerrors.New(901, "获取精选壁纸失败")
	}

	// 获取推荐的9条壁纸 暂无推荐系统，从数据库中随机拿9条
	var recommendList []model.Wallpaper
	err = l.svcCtx.Orm.Order("RAND()").Limit(9).Find(&recommendList).Error
	if err != nil {
		return nil, xerrors.New(901, "获取推荐壁纸失败")
	}

	return &types.RecommendWallpaperRsp{
		RecommendList: convert.WallpaperSliceToRecommend(recommendList, l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle),
		NewList:       convert.WallpaperSliceToRecommend(newList, l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle),
		HotList:       convert.WallpaperSliceToRecommend(hotList, l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle),
		SelectedList:  convert.WallpaperSliceToRecommend(selectedList, l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle),
	}, nil
}
