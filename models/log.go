package models

import (
	"time"

	"github.com/Lexographics/logar/internal/tableprefix"
)

type Log struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Model     Model
	Message   string
	Category  string
	Severity  Severity
}

func (Log) TableName() string {
	return tableprefix.Get() + "logs"
}

func (l Log) FieldNames() []string {
	return []string{
		"id",
		"created_at",
		"message",
		"category",
		"severity",
	}
}

type Severity int

const (
	Severity_None Severity = iota
	Severity_Trace
	Severity_Log
	Severity_Info
	Severity_Warning
	Severity_Error
	Severity_Fatal
	Severity_Max
)

var severityStrings = [...]string{"None", "Trace", "Log", "Info", "Warn", "Error", "Fatal"}

func (s Severity) String() string {
	s.Clamp()
	return severityStrings[s]
}

func (s *Severity) Clamp() {
	if *s >= Severity_Max {
		*s = Severity_Max - 1
	}
	if *s < 0 {
		*s = 0
	}
}
