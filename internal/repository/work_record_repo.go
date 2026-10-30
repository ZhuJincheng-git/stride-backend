package repository

import (
	"context"
	"time"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkRecordRepository interface {
	Create(ctx context.Context, r *model.TaskWorkRecord) error
	GetActiveForUser(ctx context.Context, userID uuid.UUID) (*model.TaskWorkRecord, error)
	Stop(ctx context.Context, id uuid.UUID, finishedAt time.Time, durationMinutes int) error
	ListForTask(ctx context.Context, userID, taskID uuid.UUID) ([]model.TaskWorkRecord, error)
	TotalMinutesForTask(ctx context.Context, userID, taskID uuid.UUID) (int, error)
}

type workRecordRepo struct{ db *gorm.DB }

func NewWorkRecordRepository(db *gorm.DB) WorkRecordRepository { return &workRecordRepo{db: db} }

func (r *workRecordRepo) Create(ctx context.Context, rec *model.TaskWorkRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *workRecordRepo) GetActiveForUser(ctx context.Context, userID uuid.UUID) (*model.TaskWorkRecord, error) {
	var rec model.TaskWorkRecord
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND finished_at IS NULL", userID).
		Order("created_at DESC").
		First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *workRecordRepo) Stop(ctx context.Context, id uuid.UUID, finishedAt time.Time, durationMinutes int) error {
	res := r.db.WithContext(ctx).
		Model(&model.TaskWorkRecord{}).
		Where("id = ? AND finished_at IS NULL", id).
		Updates(map[string]any{
			"finished_at": finishedAt,
			"duration_minutes": durationMinutes,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *workRecordRepo) ListForTask(ctx context.Context, userID, taskID uuid.UUID) ([]model.TaskWorkRecord, error) {
	var out []model.TaskWorkRecord
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND work_task_id = ?", userID, taskID).
		Order("created_at DESC").
		Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *workRecordRepo) TotalMinutesForTask(ctx context.Context, userID, taskID uuid.UUID) (int, error) {
	var total *int
	if err := r.db.WithContext(ctx).
	Model(&model.TaskWorkRecord{}).
	Where("user_id = ? AND work_task_id = ? AND finished_at IS NOT NULL", userID, taskID).
	Select("COALESCE(SUM(duration_minutes), 0)").
	Scan(&total).Error; err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}
	return *total, nil
}