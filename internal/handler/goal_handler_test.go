package handler_test

import (
	"net/http"
	"testing"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/ZhuJincheng-git/stride-backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestGoalLifecycle(t *testing.T) {
	h := testutil.New(t)
	token := h.Register(t, "Ann", "ann@example.com", "test-password")

	// create
	rec := h.Do(t, http.MethodPost, "/api/v1/goals", token, map[string]any{"title": "test"})
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	created := testutil.DecodeData[model.Goal](t, rec)
	require.Equal(t, "test", created.Title)
	require.NotEmpty(t, created.ID)

	// update
	rec = h.Do(t, http.MethodPatch, "/api/v1/goals/"+created.ID.String(), token, map[string]any{
		"description": "test-description",
	})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	updated := testutil.DecodeData[model.Goal](t, rec)
	require.Equal(t, "test-description", updated.Description)

	// complete
	rec = h.Do(t, http.MethodPost, "/api/v1/goals/"+created.ID.String()+"/complete", token, nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	complete := testutil.DecodeData[model.Goal](t, rec)
	require.NotNil(t, complete.FinishedAt)

	// list completed
	rec = h.Do(t, http.MethodGet, "/api/v1/goals?completed=true", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	list := testutil.DecodeData[[]model.Goal](t, rec)
	require.Len(t, list, 1)

	// list non-completed
	rec = h.Do(t, http.MethodGet, "/api/v1/goals?completed=false", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, testutil.DecodeData[[]model.Goal](t, rec))

	// uncomplete
	rec = h.Do(t, http.MethodPost, "/api/v1/goals/"+created.ID.String()+"/uncomplete", token, nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body)
	uncomplete := testutil.DecodeData[model.Goal](t, rec)
	require.Nil(t, uncomplete.FinishedAt)

	// soft delete
	rec = h.Do(t, http.MethodDelete, "/api/v1/goals/"+created.ID.String(), token, nil)
	require.Equal(t, http.StatusNoContent, rec.Code)

	// trash list
	rec = h.Do(t, http.MethodGet, "/api/v1/goals?deleted=true", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, testutil.DecodeData[[]model.Goal](t, rec), 1)

	// live list
	rec = h.Do(t, http.MethodGet, "/api/v1/goals", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, testutil.DecodeData[[]model.Goal](t, rec))

	// restore
	rec = h.Do(t, http.MethodPost, "/api/v1/goals/"+created.ID.String()+"/restore", token, nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	// permanent delete
	rec = h.Do(t, http.MethodDelete, "/api/v1/goals/"+created.ID.String()+"/permanent", token, nil)
	require.Equal(t, http.StatusNoContent, rec.Code)

	// list check
	rec = h.Do(t, http.MethodGet, "/api/v1/goals?include_deleted=true", token, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, testutil.DecodeData[[]model.Goal](t, rec))
}

func TestGoalCannotBeOwnParent(t *testing.T) {
	h := testutil.New(t)
	token := h.Register(t, "p-user", "p@example.com", "super-secret")
	rec := h.Do(t, http.MethodPost, "/api/v1/goals", token, map[string]any{"title": "test"})
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	g := testutil.DecodeData[model.Goal](t, rec)

	rec = h.Do(t, http.MethodPatch, "/api/v1/goals/"+g.ID.String(), token, map[string]any{
		"parent_goal_id": g.ID,
	})
	require.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
}

