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

type UserClickOperationCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserClickOperationCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClickOperationCountLogic {
	return &UserClickOperationCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserClickOperationCountLogic) UserClickOperationCount(req *types.UserClickOperationCountRst) (resp *types.UserClickOperationCountRsp, err error) {
	// 获取用户点赞数 收藏数 下载数
	userID := cast.ToUint64(l.ctx.Value("user_id"))

	var likeCount int64
	var collectCount int64
	var downloadCount int64

	err = l.svcCtx.Orm.Model(&model.UserLike{}).Where("user_id = ?", userID).Count(&likeCount).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户点赞数失败")
	}

	err = l.svcCtx.Orm.Model(&model.UserCollect{}).Where("user_id = ?", userID).Count(&collectCount).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户收藏数失败")
	}

	err = l.svcCtx.Orm.Model(&model.UserDownload{}).Where("user_id = ?", userID).Count(&downloadCount).Error
	if err != nil {
		return nil, xerrors.New(901, "获取用户下载数失败")
	}

	return &types.UserClickOperationCountRsp{
		LikeCount:     likeCount,
		CollectCount:  collectCount,
		DownloadCount: downloadCount,
	}, nil
}
