package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/ZhuJincheng-git/stride-backend/internal/database"
	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/ZhuJincheng-git/stride-backend/internal/repository"
	"github.com/ZhuJincheng-git/stride-backend/internal/service"
	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type goalSvcFixture struct {
	svc    *service.GoalService
	db     *gorm.DB
	userID uuid.UUID
}

func newGoalSvcFixture(t *testing.T) *goalSvcFixture {
	t.Helper()
	db, err := database.OpenSQLite()
	require.NoError(t, err)

	userID := uuid.New()
	require.NoError(t, db.Create(&model.User{
		BaseEntity:   model.BaseEntity{ID: userID},
		Username:     "g",
		Email:        "g@example.com",
		PasswordHash: "x",
	}).Error)

	return &goalSvcFixture{
		svc:    service.NewGoalService(repository.NewGoalRepository(db)),
		db:     db,
		userID: userID,
	}
}

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func TestCreateRejectsEmptyTitle(t *testing.T) {
	f := newGoalSvcFixture(t)
	for _, title := range []string{"", "   "} {
		_, err := f.svc.Create(context.Background(), f.userID, service.GoalInput{Title: strPtr(title)})
		ae, ok := apperror.AsAppError(err)
		require.True(t, ok, "title=%q should produce typed error, got %v", title, err)
		require.Equal(t, apperror.CodeInvalidArgument, ae.Code)
	}
}

func TestCreateTrimsTitleAndPersistsOptionalFields(t *testing.T) {
	f := newGoalSvcFixture(t)
	desc := "ship it"
	start := time.Now().Add(time.Hour).UTC().Truncate(time.Second)

	g, err := f.svc.Create(context.Background(), f.userID, service.GoalInput{
		Title:             strPtr("    ship mvP   "),
		Description:       &desc,
		ExpectedStartTime: &start,
	})
	require.NoError(t, err)
	require.Equal(t, "ship mvP", g.Title)
	require.Equal(t, desc, desc, g.Description)
	require.NotNil(t, g.ExpectedStartTime)
}

func TestUpdateLeavesUnspecifiedFieldAlone(t *testing.T) {
	f := newGoalSvcFixture(t)
	desc := "original"
	g, err := f.svc.Create(context.Background(), f.userID, service.GoalInput{
		Title: strPtr("Title v1"),
		Description: &desc,
	})
	require.NoError(t, err)

	updated, err := f.svc.Update(context.Background(), f.userID, g.ID, service.GoalInput{
		Title: strPtr("Title v2"),
	})
	require.NoError(t, err)
	require.Equal(t, "Title v2", updated.Title)
	require.Equal(t, desc, updated.Description, "Description should not be updated")
}

func TestUpdateRejectsSelfParent(t *testing.T) {
	f := newGoalSvcFixture(t)
	g, err := f.svc.Create(context.Background(), f.userID, service.GoalInput{Title: strPtr("g")})
	require.NoError(t, err)

	_, err = f.svc.Update(context.Background(), f.userID, g.ID, service.GoalInput{ParentGoalID: &g.ID})
	ae, ok := apperror.AsAppError(err)
	require.True(t, ok)
	require.Equal(t, apperror.CodeInvalidArgument, ae.Code)
}

func TestUpdateMissingReturnsNotFound(t *testing.T) {
	f := newGoalSvcFixture(t)
	_, err := f.svc.Update(context.Background(), f.userID, uuid.New(), service.GoalInput{Title: strPtr("x")})
	ae, ok := apperror.AsAppError(err)
	require.True(t, ok)
	require.Equal(t, apperror.CodeNotFound, ae.Code)
}

func TestCompleteIsIdempotent(t *testing.T) {
	f := newGoalSvcFixture(t)
	g, err := f.svc.Create(context.Background(), f.userID, service.GoalInput{Title: strPtr("g")})
	require.NoError(t, err)

	first, err := f.svc.Complete(context.Background(), f.userID, g.ID)
	require.NoError(t, err)
	require.NotNil(t, first.FinishedAt)
	stamp := *first.FinishedAt

	time.Sleep(2 * time.Millisecond)

	second, err := f.svc.Complete(context.Background(), f.userID, g.ID)
	require.NoError(t, err)
	require.NotNil(t, second.FinishedAt)
	require.True(t, second.FinishedAt.Equal(stamp), "complete must be idempotent")
}

func TestUncompleteClearsFinishedAt(t *testing.T) {
	f := newGoalSvcFixture(t)
	g, err := f.svc.Create(context.Background(), f.userID, service.GoalInput{Title: strPtr("g")})
	require.NoError(t, err)
	_, err = f.svc.Complete(context.Background(), f.userID, g.ID)
	require.NoError(t, err)

	out, err := f.svc.Uncomplete(context.Background(), f.userID, g.ID)
	require.NoError(t, err)
	require.Nil(t, out.FinishedAt)
}

func TestSoftDeleteThenRestoreRoundTrip(t *testing.T) {
	f := newGoalSvcFixture(t)
	g, err := f.svc.Create(context.Background(), f.userID, service.GoalInput{Title: strPtr("g")})
	require.NoError(t, err)

	require.NoError(t, f.svc.SoftDelete(context.Background(), f.userID, g.ID))

	_, err = f.svc.Get(context.Background(), f.userID, g.ID, false)
	require.Error(t, err)
	_, err = f.svc.Get(context.Background(), f.userID, g.ID, true)
	require.NoError(t, err)

	restored, err := f.svc.Restore(context.Background(), f.userID, g.ID)
	require.NoError(t, err)
	require.NotNil(t, restored)

	live, err := f.svc.List(context.Background(), f.userID, service.ListFilter{})
	require.NoError(t, err)
	require.Len(t, live, 1)
}

func TestListFiltersByCompleted(t *testing.T) {
	f := newGoalSvcFixture(t)
	a, err := f.svc.Create(context.Background(), f.userID, service.GoalInput{Title: strPtr("a")})
	require.NoError(t, err)
	_, err = f.svc.Create(context.Background(), f.userID, service.GoalInput{Title: strPtr("b")})
	require.NoError(t, err)
	_, err = f.svc.Complete(context.Background(), f.userID, a.ID)
	require.NoError(t, err)

	completed, err := f.svc.List(context.Background(), f.userID, service.ListFilter{Completed: boolPtr(true)})
	require.NoError(t, err)
	require.Len(t, completed, 1)
	require.Equal(t, a.ID, completed[0].ID)

	open, err := f.svc.List(context.Background(), f.userID, service.ListFilter{Completed: boolPtr(false)})
	require.NoError(t, err)
	require.Len(t, open, 1)
	require.NotEqual(t, a.ID, open[0].ID)
}

func TestListIsScopedToUser(t *testing.T) {
	f := newGoalSvcFixture(t)
	_, err := f.svc.Create(context.Background(), f.userID, service.GoalInput{Title: strPtr("mine")})
	require.NoError(t, err)

	other := uuid.New()
	require.NoError(t, f.db.Create(&model.User{
		BaseEntity: model.BaseEntity{ID: other},
		Username: "x",
		Email: "x@example.com",
		PasswordHash: "p",
	}).Error)
	require.NoError(t, f.db.Create(&model.Goal{UserID: other, Title: "theirs"}).Error)

	mine, err := f.svc.List(context.Background(), f.userID, service.ListFilter{})
	require.NoError(t, err)
	require.Len(t, mine, 1)
	require.Equal(t, "mine", mine[0].Title)
}