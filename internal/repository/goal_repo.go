package repository

import (
	"context"
	"time"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GoalFilter struct {
	UserID         uuid.UUID
	OnlyCompleted  *bool
	OnlyDeleted    bool
	includeDeleted bool
	Limit          int
	Offset         int
}

type GoalRepository interface {
	Create(ctx context.Context, g *model.Goal) error
	Update(ctx context.Context, g *model.Goal) error
	GetByID(ctx context.Context, userID, id uuid.UUID, includeDeleted bool) (*model.Goal, error)
	List(ctx context.Context, f GoalFilter) ([]model.Goal, error)
	SoftDelete(ctx context.Context, userID, id uuid.UUID) error
	Restore(ctx context.Context, userID, id uuid.UUID) error
	HardDelete(ctx context.Context, userID, id uuid.UUID) error
	SetFinished(ctx context.Context, userID, id uuid.UUID, finishedAt *time.Time) error
}

type goalRepo struct{ db *gorm.DB }

func NewGoalRepository(db *gorm.DB) GoalRepository { return &goalRepo{db: db} }

func (r *goalRepo) Create(ctx context.Context, g *model.Goal) error {
	return r.db.WithContext(ctx).Create(g).Error
}

func (r *goalRepo) Update(ctx context.Context, g *model.Goal) error {
	return r.db.WithContext(ctx).Save(g).Error
}

func (r *goalRepo) GetByID(ctx context.Context, userID, id uuid.UUID, includeDeleted bool) (*model.Goal, error) {
	q := r.db.WithContext(ctx)
	if includeDeleted {
		q = q.Unscoped()
	}
	var g model.Goal
	if err := q.First(&g, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *goalRepo) List(ctx context.Context, f GoalFilter) ([]model.Goal, error) {
	q := r.db.WithContext(ctx).Model(&model.Goal{})
	if f.OnlyDeleted {
		q = q.Unscoped().Where("deleted_at IS NOT NULL")
	} else if f.includeDeleted {
		q = q.Unscoped()
	}
	q = q.Where("user_id = ?", f.UserID)
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

	var out []model.Goal
	if err := q.Order("created_at DESC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (r *goalRepo) SoftDelete(ctx context.Context, userID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Goal{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *goalRepo) Restore(ctx context.Context, userID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.Goal{}).
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

func (r *goalRepo) HardDelete(ctx context.Context, userID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Goal{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *goalRepo) SetFinished(ctx context.Context, userID, id uuid.UUID, finishedAt *time.Time) error {
	res := r.db.WithContext(ctx).
		Model(&model.Goal{}).
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
