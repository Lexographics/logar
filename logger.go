package logar

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"sadk.dev/logar/models"
)

type Logger interface {
	Common
	WithContext(ctx context.Context) Logger

	Print(model Model, message any, category string, severity models.Severity) error
	Log(model Model, message any, category string) error
	Info(model Model, message any, category string) error
	Warn(model Model, message any, category string) error
	Error(model Model, message any, category string) error
	Fatal(model Model, message any, category string) error
	Trace(model Model, message any, category string) error

	NewTimer() *Timer
}

type LoggerImpl struct {
	core *AppImpl
	ctx  context.Context
}

func (l *LoggerImpl) GetApp() App {
	return l.core
}

func (l *LoggerImpl) WithContext(ctx context.Context) Logger {
	return &LoggerImpl{core: l.core, ctx: ctx}
}

func (l *LoggerImpl) Print(model Model, message any, category string, severity models.Severity) error {
	var contextualMessage any

	values, ok := l.core.GetContextValues(l.ctx)
	if ok && values != nil && len(values) > 0 {
		originalMsgMap, isMap := message.(Map)
		if isMap {
			for k, v := range values {
				originalMsgMap[k] = v
			}
			contextualMessage = originalMsgMap
		} else {
			newMsgMap := make(Map)
			for k, v := range values {
				newMsgMap[k] = v
			}
			newMsgMap["message"] = message
			contextualMessage = newMsgMap
		}
	} else {
		contextualMessage = message
	}

	var msg string

	switch m := contextualMessage.(type) {
	case string:
		msg = m
	default:
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		encoder.Encode(m)
		msg = buf.String()
	}

	now := time.Now()
	logEntry := models.Log{
		CreatedAt: now,
		Model:     models.Model(model),
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

func (l *LoggerImpl) Log(model Model, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Log)
}

func (l *LoggerImpl) Info(model Model, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Info)
}

func (l *LoggerImpl) Warn(model Model, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Warning)
}

func (l *LoggerImpl) Error(model Model, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Error)
}

func (l *LoggerImpl) Fatal(model Model, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Fatal)
}

func (l *LoggerImpl) Trace(model Model, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Trace)
}

func (l *LoggerImpl) NewTimer() *Timer {
	return l.core.NewTimer()
}
