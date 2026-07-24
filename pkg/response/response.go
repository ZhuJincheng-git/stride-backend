package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ZhuJincheng-git/stride-backend/pkg/apperror"
)

type Envelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Code: "ok", Message: "ok", Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{Code: "ok", Message: "created", Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, err error) {
	if err == nil {
		return
	}

	if ae, ok := apperror.AsAppError(err); ok {
		c.AbortWithStatusJSON(ae.StatusCode(), Envelope{
			Code:    string(ae.Code),
			Message: ae.Message,
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, Envelope{
			Code:    string(apperror.CodeNotFound),
			Message: "resource not found",
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, Envelope{
		Code:    string(apperror.CodeInternal),
		Message: "internal server error",
	})
	return
}
