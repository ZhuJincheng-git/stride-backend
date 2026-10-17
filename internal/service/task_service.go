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

type TaskInput struct {
	Title             *string
	Description       *string
	ExpectedStartTime *time.Time
	ExpectedEndTime   *time.Time
	GoalID            *uuid.UUID
	ParentTaskID      *uuid.UUID
	ClearGoal         bool
	ClearParent       bool
}

type TaskListFilter struct {
	GoalID         *uuid.UUID
	NoGoal         bool
	Completed      *bool
	Deleted        bool
	IncludeDeleted bool
}

type TaskService struct {
	tasks repository.TaskRepository
}

func NewTaskService(tasks repository.TaskRepository) *TaskService {
	return &TaskService{tasks: tasks}
}

func (s *TaskService) Create(ctx context.Context, userID uuid.UUID, in TaskInput) (*model.Task, error) {
	if in.Title == nil || strings.TrimSpace(*in.Title) == "" {
		return nil, apperror.New(apperror.CodeInvalidArgument, "title is required")
	}
	t := &model.Task{
		UserID:            userID,
		Title:             strings.TrimSpace(*in.Title),
		ExpectedStartTime: in.ExpectedStartTime,
		ExpectedEndTime:   in.ExpectedEndTime,
		GoalID:            in.GoalID,
		ParentTaskID:      in.ParentTaskID,
	}
	if in.Description != nil {
		t.Description = *in.Description
	}
	if err := s.tasks.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TaskService) Update(ctx context.Context, userID, id uuid.UUID, in TaskInput) (*model.Task, error) {
	t, err := s.tasks.GetByID(ctx, userID, id, false)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "task not found")
		}
		return nil, err
	}
	if in.Title != nil {
		if strings.TrimSpace(*in.Title) == "" {
			return nil, apperror.New(apperror.CodeInvalidArgument, "title cannot be empty")
		}
		t.Title = strings.TrimSpace(*in.Title)
	}
	if in.Description != nil {
		t.Description = *in.Description
	}
	if in.ExpectedStartTime != nil {
		t.ExpectedStartTime = in.ExpectedStartTime
	}
	if in.ExpectedEndTime != nil {
		t.ExpectedEndTime = in.ExpectedEndTime
	}
	if in.ClearGoal {
		t.GoalID = nil
	} else if in.GoalID != nil {
		t.GoalID = in.GoalID
	}
	if in.ClearParent {
		t.ParentTaskID = nil
	} else if in.ParentTaskID != nil {
		if *in.ParentTaskID == id {
			return nil, apperror.New(apperror.CodeInvalidArgument, "a task cannot be its own parent")
		}
		t.ParentTaskID = in.ParentTaskID
	}
	if err := s.tasks.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TaskService) Get(ctx context.Context, userID, id uuid.UUID, includeDeleted bool) (*model.Task, error) {
	t, err := s.tasks.GetByID(ctx, userID, id, includeDeleted)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "task not found")
		}
		return nil, err
	}
	return t, nil
}

func (s *TaskService) List(ctx context.Context, userID uuid.UUID, f TaskListFilter) ([]model.Task, error) {
	repoFilter := repository.TaskFilter{
		UserID: userID,
		GoalID: f.GoalID,
		OnlyCompleted: f.Completed,
		OnlyDeleted: f.Deleted,
		IncludeDeleted: f.IncludeDeleted,
	}
	if f.NoGoal {
		all, err := s.tasks.List(ctx, repoFilter)
		if err != nil {
			return nil, err
		}
		out := all[:0]
		for _, t := range all {
			if t.GoalID == nil {
				out = append(out, t)
			}
		}
		return out, nil
	}
	return s.tasks.List(ctx, repoFilter)
}

func (s *TaskService) SoftDelete(ctx context.Context, userID, id uuid.UUID) error {
	if err := s.tasks.SoftDelete(ctx, userID, id); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "task not found")
		}
		return err
	}
	return nil
}

func (s *TaskService) Restore(ctx context.Context, userID, id uuid.UUID) (*model.Task, error) {
	if err := s.tasks.Restore(ctx, userID, id); err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "no soft-deleted task with that id")
		}
		return nil, err
	}
	return s.tasks.GetByID(ctx, userID, id, false)
}

func (s *TaskService) HardDelete(ctx context.Context, userID, id uuid.UUID) error {
	if err := s.tasks.HardDelete(ctx, userID, id); err != nil {
		if repository.IsNotFound(err) {
			return apperror.New(apperror.CodeNotFound, "task not found")
		}
		return err
	}
	return nil
}

func (s *TaskService) Complete(ctx context.Context, userID, id uuid.UUID) (*model.Task, error) {
	t, err := s.tasks.GetByID(ctx, userID, id, false)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "task not found")
		}
		return nil, err
	}
	now := time.Now()
	if err := s.tasks.SetFinished(ctx, userID, id, &now); err != nil {
		return nil, err
	}
	t.FinishedAt = &now
	return t, nil
}

func (s *TaskService) Uncomplete(ctx context.Context, userID, id uuid.UUID) (*model.Task, error) {
	t, err := s.tasks.GetByID(ctx, userID, id, false)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, apperror.New(apperror.CodeNotFound, "task not found")
		}
		return nil, err
	}
	if err := s.tasks.SetFinished(ctx, userID, id, nil); err != nil {
		return nil, err
	}
	t.FinishedAt = nil
	return t, nil
}