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

