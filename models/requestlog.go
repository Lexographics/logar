package models

import (
	"time"

	"sadk.dev/logar/internal/tableprefix"
)

// RequestLog stores information about a single request.
// It's also a database model for database persistence.
type RequestLog struct {
	ID uint `gorm:"primarykey"`

	Timestamp  time.Time `gorm:"index"`
	VisitorID  string    `gorm:"index"`
	Instance   string
	Path       string
	Latency    time.Duration // Stored as int64 (nanoseconds)
	StatusCode int
	UserAgent  string
	OS         string
	Browser    string
	Referer    string
	BytesSent  int64
	BytesRecv  int64
}

func (RequestLog) TableName() string {
	return tableprefix.Get() + "request_logs"
}
