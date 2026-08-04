package model

import (
	"time"

	"github.com/google/uuid"
)

// Task is a concrete unit of work that, optionally, advances a Goal.
type Task struct {
	BaseEntity
	SoftDeletable
	Completable

	UserID            uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	Title             string     `gorm:"type:varchar(100);not null" json:"title"`
	Description       string     `gorm:"type:text" json:"description,omitempty"`
	ExpectedStartTime *time.Time `json:"expected_start_time,omitempty"`
	ExpectedEndTime   *time.Time `json:"expected_end_time,omitempty"`
	GoalID            *uuid.UUID `gorm:"type:char(36)" json:"goal_id,omitempty"`
	ParentTaskID      *uuid.UUID `gorm:"type:char(36)" json:"parent_task_id,omitempty"`

	Tags        []TaskTag        `gorm:"many2many:task_tag_rel;joinForeignKey:TaskID;joinReferences:TaskTagID" json:"tags,omitempty"`
	Metrics     []TaskMetric     `gorm:"many2many:task_metric_rel;joinForeignKey:TaskID;joinReferences:TaskMetricID" join:"metrics,omitempty"`
	WorkRecords []TaskWorkRecord `gorm:"foreignKey:WorkTaskID"`
}

func (Task) TableName() string { return "tasks" }
