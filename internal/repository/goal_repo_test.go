package repository_test

import (
	"context"
	"testing"

	"github.com/ZhuJincheng-git/stride-backend/internal/database"
	"github.com/ZhuJincheng-git/stride-backend/internal/model"
	"github.com/ZhuJincheng-git/stride-backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type goalRepoFixture struct {
	repo   repository.GoalRepository
	db     *gorm.DB
	userID uuid.UUID
}

func newGoalRepoFixture(t *testing.T) *goalRepoFixture {
	t.Helper()
	db, err := database.OpenSQLite()
	require.NoError(t, err)

	userID := uuid.New()
	require.NoError(t, db.Create(&model.User{
		BaseEntity:   model.BaseEntity{ID: userID},
		Username:     "r",
		Email:        "r@example.com",
		PasswordHash: "x",
	}).Error)

	return &goalRepoFixture{
		repo:   repository.NewGoalRepository(db),
		db:     db,
		userID: userID,
	}
}

func TestGoalRepoCreateAssignsUUIDOnInsert(t *testing.T) {
	f := newGoalRepoFixture(t)
	g := &model.Goal{UserID: f.userID, Title: "g"}
	require.NoError(t, f.repo.Create(context.Background(), g))
	require.NotEqual(t, uuid.Nil, g.ID)
}

func TestGoalRepoGetByIDIsScopedToUser(t *testing.T) {
	f := newGoalRepoFixture(t)
	g := &model.Goal{UserID: f.userID, Title: "mine"}
	require.NoError(t, f.repo.Create(context.Background(), g))

	_, err := f.repo.GetByID(context.Background(), uuid.New(), g.ID, false)
	require.True(t, repository.IsNotFound(err), "another user must see NotFound, got %v", err)
}

func TestGoalRepoRestoreOnlyAffectsDeletedRows(t *testing.T) {
	f := newGoalRepoFixture(t)
	g := &model.Goal{UserID: f.userID, Title: "g"}
	require.NoError(t, f.repo.Create(context.Background(), g))

	err := f.repo.Restore(context.Background(), f.userID, g.ID)
	require.True(t, repository.IsNotFound(err), "restore on a live row must return NotFound, got %v", err)
}

func TestGoalRepoHardDeleteRemovesEvenSoftDeletedRows(t *testing.T) {
	f := newGoalRepoFixture(t)
	g := &model.Goal{UserID: f.userID, Title: "g"}
	require.NoError(t, f.repo.Create(context.Background(), g))
	require.NoError(t, f.repo.SoftDelete(context.Background(), f.userID, g.ID))

	require.NoError(t, f.repo.HardDelete(context.Background(), f.userID, g.ID))
	
	_, err := f.repo.GetByID(context.Background(), f.userID, g.ID, true)
	require.True(t, repository.IsNotFound(err))
}

func TestGoalRepoListOnlyDeletedReturnsTrash(t *testing.T) {
	f := newGoalRepoFixture(t)
	live := &model.Goal{UserID: f.userID, Title: "live"}
	dead := &model.Goal{UserID: f.userID, Title: "dead"}
	require.NoError(t, f.repo.Create(context.Background(), live))
	require.NoError(t, f.repo.Create(context.Background(), dead))
	require.NoError(t, f.repo.SoftDelete(context.Background(), f.userID, dead.ID))

	out, err := f.repo.List(context.Background(), repository.GoalFilter{
		UserID: f.userID,
		OnlyDeleted: true,
	})
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, dead.ID, out[0].ID)
}