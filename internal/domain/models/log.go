package models

import (
	"time"
)

type Log struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Model     string
	Message   string
	Category  string
	Severity  Severity
}

type Severity int

const (
	Severity_None Severity = iota
	Severity_Log
	Severity_Info
	Severity_Warning
	Severity_Error
	Severity_Fatal
	Severity_Trace
	Severity_Max
)

var severityStrings = [...]string{"None", "Log", "Info", "Warn", "Error", "Fatal", "Trace"}

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
