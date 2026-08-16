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
		Title:             req.Title,
		Description:       req.Description,
		ExpectedStartTime: req.ExpectedStartTime,
		ExpectedEndTime:   req.ExpectedEndTime,
		ParentGoalID:      req.ParentGoalID,
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
		Title:             req.Title,
		Description:       req.Description,
		ExpectedStartTime: req.ExpectedStartTime,
		ExpectedEndTime:   req.ExpectedEndTime,
		ParentGoalID:      req.ParentGoalID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, g)
}

// Get handlers GET /goals/:id.
func (h *GoalHandler) Get(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	includeDeleted, err := parseBoolQuery(c, "include_deleted")
	if err != nil {
		response.Error(c, err)
		return
	}
	g, err := h.svc.Get(c.Request.Context(), userID, id, boolPtrValue(includeDeleted))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, g)
}

// List handles GET /goals?completed=true|false&deleted=true|false.
func (h *GoalHandler) List(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	completed, err := parseBoolQuery(c, "completed")
	if err != nil {
		response.Error(c, err)
		return
	}
	deletedOnly, err := parseBoolQuery(c, "deleted")
	if err != nil {
		response.Error(c, err)
		return
	}
	includeDeleted, err := parseBoolQuery(c, "include_deleted")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.List(c.Request.Context(), userID, service.ListFilter{
		Completed:      completed,
		Deleted:        boolPtrValue(deletedOnly),
		IncludeDeleted: boolPtrValue(includeDeleted),
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, out)
}

// SoftDelete handles DELETE /goals/:id.
func (h *GoalHandler) SoftDelete(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.svc.SoftDelete(c.Request.Context(), userID, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// Restore handles POST /goals/:id/restore.
func (h *GoalHandler) Restore(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	g, err := h.svc.Restore(c.Request.Context(), userID, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, g)
}

// HardDelete handles DELETE /goals/:id/permanent.
func (h *GoalHandler) HardDelete(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.HardDelete(c.Request.Context(), userID, id); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// Complete handles POST /goals/:id/complete.
func (h *GoalHandler) Complete(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	g, err := h.svc.Complete(c.Request.Context(), userID, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, g)
}

// Uncomplete handles POST /goals/:id/uncomplete.
func (h *GoalHandler) Uncomplete(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	g, err := h.svc.Uncomplete(c.Request.Context(), userID, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, g)
}

// boolPtrValue returns false when ptr is nil, *ptr otherwise.
func boolPtrValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
