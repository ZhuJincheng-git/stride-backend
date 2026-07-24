package response_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
	"github.com/ZhuJincheng-git/stride-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func init() { gin.SetMode(gin.TestMode) }

func dispatch(fn func(c *gin.Context)) *httptest.ResponseRecorder {
	r := gin.New()
	r.GET("/", fn)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) response.Envelope {
	t.Helper()
	var env response.Envelope
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env), rec.Body.String())
	return env
}

func TestOKWritesEnvelopeAnd200(t *testing.T) {
	rec := dispatch(func(c *gin.Context) {
		response.OK(c, gin.H{"id": 7})
	})
	require.Equal(t, http.StatusOK, rec.Code)
	env := decode(t, rec)
	require.Equal(t, "ok", env.Code)
	require.Equal(t, "ok", env.Message)
	require.NotNil(t, env.Data)
}

func TestCreatedWritesEnvelopeAnd201(t *testing.T) {
	rec := dispatch(func(c *gin.Context) {
		response.Created(c, gin.H{"id": 1})
	})
	require.Equal(t, http.StatusCreated, rec.Code)
	env := decode(t, rec)
	require.Equal(t, "ok", env.Code)
	require.Equal(t, "created", env.Message)
}

func TestNoContentWrites204AndEmptyBody(t *testing.T) {
	rec := dispatch(func(c *gin.Context) {
		response.NoContent(c)
	})
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Empty(t, rec.Body.String())
}

func TestErrorMapsAppErrorToStatus(t *testing.T) {
	rec := dispatch(func(c *gin.Context) {
		response.Error(c, apperror.New(apperror.CodeNotFound, "user not found"))
	})
	require.Equal(t, http.StatusNotFound, rec.Code)
	env := decode(t, rec)
	require.Equal(t, "not_found", env.Code)
	require.Equal(t, "user not found", env.Message)
}

func TestErrorMapsGormNotFoundTo404(t *testing.T) {
	rec := dispatch(func(c *gin.Context) {
		response.Error(c, gorm.ErrRecordNotFound)
	})
	require.Equal(t, http.StatusNotFound, rec.Code)
	env := decode(t, rec)
	require.Equal(t, "not_found", env.Code)
}

func TestErrorMapsUnknownErrorTo500AndHidesMessage(t *testing.T) {
	rec := dispatch(func(c *gin.Context) {
		response.Error(c, errors.New("secret string"))
	})
	require.Equal(t, http.StatusInternalServerError, rec.Code)
	env := decode(t, rec)
	require.Equal(t, "internal", env.Code)
	require.Equal(t, "internal server error", env.Message)
	require.NotContains(t, rec.Body.String(), "secret", "internal details must not leak to clients")
}

func TestErrorWithNilIsNoop(t *testing.T) {
	rec := dispatch(func(c *gin.Context) {
		response.Error(c, nil)
		c.Status(http.StatusOK)
	})
	require.Equal(t, http.StatusOK, rec.Code)
}
