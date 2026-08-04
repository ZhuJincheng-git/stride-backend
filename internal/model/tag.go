package model

import "github.com/google/uuid"

// GoalTag is a user-defined label that can be attached to many goals.
type GoalTag struct {
	BaseEntity
	SoftDeletable

	UserID uuid.UUID `gorm:"type:char(36);not null;uniqueIndex:idx_goal_tag_user_name" json:"user_id"`
	Name   string    `gorm:"type:varchar(30);not null;uniqueIndex:idx_goal_tag_user_name" json:"name"`
}

func (GoalTag) TableName() string { return "goal_tags" }

// GoalTagRel is the join table between Goal and GoalTag.
type GoalTagRel struct {
	GoalID    uuid.UUID `gorm:"type:char(36);primaryKey"`
	GoalTagID uuid.UUID `gorm:"type:char(36);primaryKey"`
}

func (GoalTagRel) TableName() string { return "goal_tag_rel" }
