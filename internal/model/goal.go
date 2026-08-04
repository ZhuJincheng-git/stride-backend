package model

import (
	"time"

	"github.com/google/uuid"
)

// Goal is a long-running objective the user is pursuing.
type Goal struct {
	BaseEntity
	SoftDeletable
	Completable

	UserID            uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`
	Title             string     `gorm:"type:varchar(100);not null" json:"title"`
	Description       string     `gorm:"type:text" json:"description,omitempty"`
	ExpectedStartTime *time.Time `json:"expected_start_time,omitempty"`
	ExpectedEndTime   *time.Time `json:"expected_end_time,omitempty"`
	ParentGoalID      *uuid.UUID `gorm:"type:char(36)" json:"parent_goal_id,omitempty"`

	Tags    []GoalTag    `gorm:"many2many:goal_tag_rel;joinForeignKey:GoalID;joinReferences:GoalTagID" json:"tags:omitempty"`
	Metrics []GoalMetric `gorm:"many2many:goal_metric_rel;joinForeignKey:GoalID;joinReferences:GoalMetricID" json:"metrics,omitempty"`
}

func (Goal) TableName() string { return "goals" }
