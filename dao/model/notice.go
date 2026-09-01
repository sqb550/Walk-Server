package model

import (
	"time"

	"app/comm"
)

const TableNameNotice = "notices"

// NoticeRecord is maintained manually until gorm/gen is next run against a
// database containing the notices table.
type NoticeRecord struct {
	ID        int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64           `gorm:"column:user_id;not null" json:"user_id"`
	Type      comm.NoticeType `gorm:"column:type;not null" json:"type"`
	ActorID   *int64          `gorm:"column:actor_id" json:"actor_id"`
	ActorName string          `gorm:"column:actor_name" json:"actor_name"`
	TeamID    *int64          `gorm:"column:team_id" json:"team_id"`
	ReadAt    *time.Time      `gorm:"column:read_at" json:"read_at"`
	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (*NoticeRecord) TableName() string { return TableNameNotice }
