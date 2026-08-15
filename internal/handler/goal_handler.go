package handler

import (
	"time"

	"github.com/ZhuJincheng-git/stride-backend/internal/middleware"
	"github.com/ZhuJincheng-git/stride-backend/internal/service"
	"github.com/ZhuJincheng-git/stride-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GoalHandler struct {
	svc *service.GoalService
}

func NewGoalHandler(s *service.GoalService) *GoalHandler { return &GoalHandler{svc: s} }

type goalUpsertRequest struct {
	Title             *string    `join:"title"`
	Description       *string    `join:"description"`
	ExpectedStartTime *time.Time `join:"expected_start_time"`
	ExpectedEndTime   *time.Time `join:"expected_end_time"`
	ParentGoalID      *uuid.UUID `join:"parent_goal_id"`
}

// Create handles POST /goals.
func (h *GoalHandler) Create(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	var req goalUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	g, err := h.svc.Create(c.Request.Context(), userID, service.GoalInput{
		Title: req.Title,
		Description: req.Description,
		ExpectedStartTime: req.ExpectedStartTime,
		ExpectedEndTime: req.ExpectedEndTime,
		ParentGoalID: req.ParentGoalID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, g)
}

// Update handles PATCH /goals/:id
func (h *GoalHandler) Update(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req goalUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	g, err := h.svc.Update(c.Request.Context(), userID, id, service.GoalInput{
		Title: req.Title,
		Description: req.Description,
		ExpectedStartTime: req.ExpectedStartTime,
		ExpectedEndTime: req.ExpectedEndTime,
		ParentGoalID: req.ParentGoalID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, g)
}

