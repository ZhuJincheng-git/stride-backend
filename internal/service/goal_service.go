package service

import (
	"context"
	"strings"
	"time"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/ZhuJincheng-git/stride-backend/internal/repository"
	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
	"github.com/google/uuid"
)

type GoalInput struct {
	Title             *string
	Description       *string
	ExpectedStartTime *time.Time
	ExpectedEndTime   *time.Time
	ParentGoalID      *uuid.UUID
}

type GoalService struct {
	goals repository.GoalRepository
}

func NewGoalService(goals repository.GoalRepository) *GoalService {
	return &GoalService{goals: goals}
}

// Create persists a new goal owned by `userID`.
func (s *GoalService) Create(ctx context.Context, userID uuid.UUID, in GoalInput) (*model.Goal, error) {
	if in.Title == nil || strings.TrimSpace(*in.Title) == "" {
		return nil, apperror.New(apperror.CodeInvalidArgument, "title is required")
	}
	g := &model.Goal{
		UserID:            userID,
		Title:             strings.TrimSpace(*in.Title),
		ExpectedStartTime: in.ExpectedStartTime,
		ExpectedEndTime:   in.ExpectedEndTime,
		ParentGoalID:      in.ParentGoalID,
	}
	if in.Description != nil {
		g.Description = *in.Description
	}
	if err := s.goals.Create(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

// Update applies a partial update. Only non-nil fields in `in` are written.
func (s *GoalService) Update(ctx context.Context, userID, id uuid.UUID, in GoalInput) (*model.Goal, error) {
	g, err := s.goals.GetByID(ctx, userID, id, false)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "goal not found")
		}
		return nil, err
	}
	if in.Title != nil {
		if strings.TrimSpace(*in.Title) == "" {
			return nil, apperror.New(apperror.CodeInvalidArgument, "title cannot be empty")
		}
		g.Title = strings.TrimSpace(*in.Title)
	}
	if in.Description != nil {
		g.Description = *in.Description
	}
	if in.ExpectedStartTime != nil {
		g.ExpectedStartTime = in.ExpectedStartTime
	}
	if in.ExpectedEndTime != nil {
		g.ExpectedEndTime = in.ExpectedEndTime
	}
	if in.ParentGoalID != nil {
		if *in.ParentGoalID == id {
			return nil, apperror.New(apperror.CodeInvalidArgument, "a goal cannot be its own parent")
		}
		g.ParentGoalID = in.ParentGoalID
	}
	if err := s.goals.Update(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

// Get loads a single goal. `includeDeleted=true` is used by the "view trash"
// screen and by the restore handler.
func (s *GoalService) Get(ctx context.Context, userID, id uuid.UUID, includeDeleted bool) (*model.Goal, error) {
	g, err := s.goals.GetByID(ctx, userID, id, includeDeleted)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "goal not found")
		}
		return nil, err
	}
	return g, nil
}

// ListFilter is the public-API counterpart of repository.GoalFilter.
type ListFilter struct {
	Completed      *bool
	Deleted        bool
	IncludeDeleted bool
}

// List returns goals matching the filter, owned by `userID`.
func (s *GoalService) List(ctx context.Context, userID uuid.UUID, f ListFilter) ([]model.Goal, error) {
	return s.goals.List(ctx, repository.GoalFilter{
		UserID:         userID,
		OnlyCompleted:  f.Completed,
		OnlyDeleted:    f.Deleted,
		IncludeDeleted: f.IncludeDeleted,
	})
}

// SoftDelete moves the goal to the trash.
func (s *GoalService) SoftDelete(ctx context.Context, userID, id uuid.UUID) error {
	if err := s.goals.SoftDelete(ctx, userID, id); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "goal not found")
		}
		return err
	}
	return nil
}

// Restore brings a soft-deleted goal back.
func (s *GoalService) Restore(ctx context.Context, userID, id uuid.UUID) (*model.Goal, error) {
	if err := s.goals.Restore(ctx, userID, id); err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "no soft-deleted goal with that id")
		}
		return nil, err
	}
	return s.goals.GetByID(ctx, userID, id, false)
}

// HardDelete removes the goal permanently.
func (s *GoalService) HardDelete(ctx context.Context, userID, id uuid.UUID) error {
	if err := s.goals.HardDelete(ctx, userID, id); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "no soft-deleted goal with that id")
		}
		return err
	}
	return nil
}

// Complete marks the goal as finished.
func (s *GoalService) Complete(ctx context.Context, userID, id uuid.UUID) (*model.Goal, error) {
	g, err := s.goals.GetByID(ctx, userID, id, false)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "goal not found")
		}
		return nil, err
	}
	if g.FinishedAt != nil {
		return g, nil
	}
	now := time.Now()
	if err := s.goals.SetFinished(ctx, userID, id, &now); err != nil {
		return nil, err
	}
	g.FinishedAt = &now
	return g, nil
}

// Uncomplete reverts a goal back to "in progress".
func (s *GoalService) Uncomplete(ctx context.Context, userID, id uuid.UUID) (*model.Goal, error) {
	g, err := s.goals.GetByID(ctx, userID, id, false)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "goal not found")
		}
		return nil, err
	}
	if g.FinishedAt == nil {
		return g, nil
	}
	if err := s.goals.SetFinished(ctx, userID, id, nil); err != nil {
		return nil, err
	}
	g.FinishedAt = nil
	return g, nil
}
