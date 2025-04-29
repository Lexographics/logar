package models

import "time"

type Session struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	ExpiresAt time.Time `gorm:"not null"`

	Username string `gorm:"not null"`
	Token    string `gorm:"not null;unique"`
}
