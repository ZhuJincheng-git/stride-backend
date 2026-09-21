package handler

import (
	"github.com/ZhuJincheng-git/stride-backend/internal/middleware"
	"github.com/ZhuJincheng-git/stride-backend/internal/service"
	"github.com/ZhuJincheng-git/stride-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TagHandler exposes both goal-tag and task-tag endpoints.
type TagHandler struct {
	svc *service.TagService
}

func NewTagHandler(s *service.TagService) *TagHandler { return &TagHandler{svc: s} }

type nameRequest struct {
	Name string `json:"name" binding:"required,min=1,max=30"`
}

type tagIDsRequest struct {
	TagIDs []uuid.UUID `json:"tag_ids" binding:"required,min=1"`
}

// --- goal tags ---

func (h *TagHandler) CreateGoalTag(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	var req nameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	t, err := h.svc.CreateGoalTag(c.Request.Context(), userID, req.Name)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, t)
}

func (h *TagHandler) ListGoalTags(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	out, err := h.svc.ListGoalTags(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, out)
}

func (h *TagHandler) DeleteGoalTag(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.DeleteGoalTag(c.Request.Context(), userID, id); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

func (h *TagHandler) AttachGoalTags(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	goalID, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req tagIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	if err := h.svc.AttachGoalTags(c.Request.Context(), userID, goalID, req.TagIDs); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"attached": len(req.TagIDs)})
}

func (h *TagHandler) DetachGoalTags(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	goalID, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req tagIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	if err = h.svc.DetachGoalTags(c.Request.Context(), userID, goalID, req.TagIDs); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"detached": len(req.TagIDs)})
}

func (h *TagHandler) ListGoalTagsForGoal(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	goalID, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.LIstGoalTagsForGoal(c.Request.Context(), userID, goalID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, out)
}

// --- task tags ---

func (h *TagHandler) CreateTaskTag(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	var req nameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	t, err := h.svc.CreateTaskTag(c.Request.Context(), userID, req.Name)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, t)
}

func (h *TagHandler) ListTaskTags(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	out, err := h.svc.ListTaskTags(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, out)
}

func (h *TagHandler) DeleteTaskTag(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.DeleteTaskTag(c.Request.Context(), userID, id); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

func (h *TagHandler) AttachTaskTags(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	taskID, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req tagIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	if err := h.svc.AttachTaskTags(c.Request.Context(), userID, taskID, req.TagIDs); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"attached": len(req.TagIDs)})
}

func (h *TagHandler) DetachTaskTags(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	taskID, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req tagIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	if err = h.svc.DetachTaskTags(c.Request.Context(), userID, taskID, req.TagIDs); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"detached": len(req.TagIDs)})
}

func (h *TagHandler) ListTaskTagsForTask(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	taskID, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.ListTaskTagsForTask(c.Request.Context(), userID, taskID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, out)
}