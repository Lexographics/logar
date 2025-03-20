package logar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Logger struct {
	db     *gorm.DB
	config LoggerConfig
}

type LoggerConfig struct {
	AppName           string
	Database          string
	AutoMigrate       bool
	RequireAuth       bool
	AuthFunc          AuthFunc
	Models            LogModels
	PrintedSeverities []Severity
}

func New(opts ...ConfigOpt) (*Logger, error) {

	cfg := LoggerConfig{
		AppName:           "logger",
		Database:          "logs.db",
		AutoMigrate:       true,
		RequireAuth:       false,
		AuthFunc:          nil,
		Models:            LogModels{},
		PrintedSeverities: []Severity{Severity_Log, Severity_Info, Severity_Warning, Severity_Error, Severity_Fatal, Severity_Trace},
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
		err = db.AutoMigrate(&Log{})
		if err != nil {
			return nil, err
		}
	}

	return &Logger{
		db:     db,
		config: cfg,
	}, nil
}

type LogModel struct {
	DisplayName string
	ModelId     string
}
type LogModels []LogModel

type AuthFunc func(r *http.Request) bool
type ConfigOpt func(*LoggerConfig)

func WithAppName(appName string) ConfigOpt {
	return func(cfg *LoggerConfig) {
		cfg.AppName = appName
	}
}

func WithDatabase(database string) ConfigOpt {
	return func(cfg *LoggerConfig) {
		cfg.Database = database
	}
}

func EnableAutoMigrate(cfg *LoggerConfig) {
	cfg.AutoMigrate = true
}

func DisableAutoMigrate(cfg *LoggerConfig) {
	cfg.AutoMigrate = false
}

func WithAuth(authFunc AuthFunc) ConfigOpt {
	return func(cfg *LoggerConfig) {
		cfg.RequireAuth = true
		cfg.AuthFunc = authFunc
	}
}

func AddModel(displayName, modelId string) ConfigOpt {
	return func(cfg *LoggerConfig) {
		cfg.Models = append(cfg.Models, LogModel{
			DisplayName: displayName,
			ModelId:     modelId,
		})
	}
}

func WithPrintedSeverities(severities ...Severity) ConfigOpt {
	return func(cfg *LoggerConfig) {
		cfg.PrintedSeverities = severities
	}
}

func (l *Logger) Close() error {
	sqlDB, err := l.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (l *Logger) Print(model string, message any, category string, severity Severity) error {
	msg, error := json.MarshalIndent(message, "", "  ")
	if error != nil {
		return error
	}

	now := time.Now()
	if slices.Contains(l.config.PrintedSeverities, severity) {
		fmt.Println("["+strings.ToUpper(severity.String())+"]", now.Format(time.DateTime), string(msg))
	}

	return l.db.Create(&Log{
		CreatedAt: now,
		Model:     model,
		Message:   string(msg),
		Category:  category,
		Severity:  severity,
	}).Error
}

func (l *Logger) Log(model string, message any, category string) error {
	return l.Print(model, message, category, Severity_Log)
}

func (l *Logger) Info(model string, message any, category string) error {
	return l.Print(model, message, category, Severity_Info)
}

func (l *Logger) Warn(model string, message any, category string) error {
	return l.Print(model, message, category, Severity_Warning)
}

func (l *Logger) Error(model string, message any, category string) error {
	return l.Print(model, message, category, Severity_Error)
}

func (l *Logger) Fatal(model string, message any, category string) error {
	return l.Print(model, message, category, Severity_Fatal)
}

func (l *Logger) Trace(model string, message any, category string) error {
	return l.Print(model, message, category, Severity_Trace)
}
