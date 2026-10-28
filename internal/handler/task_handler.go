package handler

import (
	"time"

	"github.com/ZhuJincheng-git/stride-backend/internal/middleware"
	"github.com/ZhuJincheng-git/stride-backend/internal/service"
	"github.com/ZhuJincheng-git/stride-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskService(s *service.TaskService) *TaskHandler { return &TaskHandler{svc: s} }

type taskUpsertRequest struct {
	Title             *string    `json:"title" binding:"omitempty,min=1,max=100"`
	Description       *string    `json:"description" binding:"omitempty,max=10000"`
	ExpectedStartTime *time.Time `json:"expected_start_time"`
	ExpectedEndTime   *time.Time `json:"expected_end_time"`
	GoalID            *uuid.UUID `json:"goal_id"`
	ParentTaskID      *uuid.UUID `json:"parent_task_id"`
	ClearGoal         bool       `json:"clear_goal"`
	ClearParent       bool       `json:"clear_parent"`
}

// Create handles POST /tasks.
func (h *TaskHandler) Create(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	var req taskUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	t, err := h.svc.Create(c.Request.Context(), userID, service.TaskInput{
		Title:             req.Title,
		Description:       req.Description,
		ExpectedStartTime: req.ExpectedStartTime,
		ExpectedEndTime:   req.ExpectedEndTime,
		GoalID:            req.GoalID,
		ParentTaskID:      req.ParentTaskID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, t)
}

// Update handles PATCH /tasks/:id.
func (h *TaskHandler) Update(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req taskUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bindingError(err))
		return
	}
	t, err := h.svc.Update(c.Request.Context(), userID, id, service.TaskInput{
		Title:             req.Title,
		Description:       req.Description,
		ExpectedStartTime: req.ExpectedStartTime,
		ExpectedEndTime:   req.ExpectedEndTime,
		GoalID:            req.GoalID,
		ParentTaskID:      req.ParentTaskID,
		ClearGoal:         req.ClearGoal,
		ClearParent:       req.ClearParent,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, t)
}

// Get handles GET /tasks/:id.
func (h *TaskHandler) Get(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	includeDeleted, err := parseBoolQuery(c, "includeDeleted")
	if err != nil {
		response.Error(c, err)
		return
	}
	t, err := h.svc.Get(c.Request.Context(), userID, id, boolPtrValue(includeDeleted))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, t)
}

// List handles GET /tasks?goal_id=...&completed=true|false&deleted=...
func (h *TaskHandler) List(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	goalID, err := parseUUIDQuery(c, "goal_id")
	if err != nil {
		response.Error(c, err)
		return
	}
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
	noGoal, err := parseBoolQuery(c, "no_goal")
	if err != nil {
		response.Error(c, err)
		return
	}

	f := service.TaskListFilter{
		Completed:      completed,
		Deleted:        boolPtrValue(deletedOnly),
		IncludeDeleted: boolPtrValue(includeDeleted),
		NoGoal:         boolPtrValue(noGoal),
	}
	f.GoalID = goalID
	out, err := h.svc.List(c.Request.Context(), userID, f)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, out)
}

// SoftDelete handles DELETE /tasks/:id
func (h *TaskHandler) SoftDelete(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.SoftDelete(c.Request.Context(), userID, id); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// Restore handles POST /tasks/:id/restore.
func (h *TaskHandler) Restore(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	t, err := h.svc.Restore(c.Request.Context(), userID, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, t)
}

// HardDelete handles DELETE /tasks/:id/permanent.
func (h *TaskHandler) HardDelete(c *gin.Context) {
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

// Complete handles POST /tasks/:id/complete.
func (h *TaskHandler) Complete(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	t, err := h.svc.Complete(c.Request.Context(), userID, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, t)
}

// Uncomplete handles POST /tasks/:id/uncomplete.
func (h *TaskHandler) Uncomplete(c *gin.Context) {
	userID, _ := middleware.CurrentUserID(c)
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	t, err := h.svc.Uncomplete(c.Request.Context(), userID, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, t)
}
