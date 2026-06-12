package login

import (
	"context"
	"errors"
	"fmt"
	"wallpaper/internal/model"
	"wallpaper/internal/pkg/strs"
	"wallpaper/internal/pkg/token"

	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	xerrors "github.com/zeromicro/x/errors"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type WxLoginCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWxLoginCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WxLoginCodeLogic {
	return &WxLoginCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// WxLoginCode 微信小程序Code登录 微信小程序一键登录
func (l *WxLoginCodeLogic) WxLoginCode(req *types.WxLoginCodeRst) (*types.WxLoginCodeRsp, error) {

	ctx := context.Background()

	ret, err := l.svcCtx.MiniProgram.Auth.Session(ctx, req.Code)
	
	if err != nil {
		l.Errorf("[WxLoginCode] get session key error: %v", err)
		return nil, xerrors.New(901, "登录失败")
	}

	var user model.MpUser	
	err = l.svcCtx.Orm.Where("openid = ?", ret.OpenID).Take(&user).Error; 

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 用户不存在 创建新用户

		user = model.MpUser{
			Name:       fmt.Sprintf("user-%s", strs.RandomString(10)),
			Openid:     ret.OpenID,
			Unionid:    ret.UnionID,
			SessionKey: ret.SessionKey,
			AvatarUrl:  "",
		}

		if err = l.svcCtx.Orm.Create(&user).Error; err != nil {
			l.Errorf("[WxLoginCode] create user error: %v", err)
			return nil, xerrors.New(903, "登录失败")
		}

		// 签发 access token
		tokenStr, expireAt, err := token.GenerateToken(&l.svcCtx.Config.Auth, token.MyClaims{
			UserID:     user.ID,
			Name:       user.Name,
			Unionid:    user.Unionid,
			SessionKey: user.SessionKey,
		})
		if err != nil {
			l.Errorf("[WxLoginCode] generate token error: %v", err)
			return nil, xerrors.New(904, "登录失败")
		}

		// 签发 refresh token
		refreshTokenStr, err := token.GenerateRefreshToken(&l.svcCtx.Config.Auth, user.ID)
		if err != nil {
			l.Errorf("[WxLoginCode] generate refresh token error: %v", err)
			return nil, xerrors.New(904, "登录失败")
		}

		return &types.WxLoginCodeRsp{
			ID:           user.ID,
			Token:        tokenStr,
			RefreshToken: refreshTokenStr,
			UserName:     user.Name,
			AvatarUrl:    l.svcCtx.QiniuYun.AuthUrlNoStyle(user.AvatarUrl),
			Phone:        "",
			ExpireAt:     expireAt,
		}, nil
	}

	// 用户存在更新 session_key
	if err = l.svcCtx.Orm.Model(&model.MpUser{}).Where("id = ?", user.ID).UpdateColumn("session_key", ret.SessionKey).Error; err != nil {
		l.Errorf("[WxLoginCode] update user session_key error: %v", err)
		return nil, xerrors.New(905, "登录失败")
	}

	// 签发 access token
	tokenStr, expireAt, err := token.GenerateToken(&l.svcCtx.Config.Auth, token.MyClaims{
		UserID:     user.ID,
		Name:       user.Name,
		Unionid:    user.Unionid,
		SessionKey: user.SessionKey,
	})
	if err != nil {
		l.Errorf("[WxLoginCode] generate token error: %v", err)
		return nil, xerrors.New(906, "登录失败")
	}

	// 签发 refresh token
	refreshTokenStr, err := token.GenerateRefreshToken(&l.svcCtx.Config.Auth, user.ID)
	if err != nil {
		l.Errorf("[WxLoginCode] generate refresh token error: %v", err)
		return nil, xerrors.New(906, "登录失败")
	}

	return &types.WxLoginCodeRsp{
		ID:           user.ID,
		Token:        tokenStr,
		RefreshToken: refreshTokenStr,
		UserName:     user.Name,
		AvatarUrl:    l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle(user.AvatarUrl),
		Phone:        strs.HidePhone(user.Phone),
		ExpireAt:     expireAt,
	}, nil
}