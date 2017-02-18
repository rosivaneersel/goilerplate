package users

import (
	"time"
)

type Profile struct {
	GID uint          `gorm:"column:id; primary_key" bson:"-"`
	UserID   uint	`bson:"-"`
	FirstName string
	LastName string
	AvatarURL string
	CreatedAt 	time.Time
	UpdatedAt	time.Time
	DeletedAt	*time.Time `sql:”index”`
}
