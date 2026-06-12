package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/pkg/cache"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type CategoriesListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoriesListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoriesListLogic {
	return &CategoriesListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoriesListLogic) CategoriesList(req *types.CategoriesListRst) (resp []types.CategoriesListRsp, err error) {

	var classifyList []model.Classification

	// 从缓存中获取分类列表
	classifyList, err = cache.GetClassification(l.svcCtx.Cache)

	if err != nil {
		err = l.svcCtx.Orm.Find(&classifyList).Error

		if err != nil {
			return nil, xerrors.New(901, "获取分类列表失败")
		}
		cache.SetClassification(l.svcCtx.Cache, classifyList)
	}

	var classifyListRsp = make([]types.CategoriesListRsp, 0, len(classifyList))
	for _, classify := range classifyList {
		classifyListRsp = append(classifyListRsp, types.CategoriesListRsp{
			ID:       classify.ID,
			Name:     classify.Name,
			Url:      l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle(classify.Url),
			Count:    classify.Count,
			Icon:     classify.Icon,
			Gradient: classify.Gradient,
		})
	}

	return classifyListRsp, nil
}
