package repository

import (
	"context"
	"time"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskFilter struct {
	UserID         uuid.UUID
	GoalID         *uuid.UUID // nil = all goals (or none)
	OnlyCompleted  *bool
	OnlyDeleted    bool
	IncludeDeleted bool
	Limit          int
	Offset         int
}

type TaskRepository interface {
	Create(ctx context.Context, t *model.Task) error
	Update(ctx context.Context, t *model.Task) error
	GetByID(ctx context.Context, userID, id uuid.UUID, IncludeDeleted bool) (*model.Task, error)
	List(ctx context.Context, f TaskFilter) ([]model.Task, error)
	SoftDelete(ctx context.Context, userID, id uuid.UUID) error
	Restore(ctx context.Context, userID, id uuid.UUID) error
	HardDelete(ctx context.Context, userID, id uuid.UUID)
	SetFinished(ctx context.Context, userID, id uuid.UUID, finishedAt *time.Time) error
}

type taskRepo struct{ db *gorm.DB }

func NewTaskRepository(db *gorm.DB) TaskRepository { return &taskRepo{db: db} }

func (r *taskRepo) Create(ctx context.Context, t *model.Task) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *taskRepo) Update(ctx context.Context, t *model.Task) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *taskRepo) GetByID(ctx context.Context, userID, id uuid.UUID, IncludeDeleted bool) (*model.Task, error) {
	q := r.db.WithContext(ctx)
	if IncludeDeleted {
		q = q.Unscoped()
	}
	var t model.Task
	if err := q.First(&t, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *taskRepo) List(ctx context.Context, f TaskFilter) ([]model.Task, error) {
	q := r.db.WithContext(ctx).Model(&model.Task{})
	if f.OnlyDeleted {
		q = q.Unscoped().Where("deleted_at IS NOT NULL")
	} else if f.IncludeDeleted {
		q = q.Unscoped()
	}
	q = q.Where("user_id = ?", f.UserID)
	if f.GoalID != nil {
		q = q.Where("goal_id = ?", *f.GoalID)
	}
	if f.OnlyCompleted != nil {
		if *f.OnlyCompleted {
			q = q.Where("finished_at IS NOT NULL")
		} else {
			q = q.Where("finished_at IS NULL")
		}
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	if f.Offset > 0 {
		q = q.Offset(f.Offset)
	}

	var out []model.Task
	if err := q.Order("created_at DESC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *taskRepo) SoftDelete(ctx context.Context, userID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Task{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *taskRepo) Restore(ctx context.Context, userID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.Task{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", id, userID).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *taskRepo) HardDelete(ctx context.Context, userID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Task{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *taskRepo) SetFinished(ctx context.Context, userID, id uuid.UUID, finishedAt *time.Time) error {
	res := r.db.WithContext(ctx).
		Model(&model.Task{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("finished_at", finishedAt)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}