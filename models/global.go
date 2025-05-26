package models

import (
	"time"

	"github.com/Lexographics/logar/internal/tableprefix"
)

type Global struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Key   string `gorm:"unique"`
	Value string // can be any JSON-serializable value
}

func (Global) TableName() string {
	return tableprefix.Get() + "globals"
}
