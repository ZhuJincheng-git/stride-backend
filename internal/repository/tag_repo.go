package repository

import (
	"context"
	"errors"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TagRepository manages user-defined labels for goals and tasks.
type TagRepository interface {
	// --- goal tags ---
	CreateGoalTag(ctx context.Context, t *model.GoalTag) error
	ListGoalTags(ctx context.Context, userID uuid.UUID) ([]model.GoalTag, error)
	HardDeleteGoalTag(ctx context.Context, userID, id uuid.UUID) error
	AttachGoalTags(ctx context.Context, userID, goalID uuid.UUID, tagIDs []uuid.UUID) error
	DetachGoalTags(ctx context.Context, userID, goalID uuid.UUID, tagIDs []uuid.UUID) error
	ListTagsForGoal(ctx context.Context, userID, goalID uuid.UUID) ([]model.GoalTag, error)

	// --- task tags ---
	CreateTaskTag(ctx context.Context, t *model.TaskTag) error
	ListTaskTags(ctx context.Context, userID uuid.UUID) ([]model.TaskTag, error)
	HardDeleteTaskTag(ctx context.Context, userID, id uuid.UUID) error
	AttachTaskTags(ctx context.Context, userID, taskID uuid.UUID, tagIDs []uuid.UUID) error
	DetachTaskTags(ctx context.Context, userID, taskID uuid.UUID, tagIDs []uuid.UUID) error
	ListTagsForTask(ctx context.Context, userID, taskID uuid.UUID) ([]model.TaskTag, error)
}

type tagRepo struct{ db *gorm.DB }

func NewTagRepository(db *gorm.DB) TagRepository { return &tagRepo{db: db} }

// --- goal tag implementations ---

func (r *tagRepo) CreateGoalTag(ctx context.Context, t *model.GoalTag) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *tagRepo) ListGoalTags(ctx context.Context, userID uuid.UUID) ([]model.GoalTag, error) {
	var out []model.GoalTag
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&out).Error
	return out, err
}

func (r *tagRepo) HardDeleteGoalTag(ctx context.Context, userID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&model.GoalTag{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.WithContext(ctx).Where("goal_tag_id = ?", id).Delete(&model.GoalTagRel{}).Error
}

func (r *tagRepo) AttachGoalTags(ctx context.Context, userID, goalID uuid.UUID, tagIDs []uuid.UUID) error {
	if err := r.assertGoalOwned(ctx, userID, goalID); err != nil {
		return err
	}
	if err := r.assertGoalTagsOwned(ctx, userID, tagIDs); err != nil {
		return err
	}
	rels := make([]model.GoalTagRel, 0, len(tagIDs))
	for _, t := range tagIDs {
		rels = append(rels, model.GoalTagRel{GoalID: goalID, GoalTagID: t})
	}
	return r.db.WithContext(ctx).
		Session(&gorm.Session{}).
		Clauses(onConflictDoNothing("goal_id", "goal_tag_id")...).
		Create(&rels).Error
}

func (r *tagRepo) DetachGoalTags(ctx context.Context, userID, goalID uuid.UUID, tagIDs []uuid.UUID) error {
	if err := r.assertGoalOwned(ctx, userID, goalID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Where("goal_id = ? AND goal_tag_id IN ?", goalID, tagIDs).
		Delete(&model.GoalTagRel{}).Error
}

func (r *tagRepo) ListTagsForGoal(ctx context.Context, userID, goalID uuid.UUID) ([]model.GoalTag, error) {
	if err := r.assertGoalOwned(ctx, userID, goalID); err != nil {
		return nil, err
	}
	var out []model.GoalTag
	err := r.db.WithContext(ctx).
		Joins("JOIN goal_tag_rel ON goal_tag_rel.goal_tag_id = goal_tags.id").
		Where("goal_tag_rel.goal_id = ?", goalID).
		Find(&out).Error
	return out, err
}

// --- task tag implementations ---

func (r *tagRepo) CreateTaskTag(ctx context.Context, t *model.TaskTag) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *tagRepo) ListTaskTags(ctx context.Context, userID uuid.UUID) ([]model.TaskTag, error) {
	var out []model.TaskTag
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&out).Error
	return out, err
}

func (r *tagRepo) HardDeleteTaskTag(ctx context.Context, userID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&model.TaskTag{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.WithContext(ctx).Where("task_tag_id = ?", id).Delete(&model.TaskTagRel{}).Error
}

func (r *tagRepo) AttachTaskTags(ctx context.Context, userID, taskID uuid.UUID, tagIDs []uuid.UUID) error {
	if err := r.assertTaskOwned(ctx, userID, taskID); err != nil {
		return err
	}
	if err := r.assertTaskTagsOwned(ctx, userID, tagIDs); err != nil {
		return err
	}
	rels := make([]model.TaskTagRel, 0, len(taskID))
	for _, t := range tagIDs {
		rels = append(rels, model.TaskTagRel{TaskID: taskID, TaskTagID: t})
	}
	return r.db.WithContext(ctx).Session(&gorm.Session{}).Clauses(onConflictDoNothing("task_id", "task_tag_id")...).Create(&rels).Error
}

func (r *tagRepo) DetachTaskTags(ctx context.Context, userID, taskID uuid.UUID, tagIDs []uuid.UUID) error {
	if err := r.assertTaskOwned(ctx, userID, taskID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("task_id = ? AND task_tag_id IN ?", taskID, tagIDs).Delete(&model.TaskTagRel{}).Error
}

func (r *tagRepo) ListTagsForTask(ctx context.Context, userID, taskID uuid.UUID) ([]model.TaskTag, error) {
	if err := r.assertTaskOwned(ctx, userID, taskID); err != nil {
		return nil, err
	}
	var out []model.TaskTag
	err := r.db.WithContext(ctx).
		Joins("JOIN task_tag_rel ON task_tag_rel.task_tag_id = task_tags.id").
		Where("task_tag_rel.task_id = ?", taskID).
		Find(&out).Error
	return out, err
}

// --- ownership guards ---

func (r *tagRepo) assertGoalOwned(ctx context.Context, userID, goalID uuid.UUID) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Goal{}).Where("id = ? AND user_id = ?", goalID, userID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *tagRepo) assertTaskOwned(ctx context.Context, userID, taskID uuid.UUID) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ? AND user_id = ?", taskID, userID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *tagRepo) assertGoalTagsOwned(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no tag ids supplied")
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.GoalTag{}).Where("user_id = ? AND id IN ?", userID, ids).Count(&count).Error; err != nil {
		return err
	}
	if int(count) != len(ids) {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *tagRepo) assertTaskTagsOwned(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no tag ids supplied")
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.GoalTag{}).Where("user_id = ? AND id IN ?", userID, ids).Count(&count).Error; err != nil {
		return err
	}
	if int(count) != len(ids) {
		return gorm.ErrRecordNotFound
	}
	return nil
}
