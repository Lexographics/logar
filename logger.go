package logar

import (
	"encoding/json"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
)

type Logger interface {
	Common

	Print(model string, message any, category string, severity models.Severity) error
	Log(model string, message any, category string) error
	Info(model string, message any, category string) error
	Warn(model string, message any, category string) error
	Error(model string, message any, category string) error
	Fatal(model string, message any, category string) error
	Trace(model string, message any, category string) error

	NewTimer() *Timer
}

type LoggerImpl struct {
	core *AppImpl
}

func (l *LoggerImpl) GetApp() App {
	return l.core
}

func (l *LoggerImpl) Print(model string, message any, category string, severity models.Severity) error {
	var msg string

	switch m := message.(type) {
	case string:
		msg = m
	default:
		msgBytes, err := json.Marshal(message)
		if err != nil {
			return err
		}
		msg = string(msgBytes)
	}

	now := time.Now()
	logEntry := models.Log{
		CreatedAt: now,
		Model:     model,
		Message:   msg,
		Category:  category,
		Severity:  severity,
	}
	err := l.core.db.Create(&logEntry).Error
	if err != nil {
		return err
	}

	for _, p := range l.core.proxies {
		p.TrySend(logEntry, msg)
	}

	return nil
}

func (l *LoggerImpl) Log(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Log)
}

func (l *LoggerImpl) Info(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Info)
}

func (l *LoggerImpl) Warn(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Warning)
}

func (l *LoggerImpl) Error(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Error)
}

func (l *LoggerImpl) Fatal(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Fatal)
}

func (l *LoggerImpl) Trace(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Trace)
}

func (l *LoggerImpl) NewTimer() *Timer {
	return l.core.NewTimer()
}
