package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
	"gorm.io/gorm"
)

type UserClickOperationCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserClickOperationCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClickOperationCheckLogic {
	return &UserClickOperationCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserClickOperationCheckLogic) UserClickOperationCheck(req *types.UserClickOperationCheckRst) (resp *types.UserClickOperationCheckRsp, err error) {

	userID := cast.ToUint64(l.ctx.Value("user_id"))

	var user model.MpUser
	err = l.svcCtx.Orm.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, xerrors.New(901, "用户不存在")
	}

	// 判断壁纸是否存在
	err = l.svcCtx.Orm.Model(&model.Wallpaper{}).Where("id = ?", req.WallpaperId).First(&model.Wallpaper{}).Error
	if err != nil {
		return nil, xerrors.New(901, "壁纸不存在")
	}

	var (
		isLike     bool
		isCollect  bool
		isDownload bool
	)

	err = l.svcCtx.Orm.Model(&model.UserLike{}).Where("user_id = ? AND wallpaper_id = ?", userID, req.WallpaperId).First(&model.UserLike{}).Error

	switch err {
	case gorm.ErrRecordNotFound:
		isLike = false
	case nil:
		isLike = true
	default:
		logx.Errorf("[UserClickOperationCheckLogic] get user like error: %v", err)
		return nil, xerrors.New(901, "获取用户点赞状态失败")
	}

	err = l.svcCtx.Orm.Model(&model.UserCollect{}).Where("user_id = ? AND wallpaper_id = ?", userID, req.WallpaperId).First(&model.UserCollect{}).Error
	switch err {
	case gorm.ErrRecordNotFound:
		isCollect = false
	case nil:
		isCollect = true
	default:
		logx.Errorf("[UserClickOperationCheckLogic] get user collect error: %v", err)
		return nil, xerrors.New(901, "获取用户收藏状态失败")
	}

	err = l.svcCtx.Orm.Model(&model.UserDownload{}).Where("user_id = ? AND wallpaper_id = ?", userID, req.WallpaperId).First(&model.UserDownload{}).Error
	switch err {
	case gorm.ErrRecordNotFound:
		isDownload = false
	case nil:
		isDownload = true
	default:
		logx.Errorf("[UserClickOperationCheckLogic] get user download error: %v", err)
		return nil, xerrors.New(901, "获取用户下载状态失败")
	}

	return &types.UserClickOperationCheckRsp{
		IsLike:     isLike,
		IsCollect:  isCollect,
		IsDownload: isDownload,
	}, nil
}
