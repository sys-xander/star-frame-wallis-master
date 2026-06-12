package mp

import (
	"context"

	"wallpaper/internal/model"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoRst) (resp *types.UpdateUserInfoRsp, err error) {
	
	// 进行校验，如果req.UserName和req.AvatarUrl和req.Phone都为空，则返回错误
	if req.UserName == "" && req.AvatarUrl == "" && req.Phone == "" {
		return nil, xerrors.New(902, "更新用户信息失败")
	}

	var user model.MpUser

	err = l.svcCtx.Orm.Where("id = ?", l.ctx.Value("user_id")).First(&user).Error;
	if err != nil {
		return nil, xerrors.New(901, "用户不存在")
	}

	if req.UserName != "" {
		user.Name = req.UserName
	}

	if req.AvatarUrl != "" {
		user.AvatarUrl = req.AvatarUrl
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	err = l.svcCtx.Orm.Save(&user).Error;
	
	if err != nil {
		return nil, xerrors.New(902, "更新用户信息失败")
	}

	return
}
