package models

import "time"

type User struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`

	Username    string `json:"username" gorm:"not null;unique"`
	DisplayName string `json:"display_name"`
	Password    string `json:"-" gorm:"not null"`
	IsAdmin     bool   `json:"is_admin" gorm:"not null;default:false"`
}
