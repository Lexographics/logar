package logar

import (
	"encoding/json"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
)

type Logger interface {
	Print(model string, message any, category string, severity models.Severity) error
	Log(model string, message any, category string) error
	Info(model string, message any, category string) error
	Warn(model string, message any, category string) error
	Error(model string, message any, category string) error
	Fatal(model string, message any, category string) error
	Trace(model string, message any, category string) error

	NewTimer() *Timer
}

func (l *AppImpl) Print(model string, message any, category string, severity models.Severity) error {
	var msg string

	switch message := message.(type) {
	case string:
		msg = message
	default:
		msgBytes, err := json.Marshal(message)
		if err != nil {
			return err
		}
		msg = string(msgBytes)
	}

	now := time.Now()
	log := models.Log{
		CreatedAt: now,
		Model:     model,
		Message:   msg,
		Category:  category,
		Severity:  severity,
	}
	err := l.db.Create(&log).Error
	if err != nil {
		return err
	}

	for _, proxy := range l.proxies {
		proxy.TrySend(log, msg)
	}

	return nil
}

func (l *AppImpl) Log(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Log)
}

func (l *AppImpl) Info(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Info)
}

func (l *AppImpl) Warn(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Warning)
}

func (l *AppImpl) Error(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Error)
}

func (l *AppImpl) Fatal(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Fatal)
}

func (l *AppImpl) Trace(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Trace)
}
