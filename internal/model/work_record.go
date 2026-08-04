package model

import "github.com/google/uuid"

type TaskWorkRecord struct {
	BaseEntity
	Completable

	WorkTaskID     uuid.UUID `gorm:"type:char(36);not null" json:"work_task_id"`
	UserID         uuid.UUID `gorm:"type:char(36);not null" json:"user_id"`
	DurationMinutes int       `json:"duration_minutes"`
}

func (TaskWorkRecord) TableName() string { return "task_work_records" }

func (r TaskWorkRecord) IsActive() bool { return r.FinishedAt == nil }