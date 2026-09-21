package service

import (
	"context"
	"strings"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/ZhuJincheng-git/stride-backend/internal/repository"
	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
	"github.com/google/uuid"
)

// TagService handles user-defined goal/task labels.
type TagService struct {
	tags repository.TagRepository
}

func NewTagService(tags repository.TagRepository) *TagService { return &TagService{tags: tags} }

// --- goal tag operations ---

func (s *TagService) CreateGoalTag(ctx context.Context, userID uuid.UUID, name string) (*model.GoalTag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperror.New(apperror.CodeInvalidArgument, "name is required")
	}
	t := &model.GoalTag{UserID: userID, Name: name}
	if err := s.tags.CreateGoalTag(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TagService) ListGoalTags(ctx context.Context, userID uuid.UUID) ([]model.GoalTag, error) {
	return s.tags.ListGoalTags(ctx, userID)
}

func (s *TagService) DeleteGoalTag(ctx context.Context, userID, id uuid.UUID) error {
	if err := s.tags.HardDeleteGoalTag(ctx, userID, id); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "tag not found")
		}
		return err
	}
	return nil
}

func (s *TagService) AttachGoalTags(ctx context.Context, userID, goalID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return apperror.New(apperror.CodeInvalidArgument, "at least one tag id required")
	}
	if err := s.tags.AttachGoalTags(ctx, userID, goalID, ids); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "goal not found")
		}
		return err
	}
	return nil
}

func (s *TagService) DetachGoalTags(ctx context.Context, userID, goalID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return apperror.New(apperror.CodeInvalidArgument, "at least one tag id required")
	}
	if err := s.tags.DetachGoalTags(ctx, userID, goalID, ids); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "goal not found")
		}
		return err
	}
	return nil
}

func (s *TagService) LIstGoalTagsForGoal(ctx context.Context, userID, goalID uuid.UUID) ([]model.GoalTag, error) {
	out, err := s.tags.ListTagsForGoal(ctx, userID, goalID)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "goal not found")
		}
		return nil, err
	}
	return out, nil
}

// --- task tag operations ---

func (s *TagService) CreateTaskTag(ctx context.Context, userID uuid.UUID, name string) (*model.TaskTag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apperror.New(apperror.CodeInvalidArgument, "name is required")
	}
	t := &model.TaskTag{UserID: userID, Name: name}
	if err := s.tags.CreateTaskTag(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TagService) ListTaskTags(ctx context.Context, userID uuid.UUID) ([]model.TaskTag, error) {
	return s.tags.ListTaskTags(ctx, userID)
}

func (s *TagService) DeleteTaskTag(ctx context.Context, userID, id uuid.UUID) error {
	if err := s.tags.HardDeleteTaskTag(ctx, userID, id); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "tag not found")
		}
		return err
	}
	return nil
}

func (s *TagService) AttachTaskTags(ctx context.Context, userID, taskID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return apperror.New(apperror.CodeInvalidArgument, "at least one tag id required")
	}
	if err := s.tags.AttachTaskTags(ctx, userID, taskID, ids); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "task or tag(s) not found")
		}
		return err
	}
	return nil
}

func (s *TagService) DetachTaskTags(ctx context.Context, userID, taskID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return apperror.New(apperror.CodeInvalidArgument, "at least one tag id required")
	}
	if err := s.tags.DetachTaskTags(ctx, userID, taskID, ids); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "task not found")
		}
		return err
	}
	return nil
}

func (s *TagService) ListTaskTagsForTask(ctx context.Context, userID, taskID uuid.UUID) ([]model.TaskTag, error) {
	out, err := s.tags.ListTagsForTask(ctx, userID, taskID)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "task not found")
		}
		return nil, err
	}
	return out, nil
}
