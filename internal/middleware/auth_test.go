package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ZhuJincheng-git/stride-backend/internal/middleware"
	"github.com/ZhuJincheng-git/stride-backend/pkg/jwt"
)

func init() { gin.SetMode(gin.TestMode) }

func buildAuthRouter(t *testing.T) (*gin.Engine, *jwt.Manager) {
	t.Helper()
	tokens := jwt.New("test-secret", time.Hour, "stride-test")
	r := gin.New()
	r.GET("/me", middleware.AuthRequired(tokens), func(c *gin.Context) {
		uid, ok := middleware.CurrentUserID(c)
		require.True(t, ok)
		c.JSON(http.StatusOK, gin.H{"uid": uid.String()})
	})
	return r, tokens
}

func TestAuthRequiredRejectsMissingHeader(t *testing.T) {
	r, _ := buildAuthRouter(t)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/me", nil))
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthRequiredRejectsMalformedHeader(t *testing.T) {
	cases := map[string]string{
		"only scheme": "Bearer",
		"empty token": "Bearer ",
		"wrong scheme": "Basic abc",
		"no scheme prefix": "abc",
	}
	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			r, _ := buildAuthRouter(t)
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			req.Header.Set("Authorization", header)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			require.Equal(t, http.StatusUnauthorized, rec.Code, "header=%q", header)
		})
	}
}

func TestAuthRequiredRejectsInvalidToken(t *testing.T) {
	r, _ := buildAuthRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer this-is-not-a-jwt")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthRequiredAcceptsValidTokenAndExposesUserID(t *testing.T) {
	r, tokens := buildAuthRouter(t)
	uid := uuid.New()
	tok, err := tokens.Generate(uid)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), uid.String())
}

func TestCurrentUserIDReturnsFalseWithMiddleWareNotRun(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	id, ok := middleware.CurrentUserID(c)
	require.False(t, ok)
	require.Equal(t, uuid.Nil, id)
}
