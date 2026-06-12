package login

import (
	"context"
	"errors"

	"wallpaper/internal/model"
	"wallpaper/internal/pkg/token"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRst) (resp *types.RefreshTokenRsp, err error) {
	// 解析 refresh token
	claims, err := token.ParseRefreshToken(l.svcCtx.Config.Auth.RefreshSecret, req.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, xerrors.New(901, "刷新令牌已过期")
		}
		return nil, xerrors.New(901, "无效的刷新令牌")
	}

	// 查询用户是否存在
	var user model.MpUser
	if err = l.svcCtx.Orm.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		l.Errorf("[RefreshToken] user not found: %v", err)
		return nil, xerrors.New(902, "用户不存在")
	}

	// 签发新的 access token
	tokenStr, expireAt, err := token.GenerateToken(&l.svcCtx.Config.Auth, token.MyClaims{
		UserID:     user.ID,
		Name:       user.Name,
		Unionid:    user.Unionid,
		SessionKey: user.SessionKey,
	})
	if err != nil {
		l.Errorf("[RefreshToken] generate token error: %v", err)
		return nil, xerrors.New(903, "签发令牌失败")
	}

	// 签发新的 refresh token
	refreshTokenStr, err := token.GenerateRefreshToken(&l.svcCtx.Config.Auth, user.ID)
	if err != nil {
		l.Errorf("[RefreshToken] generate refresh token error: %v", err)
		return nil, xerrors.New(903, "签发令牌失败")
	}

	return &types.RefreshTokenRsp{
		Token:        tokenStr,
		RefreshToken: refreshTokenStr,
		ExpireAt:     expireAt,
	}, nil
}
