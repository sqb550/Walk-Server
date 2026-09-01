package repo

import (
	"context"
	"errors"
	"time"

	"app/comm"
	"app/dao/model"

	"github.com/zjutjh/mygo/ndb"
	"gorm.io/gorm"
)

const unreadNoticeLimit = 50

type NoticeRepo struct {
	db *gorm.DB
}

func NewNoticeRepo() *NoticeRepo {
	return &NoticeRepo{db: ndb.Pick()}
}

func NewNoticeRepoWithDB(db *gorm.DB) *NoticeRepo {
	return &NoticeRepo{db: db}
}

func (r *NoticeRepo) ListUnread(ctx context.Context, userID int64) ([]model.NoticeRecord, error) {
	var notices []model.NoticeRecord
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND read_at IS NULL", userID).
		Order("created_at ASC, id ASC").
		Limit(unreadNoticeLimit).
		Find(&notices).Error
	return notices, err
}

func (r *NoticeRepo) Ack(ctx context.Context, userID, noticeID int64) error {
	return r.db.WithContext(ctx).Model(&model.NoticeRecord{}).
		Where("user_id = ? AND id = ? AND read_at IS NULL", userID, noticeID).
		Update("read_at", time.Now()).Error
}

func (r *NoticeRepo) UpsertUnread(ctx context.Context, notices []model.NoticeRecord) error {
	if len(notices) == 0 {
		return nil
	}
	for index := range notices {
		notice := &notices[index]
		var existing model.NoticeRecord
		err := r.db.WithContext(ctx).
			Where("user_id = ? AND type = ? AND read_at IS NULL", notice.UserID, notice.Type).
			Order("id DESC").
			First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.WithContext(ctx).Create(notice).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if err := r.db.WithContext(ctx).Model(&existing).Updates(map[string]any{
			"actor_id":   notice.ActorID,
			"actor_name": notice.ActorName,
			"team_id":    notice.TeamID,
			"created_at": time.Now(),
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *NoticeRepo) DeleteUnreadTypes(ctx context.Context, userID int64, types ...comm.NoticeType) error {
	if len(types) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Where("user_id = ? AND type IN ? AND read_at IS NULL", userID, types).
		Delete(&model.NoticeRecord{}).Error
}
