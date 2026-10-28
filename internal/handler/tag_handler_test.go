package handler_test

import (
	"net/http"
	"testing"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/ZhuJincheng-git/stride-backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

// createTask is a tiny helper.
func createTask(t *testing.T, h *testutil.Harness, token, title string) model.Task {
	t.Helper()
	rec := h.Do(t, http.MethodPost, "/api/v1/tasks", token, map[string]any{"title": title})
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	return testutil.DecodeData[model.Task](t, rec)
}

func TestTaskTagsAttachAndDetach(t *testing.T) {
	h := testutil.New(t)
	token := h.Register(t, "tag-user", "tag@example.com", "super-secret")

	task := createTask(t, h, token, "Tagged task")

	// create a tag
	rec := h.Do(t, http.MethodPost, "/api/v1/task-tags", token, map[string]any{"name": "deep-work"})
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	tag := testutil.DecodeData[model.TaskTag](t, rec)

	// attach
	rec = h.Do(t, http.MethodPost, "/api/v1/tasks/"+task.ID.String()+"/tags", token, map[string]any{
		"tag_ids": []string{tag.ID.String()},
	})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	// list
	rec = h.Do(t, http.MethodGet, "/api/v1/tasks/"+task.ID.String()+"/tags", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, testutil.DecodeData[[]model.TaskTag](t, rec), 1)

	// idempotent attach
	rec = h.Do(t, http.MethodPost, "/api/v1/tasks/"+task.ID.String()+"/tags", token, map[string]any{
		"tag_ids": []string{tag.ID.String()},
	})
	require.Equal(t, http.StatusOK, rec.Code)

	// detach
	rec = h.Do(t, http.MethodDelete, "/api/v1/tasks/"+task.ID.String()+"/tags", token, map[string]any{
		"tag_ids": []string{tag.ID.String()},
	})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	// list empty after detach
	rec = h.Do(t, http.MethodGet, "/api/v1/tasks/"+task.ID.String()+"/tags", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, testutil.DecodeData[[]model.TaskTag](t, rec))

	// permanent delete
	rec = h.Do(t, http.MethodDelete, "/api/v1/task-tags/"+tag.ID.String(), token, nil)
	require.Equal(t, http.StatusNoContent, rec.Code)
}
