package model

import (
	"time"

	"github.com/google/uuid"
)

// GoalMetric is a user-defined numeric metric for goals.
type GoalMetric struct {
	BaseEntity
	SoftDeletable

	UserID uuid.UUID `gorm:"type:char(30);not null" json:"user_id"`
	Name   string    `gorm:"type:varchar(30);not null" json:"name"`
}

func (GoalMetric) TableName() string { return "goal_metrics" }

// GoalMetricRel records a metric value for a goal at a point in time.
type GoalMetricRel struct {
	GoalID             uuid.UUID `gorm:"type:char(36);primaryKey" json:"goal_id"`
	GoalMetricID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"goal_metric_id"`
	EvaluatedAt        time.Time `gorm:"primaryKey;autoCreateTime" json:"evaluate_at"`
	ProspectiveValue   float64   `gorm:"not null" json:"prospective_value"`
	RetrospectiveValue *float64  `json:"retrospective_value,omitempty"`
}

func (GoalMetricRel) TableName() string { return "goal_metric_rel" }

type TaskMetric struct {
	BaseEntity
	SoftDeletable

	UserID uuid.UUID `gorm:"type:char(30);not null" json:"user_id"`
	Name   string    `gorm:"type:varchar(30);not null" json:"name"`
}

func (TaskMetric) TableName() string { return "goal_metrics" }

type TaskMetricRel struct {
	TaskID             uuid.UUID `gorm:"type:char(36);primaryKey" json:"task_id"`
	TaskMetricID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"task_metric_id"`
	EvaluatedAt        time.Time `gorm:"primaryKey;autoCreateTime" json:"evaluate_at"`
	ProspectiveValue   float64   `gorm:"not null" json:"prospective_value"`
	RetrospectiveValue *float64  `json:"retrospective_value,omitempty"`
}

func (TaskMetricRel) TableName() string { return "task_metric_rel" }
