package qiniu

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wallpaper/internal/config"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

const (
	ossAvatarPrefix = "avatar/"
	ossStylePreview = "preview"
)

type QiniuYun struct {
	config config.Qiniu
	cred   *credentials.Credentials
	mac    *auth.Credentials
}

func NewQiniuYun(config config.Qiniu) *QiniuYun {
	// 创建凭证（AK/SK 从 https://portal.qiniu.com/user/key 获取）
	cred := credentials.NewCredentials(config.AccessKey, config.SecretKey)

	mac := auth.New(config.AccessKey, config.SecretKey)

	return &QiniuYun{
		config: config,
		cred:   cred,
		mac:    mac,
	}
}

// 上传策略 https://developer.qiniu.com/kodo/1206/put-policy
// scope 上传空间 使用 <bucket>:<key> 格式
// UploadUserAvatarToken 生成上传用户头像的oss token
func (qy *QiniuYun) UploadUserAvatarToken(key string) (string, string, error) {
	fullKey := fmt.Sprintf("%s%s", ossAvatarPrefix, key)
	putPolicy, _ := uptoken.NewPutPolicyWithKeyPrefix(qy.config.Bucket, fullKey, time.Now().Add(time.Hour))
	// isPrefixalScope 为 1 时无法覆盖上传。oss已开启cdn，如开启覆盖上传，需刷新cdn缓存
	putPolicy.SetIsPrefixalScope(1)
	// 限制图片大小
	putPolicy.SetFsizeLimit(1024 * 1024 * 1)
	signer := uptoken.NewSigner(putPolicy, qy.cred)

	token, err := signer.GetUpToken(context.Background())
	if err != nil {
		return "", "", err
	}

	return token, fullKey, nil
}

// AuthUrl 生成七牛云的鉴权URL 私有空间url签名
// url 图片地址
// style 样式
// expires 过期时间
func (qy *QiniuYun) authUrl(url, style string, expires time.Duration) string {
	if url == "" {
		return ""
	}

	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
	} else {
		url = fmt.Sprintf("%s/%s", qy.config.Host, url)
	}

	var ex time.Time

	if expires <= 0 {
		ex = time.Now().Add(time.Duration(qy.config.Expire) * time.Second)
	} else {
		ex = time.Now().Add(expires)
	}

	if style == "" {
		url = fmt.Sprintf("%s?e=%d", url, ex.Unix())
	} else {
		url = fmt.Sprintf("%s-%s?e=%d", url, style, ex.Unix())
	}

	sign := qy.mac.Sign([]byte(url))

	return fmt.Sprintf("%s&token=%s", url, sign)
}

func (qy *QiniuYun) AuthUrlNoStyle(url string) string {
	return qy.authUrl(url, "", 0)
}

func (qy *QiniuYun) AuthUrlWithPreviewStyle(url string) string {
	return qy.authUrl(url, ossStylePreview, time.Hour*24)
}
