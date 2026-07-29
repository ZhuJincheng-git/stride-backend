package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"

	"github.com/ZhuJincheng-git/stride-backend/internal/app"
	"github.com/ZhuJincheng-git/stride-backend/internal/config"
	"github.com/ZhuJincheng-git/stride-backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type Harness struct {
	App *app.App
	DB *gorm.DB
}

func New(t *testing.T) *Harness {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := database.OpenSQLite()
	require.NoError(t, err)
	cfg := &config.Config{
		AppEnv: config.Test,
		JWTSecret: "test-secret",
		JWTExpiresHours: 1,
	}
	return &Harness{App: app.New(cfg, db), DB: db}
}

func (h *Harness) Do(t *testing.T, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.App.Engine.ServeHTTP(rec, req)
	return rec
}

func (h *Harness) Register(t *testing.T, username, email, password string) string {
	t.Helper()
	rec := h.Do(t, http.MethodPost, "/api/v1/auth/register", "", map[string]any{
		"username": username,
		"email": email,
		"password": password,
	})
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	var env struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&env))
	require.NotEmpty(t, env.Data.Token)
	return env.Data.Token
}

func DecodeData[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var env struct {
		Data T `json:"data"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&env), rec.Body.String())
	return env.Data
}