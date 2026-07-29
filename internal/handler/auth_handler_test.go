package handler_test

import (
	"net/http"
	"testing"

	"github.com/ZhuJincheng-git/stride-backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestAuthRegisterAndLogin(t *testing.T) {
	h := testutil.New(t)

	// --- register ---
	rec := h.Do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]any{
		"username": "alice",
		"email": "alice@example.com",
		"password": "correct-horse",
	})
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	// --- duplicate username fails ---
	rec = h.Do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]any{
		"username": "alice",
		"email": "alice2@example.com",
		"password": "correct-horse",
	})
	require.NotEqual(t, http.StatusCreated, rec.Code, rec.Body.String())

	// --- login with username ---
	rec = h.Do(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"identifier": "alice",
		"password": "correct-horse",
	})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	// --- wrong password ---
	rec = h.Do(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"identifier": "alice",
		"password": "wrong-password",
	})
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())
}

func TestAuthMeRequiresToken(t *testing.T) {
	h := testutil.New(t)
	rec := h.Do(t, http.MethodGet, "/api/v1/auth/me", "", nil)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMeReturnsCurrentUser(t *testing.T) {
	h := testutil.New(t)
	token := h.Register(t, "bob", "bob@example.com", "password")
	rec := h.Do(t, http.MethodGet, "/api/v1/auth/me", token, nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
}