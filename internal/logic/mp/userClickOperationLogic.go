package mp

import (
	"context"
	"errors"

	"wallpaper/internal/model"
	"wallpaper/internal/svc"
	"wallpaper/internal/types"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	xerrors "github.com/zeromicro/x/errors"
	"gorm.io/gorm"
)

type UserClickOperationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserClickOperationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClickOperationLogic {
	return &UserClickOperationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserClickOperationLogic) UserClickOperation(req *types.UserClickOperationRst) (resp *types.UserClickOperationRsp, err error) {
	userID := cast.ToUint64(l.ctx.Value("user_id"))

	var user model.MpUser
	err = l.svcCtx.Orm.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, xerrors.New(901, "用户不存在")
	}

	tx := l.svcCtx.Orm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	switch req.Operation {
	case "like":
		err = l.handleLikeOperation(tx, userID, req.WallpaperId, req.OperationType)
	case "collect":
		err = l.handleCollectOperation(tx, userID, req.WallpaperId, req.OperationType)
	case "download":
		err = l.handleDownloadOperation(tx, userID, req.WallpaperId, req.OperationType)
	default:
		tx.Rollback()
		return nil, xerrors.New(901, "操作不存在")
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &types.UserClickOperationRsp{}, nil
}

// handleLikeOperation 处理点赞/取消点赞
func (l *UserClickOperationLogic) handleLikeOperation(tx *gorm.DB, userID, wallpaperID uint64, opType int) error {
	if opType == 1 { // 执行点赞
		err := tx.Create(&model.UserLike{UserID: userID, WallpaperID: wallpaperID}).Error
		if err != nil {
			// 判断是否唯一键冲突
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				// 已点赞过，正常返回即可
				return nil
			}
			return xerrors.New(902, "点赞记录创建失败")
		}
		if err := tx.Model(&model.Wallpaper{}).Where("id = ?", wallpaperID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
			return xerrors.New(902, "更新点赞数失败")
		}
	} else { // 取消点赞
		result := tx.Where("user_id = ? AND wallpaper_id = ?", userID, wallpaperID).Delete(&model.UserLike{})
		if result.Error != nil {
			return xerrors.New(902, "取消点赞失败")
		}
		if result.RowsAffected == 0 {
			return xerrors.New(904, "未点赞无法取消")
		}
		if err := tx.Model(&model.Wallpaper{}).Where("id = ?", wallpaperID).UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error; err != nil {
			return xerrors.New(902, "更新点赞数失败")
		}
	}
	return nil
}

// handleCollectOperation 处理收藏/取消收藏
func (l *UserClickOperationLogic) handleCollectOperation(tx *gorm.DB, userID, wallpaperID uint64, opType int) error {
	if opType == 1 { // 执行收藏
		err := tx.Create(&model.UserCollect{UserID: userID, WallpaperID: wallpaperID}).Error
		if err != nil {
			// 判断是否唯一键冲突
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				// 已收藏过，正常返回即可
				return nil
			}
			return xerrors.New(902, "收藏记录创建失败")
		}
		if err := tx.Model(&model.Wallpaper{}).Where("id = ?", wallpaperID).UpdateColumn("collect_count", gorm.Expr("collect_count + 1")).Error; err != nil {
			return xerrors.New(902, "更新收藏数失败")
		}
	} else { // 取消收藏
		result := tx.Where("user_id = ? AND wallpaper_id = ?", userID, wallpaperID).Delete(&model.UserCollect{})
		if result.Error != nil {
			return xerrors.New(902, "取消收藏失败")
		}
		if result.RowsAffected == 0 {
			return xerrors.New(904, "未收藏无法取消")
		}
		if err := tx.Model(&model.Wallpaper{}).Where("id = ?", wallpaperID).UpdateColumn("collect_count", gorm.Expr("GREATEST(collect_count - 1, 0)")).Error; err != nil {
			return xerrors.New(902, "更新收藏数失败")
		}
	}
	return nil
}

// handleDownloadOperation 处理下载/取消下载
func (l *UserClickOperationLogic) handleDownloadOperation(tx *gorm.DB, userID, wallpaperID uint64, opType int) error {
	if opType == 1 { // 执行下载
		err := tx.Create(&model.UserDownload{UserID: userID, WallpaperID: wallpaperID}).Error
		if err != nil {
			// 判断是否唯一键冲突
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				// 已下载过，正常返回即可
				return nil
			}
			return xerrors.New(902, "下载记录创建失败")
		}
		if err := tx.Model(&model.Wallpaper{}).Where("id = ?", wallpaperID).UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error; err != nil {
			return xerrors.New(902, "更新下载数失败")
		}
	} else { // 取消下载
		result := tx.Where("user_id = ? AND wallpaper_id = ?", userID, wallpaperID).Delete(&model.UserDownload{})
		if result.Error != nil {
			return xerrors.New(902, "取消下载失败")
		}
		if result.RowsAffected == 0 {
			return xerrors.New(904, "未下载无法取消")
		}
		if err := tx.Model(&model.Wallpaper{}).Where("id = ?", wallpaperID).UpdateColumn("download_count", gorm.Expr("GREATEST(download_count - 1, 0)")).Error; err != nil {
			return xerrors.New(902, "更新下载数失败")
		}
	}
	return nil
}
