package logar

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/Lexographics/logar/proxy"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Logger struct {
	db      *gorm.DB
	config  LoggerConfig
	proxies []proxy.Proxy
}

type LoggerConfig struct {
	AppName     string
	Database    string
	AutoMigrate bool
	RequireAuth bool
	AuthFunc    AuthFunc
	Models      LogModels
	Proxies     []proxy.Proxy
}

func New(opts ...ConfigOpt) (*Logger, error) {

	cfg := LoggerConfig{
		AppName:     "logger",
		Database:    "logs.db",
		AutoMigrate: true,
		RequireAuth: false,
		AuthFunc:    nil,
		Models:      LogModels{},
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

func AddProxy(proxy proxy.Proxy) ConfigOpt {
	return func(cfg *LoggerConfig) {
		cfg.Proxies = append(cfg.Proxies, proxy)
	}
}

func Combine(opts ...ConfigOpt) ConfigOpt {
	return func(cfg *LoggerConfig) {
		for _, opt := range opts {
			opt(cfg)
		}
	}
}

func If(condition bool, opts ...ConfigOpt) ConfigOpt {
	if condition {
		return Combine(opts...)
	}
	return func(cfg *LoggerConfig) {}
}

func IfElse(condition bool, ifOpts ConfigOpt, elseOpts ...ConfigOpt) ConfigOpt {
	if condition {
		return ifOpts
	}
	return Combine(elseOpts...)
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
	log := models.Log{
		CreatedAt: now,
		Model:     model,
		Message:   string(msg),
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
