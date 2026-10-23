package handler_test

import (
	"net/http"
	"testing"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/ZhuJincheng-git/stride-backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestTaskLifecycle(t *testing.T) {
	h := testutil.New(t)
	token := h.Register(t, "t-user", "t@example.com", "super-secret")

	// goal first
	rec := h.Do(t, http.MethodPost, "/api/v1/goals", token, map[string]any{"title": "Ship MVP"})
	require.Equal(t, http.StatusCreated, rec.Code)
	goal := testutil.DecodeData[model.Goal](t, rec)

	// task with goal
	rec = h.Do(t, http.MethodPost, "/api/v1/tasks", token, map[string]any{
		"title": "Wire DB",
		"goal_id": goal.ID,
	})
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	task := testutil.DecodeData[model.Task](t, rec)
	require.NotNil(t, task.GoalID)
	require.Equal(t, goal.ID, *task.GoalID)

	// task without a goal
	rec = h.Do(t, http.MethodPost, "/api/v1/tasks", token, map[string]any{"title": "Floating task"})
	require.Equal(t, http.StatusCreated, rec.Code)

	// list with goal_id filter returns one
	rec = h.Do(t, http.MethodGet, "/api/v1/tasks?goal_id="+goal.ID.String(), token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, testutil.DecodeData[[]model.Task](t, rec), 1)

	// list no_goal returns the orphan task
	rec = h.Do(t, http.MethodGet, "/api/v1/tasks?no_goal=true", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, testutil.DecodeData[[]model.Task](t, rec), 1)

	// clear_goal updates a task to have no goal
	rec = h.Do(t, http.MethodPatch, "/api/v1/tasks/"+task.ID.String(), token, map[string]any{
		"clear_goal": true,
	})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	cleared := testutil.DecodeData[model.Task](t, rec)
	require.Nil(t, cleared.GoalID)

	// soft delete + restore + permanent delete
	rec = h.Do(t, http.MethodDelete, "/api/v1/tasks/"+task.ID.String(), token, nil)
	require.Equal(t, http.StatusNoContent, rec.Code)
	rec = h.Do(t, http.MethodPost, "/api/v1/tasks/"+task.ID.String()+"/restore", token, nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	rec = h.Do(t, http.MethodDelete, "/api/v1/tasks/"+task.ID.String()+"/permanent", token, nil)
	require.Equal(t, http.StatusNoContent, rec.Code)
}