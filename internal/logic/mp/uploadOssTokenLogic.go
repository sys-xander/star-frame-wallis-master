package mp

import (
	"context"
	"fmt"
	"time"

	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type UploadOssTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadOssTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadOssTokenLogic {
	return &UploadOssTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadOssToken 生成上传时的oss token 和 url
func (l *UploadOssTokenLogic) UploadOssToken(req *types.UploadOssTokenRst) (resp *types.UploadOssTokenRsp, err error) {
	var key string

	// 获取用户id
	userId := l.ctx.Value("user_id")

	switch req.TypeKey {
	case "avatar":
		// 使用用户id和时间戳生成key
		key = fmt.Sprintf("%s-%d.%s", userId, time.Now().Unix(), req.Format)
	default:
		return nil, xerrors.New(901, "key不合法")
	}

	token, fullKey, err := l.svcCtx.QiniuYun.UploadUserAvatarToken(key)
	if err != nil {
		return nil, xerrors.New(902, "生成oss token失败")
	}

	return &types.UploadOssTokenRsp{
		Token: token,
		Url:   l.svcCtx.QiniuYun.AuthUrlNoStyle(fullKey),
		Key:   fullKey,
	}, nil
}
