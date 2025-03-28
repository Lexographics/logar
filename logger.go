package logar

import (
	"encoding/json"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/Lexographics/logar/internal/options/config"
	"github.com/Lexographics/logar/proxy"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Logger struct {
	db      *gorm.DB
	config  config.Config
	proxies []proxy.Proxy
}

func New(opts ...config.ConfigOpt) (*Logger, error) {

	cfg := config.Config{
		AppName:     "logger",
		Database:    "logs.db",
		AutoMigrate: true,
		RequireAuth: false,
		AuthFunc:    nil,
		Models:      config.LogModels{},
		Proxies:     []proxy.Proxy{},
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	db, err := gorm.Open(sqlite.Open("file:"+cfg.Database+"?cache=shared&mode=rwc&_journal_mode=WAL"), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		return nil, err
	}

	if cfg.AutoMigrate {
		err = db.AutoMigrate(&models.Log{})
		if err != nil {
			return nil, err
		}
	}

	return &Logger{
		db:      db,
		config:  cfg,
		proxies: cfg.Proxies,
	}, nil
}

func (l *Logger) Close() error {
	sqlDB, err := l.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (l *Logger) Print(model string, message any, category string, severity models.Severity) error {
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

func (l *Logger) Log(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Log)
}

func (l *Logger) Info(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Info)
}

func (l *Logger) Warn(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Warning)
}

func (l *Logger) Error(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Error)
}

func (l *Logger) Fatal(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Fatal)
}

func (l *Logger) Trace(model string, message any, category string) error {
	return l.Print(model, message, category, models.Severity_Trace)
}
