package repository

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }

// onConflictDoNothing builds the clause used by bulk-inserts that should be
// idempotent: "INSERT ... ON CONFLICT(<cols>) DO NOTHING".
func onConflictDoNothing(cols ...string) []clause.Expression {
	c := clause.OnConflict{
		Columns: make([]clause.Column, 0, len(cols)),
		DoNothing: true,
	}
	for _, col := range cols {
		c.Columns = append(c.Columns, clause.Column{Name: col})
	}
	return []clause.Expression{c}
}