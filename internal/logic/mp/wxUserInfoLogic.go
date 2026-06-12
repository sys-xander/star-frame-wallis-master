package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/pkg/strs"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type WxUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWxUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WxUserInfoLogic {
	return &WxUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WxUserInfoLogic) WxUserInfo(req *types.WxUserInfoRst) (resp *types.WxUserInfoRsp, err error) {
	
	var user model.MpUser

	err = l.svcCtx.Orm.Where("id = ?", l.ctx.Value("user_id")).First(&user).Error;
	if err != nil {
		return nil, xerrors.New(901, "用户不存在")
	}

	return &types.WxUserInfoRsp{
		ID:        user.ID,
		UserName:  user.Name,
		AvatarUrl: l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle(user.AvatarUrl),
		Phone:     strs.HidePhone(user.Phone),
	}, nil
}
