package models

import "time"

type Session struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	ExpiresAt    time.Time `gorm:"not null"`
	LastActivity time.Time `gorm:"not null"`
	Device       string

	UserID uint   `gorm:"not null"`
	Token  string `gorm:"not null;unique"`
}
