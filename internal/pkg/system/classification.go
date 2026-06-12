package system

import (
	"wallpaper/internal/model"
	"wallpaper/internal/pkg/cache"
	"wallpaper/internal/svc"
)

func InitClassification(ctx *svc.ServiceContext) error {
	var classifyList []model.Classification
	err := ctx.Orm.Find(&classifyList).Error
	if err != nil {
		return err
	}
	cache.SetClassification(ctx.Cache, classifyList)
	return nil
}
