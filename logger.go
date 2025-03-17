package logar

import (
	"encoding/json"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Logger struct {
	db     *gorm.DB
	config LoggerConfig
}

type LoggerConfig struct {
	AppName     string
	Database    string
	AutoMigrate bool
	RequireAuth bool
	AuthFunc    AuthFunc
	Models      LogModels
}

func New(opts ...ConfigOpt) (*Logger, error) {
	cfg := LoggerConfig{
		AppName:     "logger",
		Database:    "logs.db",
		AutoMigrate: true,
		RequireAuth: false,
		AuthFunc:    nil,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	db, err := gorm.Open(sqlite.Open("file:" + cfg.Database + "?cache=shared&mode=rwc&_journal_mode=WAL"))
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

	return l.db.Create(&Log{
		Model:    model,
		Message:  string(msg),
		Category: category,
		Severity: severity,
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
