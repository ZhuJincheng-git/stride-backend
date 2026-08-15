package handler

import (
	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// bindingError wraps a JSON-binding error into domain error.
func bindingError(err error) error {
	if err == nil {
		return nil
	}
	return apperror.New(apperror.CodeInvalidArgument, err.Error())
}

// parseUUIDParam reads a URL path parameter as a UUID.
func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, error) {
	raw := c.Param(name)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, apperror.Newf(apperror.CodeInvalidArgument, "invalid %s: %q", name, raw)
	}
	return id, nil
}