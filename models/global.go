package models

import (
	"time"

	"sadk.dev/logar/internal/tableprefix"
)

type Global struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Key      string `gorm:"unique"`
	Value    string // can be any JSON-serializable value
	Exported bool   `gorm:"default:false"` // if true, this global will be exported to feature flags
}

func (Global) TableName() string {
	return tableprefix.Get() + "globals"
}
