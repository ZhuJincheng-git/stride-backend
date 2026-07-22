package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ZhuJincheng-git/stride-backend/internal/database"
	"github.com/ZhuJincheng-git/stride-backend/internal/repository"
	"github.com/ZhuJincheng-git/stride-backend/internal/service"
	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
	"github.com/ZhuJincheng-git/stride-backend/pkg/jwt"
)

func newAuthSvc(t *testing.T) (*service.AuthService, *jwt.Manager) {
	t.Helper()
	db, err := database.OpenSQLite()
	require.NoError(t, err)
	tokens := jwt.New("auth-test-secret", time.Hour, "stride-tests")
	svc := service.NewAuthService(
		repository.NewUserRepository(db),
		tokens,
	)
	return svc, tokens
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	cases := []struct {
		name string
		in   service.RegisterInput
	}{
		{"missing username", service.RegisterInput{Email: "a@a.com", Password: "12345678"}},
		{"missing email", service.RegisterInput{Username: "a", Password: "12345678"}},
		{"missing password", service.RegisterInput{Username: "a", Email: "a@a.com"}},
		{"short password", service.RegisterInput{Username: "a", Email: "a@a.com", Password: "1234567"}},
		{"whitespace username", service.RegisterInput{Username: "   ", Email: "a@a.com", Password: "12345678"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, _ := newAuthSvc(t)
			_, err := svc.Register(context.Background(), tc.in)
			ae, ok := apperror.AsAppError(err)
			require.True(t, ok, "expected typed error, got %v", err)
			require.Equal(t, apperror.CodeInvalidArgument, ae.Code)
		})
	}
}

func TestRegisterCreatesUserAndReturnsParsableToken(t *testing.T) {
	svc, tokens := newAuthSvc(t)

	out, err := svc.Register(context.Background(), service.RegisterInput{
		Username: "alice",
		Email: "Alice@Example.com",
		Password: "hunter22hunter22",
	})
	require.NoError(t, err)
	require.NotEmpty(t, out.Token)
	require.Equal(t, "alice", out.User.Username)
	require.Equal(t, "alice@example.com", out.User.Email, "email must be lowercased")
	require.NotEqual(t, "hunter22hunter22", out.User.PasswordHash, "password must be hashed, not stored as plaintext")

	claims, err := tokens.Parse(out.Token)
	require.NoError(t, err)
	require.Equal(t, out.User.ID, claims.UserID)
}

func TestRegisterRejectsDuplicateUsername(t *testing.T) {
	svc, _ := newAuthSvc(t)
	_, err := svc.Register(context.Background(), service.RegisterInput{
		Username: "alice", Email: "a@a.com", Password: "hunter22hunter22",
	})
	require.NoError(t, err)

	_, err = svc.Register(context.Background(), service.RegisterInput{
		Username: "alice", Email: "b@b.com", Password: "hunter22hunter22",
	})
	require.Error(t, err)
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	svc, _ := newAuthSvc(t)
	_, err := svc.Register(context.Background(), service.RegisterInput{
		Username: "alice", Email: "a@a.com", Password: "hunter22hunter22",
	})
	require.NoError(t, err)

	_, err = svc.Register(context.Background(), service.RegisterInput{
		Username: "bob", Email: "A@a.com", Password: "hunter22hunter22",
	})
	require.Error(t, err)
}

func TestLoginAcceptsUsernameOrEmail(t *testing.T) {
	svc, _ := newAuthSvc(t)
	_, err := svc.Register(context.Background(), service.RegisterInput{
		Username: "alice", Email: "alice@example.com", Password: "hunter22hunter22",
	})
	require.NoError(t, err)

	for _, identifier := range []string{"alice", "alice@example.com", "Alice@Example.com"} {
		out, err := svc.Login(context.Background(), service.LoginInput{
			Identifier: identifier, Password: "hunter22hunter22",
		})
		require.NoErrorf(t, err, "login with %q failed", identifier)
		require.NotEmpty(t, out.Token)
	}
}

func TestLoginReturnsUnauthenticatedOnUnknownUser(t *testing.T) {
	svc, _ := newAuthSvc(t)
	_, err := svc.Login(context.Background(), service.LoginInput{
		Identifier: "ghost", Password: "whatever1234",
	})
	ae, ok := apperror.AsAppError(err)
	require.True(t, ok)
	require.Equal(t, apperror.CodeUnauthenticated, ae.Code)
}

func TestCurrentUserReturnsUserForValidID(t *testing.T) {
	svc, _ := newAuthSvc(t)
	reg, err := svc.Register(context.Background(), service.RegisterInput{
		Username: "alice", Email: "alice@example.com", Password: "hunter22hunter22",
	})
	require.NoError(t, err)

	got, err := svc.CurrentUser(context.Background(), reg.User.ID)
	require.NoError(t, err)
	require.Equal(t, reg.User.ID, got.ID)
}

func TestCurrentUserReturnsUnauthenticatedForUnknownID(t *testing.T) {
	svc, _ := newAuthSvc(t)
	_, err := svc.CurrentUser(context.Background(), uuid.New())
	ae, ok := apperror.AsAppError(err)
	require.True(t, ok)
	require.Equal(t, apperror.CodeUnauthenticated, ae.Code)
}