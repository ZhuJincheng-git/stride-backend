package model

type User struct {
	BaseEntity
	SoftDeletable

	Username     string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
}

func(User) TableName() string { return "users" }