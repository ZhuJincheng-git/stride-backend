package handler

import (
	"strconv"

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

// parseBoolQuery handles the `?completed=true|false` like params.
func parseBoolQuery(c *gin.Context, name string) (*bool, error) {
	raw := c.Query(name)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, apperror.Newf(apperror.CodeInvalidArgument, "invalid %s: %q", name, raw)
	}
	return &v, nil
}