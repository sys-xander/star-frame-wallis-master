package common

import (
	"context"
	"time"

	"wallpaper/internal/model"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
)

type AnnouncementListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAnnouncementListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnnouncementListLogic {
	return &AnnouncementListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AnnouncementListLogic) AnnouncementList(req *types.AnnouncementListRst) (resp []types.AnnouncementListRsp, err error) {
	
	var announcements []model.Announcement

	err = l.svcCtx.Orm.Where("status = 1").Order("sort DESC").Find(&announcements).Error
	if err != nil {
		return nil, xerrors.New(901, "获取公告列表失败")
	}

	var announcementListRsp []types.AnnouncementListRsp
	for _, announcement := range announcements {
		announcementListRsp = append(announcementListRsp, types.AnnouncementListRsp{
			ID:        announcement.ID,
			Title:     announcement.Title,
			Content:   announcement.Content,
			Type:      announcement.Type,
			Status:    announcement.Status,
			Cover:     l.svcCtx.QiniuYun.AuthUrlWithPreviewStyle(announcement.Cover),
			CreatedAt: announcement.CreatedAt.Format(time.DateTime),
			UpdatedAt: announcement.UpdatedAt.Format(time.DateTime),
		})
	}

	return announcementListRsp, nil
}
