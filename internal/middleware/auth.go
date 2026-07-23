package middleware

import (
	"net/http"
	"strings"

	"github.com/ZhuJincheng-git/stride-backend/pkg/jwt"
	"github.com/ZhuJincheng-git/stride-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	CtxKeyUserID = "auth.user_id"
)

func AuthRequired(tokens *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			abort(c, http.StatusUnauthorized, "missing Authorization header")
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			abort(c, http.StatusUnauthorized, "invalid Authorization header")
			return
		}
		claims, err := tokens.Parse(parts[1])
		if err != nil {
			abort(c, http.StatusUnauthorized, "invalid or expired token")
			return 
		}
		c.Set(CtxKeyUserID, claims.UserID)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(CtxKeyUserID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

func abort(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, response.Envelope{
		Code: "unauthenticated",
		Message: msg,
	})
}