package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ZhuJincheng-git/stride-backend/internal/service"
	"github.com/ZhuJincheng-git/stride-backend/pkg/response"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{auth: s} }

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// Register handlers POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	res, err := h.auth.Register(c.Request.Context(), service.RegisterInput{
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, res)
}