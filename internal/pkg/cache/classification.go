package cache

import (
	"errors"
	"time"
	"wallpaper/internal/model"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/singleflight"
)

// 实现对壁纸分类的缓存和查询

const CacheClassificationKey = "classification"

var sf = new(singleflight.Group)

func GetClassification(cache *ristretto.Cache[string, any]) ([]model.Classification, error) {
	cacheValue, ok := cache.Get(CacheClassificationKey)
	if ok {
		return cacheValue.([]model.Classification), nil
	}
	return nil, errors.New("cache classification not found")
}

func SetClassification(cache *ristretto.Cache[string, any], classify []model.Classification) {
	sf.Do(CacheClassificationKey, func() (interface{}, error) {
		ok := cache.SetWithTTL(CacheClassificationKey, classify, 0, time.Hour*24)
		if !ok {
			logx.Errorw("set classification failed", logx.Field("CacheClassificationKey", CacheClassificationKey), logx.Field("classify", classify))
		}
		return nil, nil
	})
}

func GetClassificationById(cache *ristretto.Cache[string, any], id int) (string, error) {
	cacheValue, ok := cache.Get(CacheClassificationKey)
	if !ok{
		return "", errors.New("cache classification not found")
	}
	for _, classify := range cacheValue.([]model.Classification) {
		if classify.ID == uint(id) {
			return classify.Name, nil
		}
	}
	return "", errors.New("classification not found")
}
