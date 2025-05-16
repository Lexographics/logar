package models

import "github.com/Lexographics/logar/internal/tableprefix"

type FeatureFlag struct {
	ID uint `json:"id" gorm:"primary_key"`

	Enabled   bool   `json:"enabled"`
	Name      string `json:"name" gorm:"unique"`
	Condition string `json:"condition"` // expr-lang expression
}

func (FeatureFlag) TableName() string {
	return tableprefix.Get() + "feature_flags"
}
