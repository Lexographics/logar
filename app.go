package logar

import (
	"reflect"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/Lexographics/logar/proxy"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Common interface {
	GetApp() App
}

// App is the main struct that contains library data for things like logging, actions, etc.
type App interface {
	Common
	Logger
	ActionManager
	WebPanel
	Analytics

	Close() error
	GetAllModels() LogModels
	SetTypeKind(type_ reflect.Type, kind TypeKind)
	SetTypeKindString(type_ string, kind TypeKind)
	GetTypeKind(type_ reflect.Type) (TypeKind, bool)
	GetTypeKindString(type_ string) (TypeKind, bool)
}

type AppImpl struct {
	db        *gorm.DB
	config    Config
	proxies   []proxy.Proxy
	actions   Actions
	typeKinds map[string]TypeKind
}

func New(opts ...ConfigOpt) (App, error) {
	cfg := Config{
		AppName:         "logger",
		Database:        "logs.db",
		RequireAuth:     false,
		AuthFunc:        nil,
		Models:          LogModels{},
		Proxies:         []proxy.Proxy{},
		Actions:         Actions{},
		DefaultLanguage: English,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	db, err := gorm.Open(sqlite.Open("file:"+cfg.Database+"?cache=shared&mode=rwc&_journal_mode=WAL"), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		return nil, err
	}

	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqldb.SetMaxOpenConns(1)

	err = db.AutoMigrate(
		&models.Log{},
		&models.Session{},
		&models.User{},
		&RequestLog{},
	)
	if err != nil {
		return nil, err
	}

	// Delete expired sessions
	err = db.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
	if err != nil {
		return nil, err
	}

	logger := &AppImpl{
		db:        db,
		config:    cfg,
		proxies:   cfg.Proxies,
		actions:   cfg.Actions,
		typeKinds: map[string]TypeKind{},
	}

	// Default type kinds
	logger.SetTypeKind(reflect.TypeOf(string("")), TypeKind_Text)
	logger.SetTypeKind(reflect.TypeOf(int(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(int8(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(int16(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(int32(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(int64(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(uint(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(uint8(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(uint16(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(uint32(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(uint64(0)), TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(float32(0)), TypeKind_Float)
	logger.SetTypeKind(reflect.TypeOf(float64(0)), TypeKind_Float)
	logger.SetTypeKind(reflect.TypeOf(bool(false)), TypeKind_Bool)
	logger.SetTypeKind(reflect.TypeOf(time.Time{}), TypeKind_Time)
	logger.SetTypeKind(reflect.TypeOf(time.Duration(0)), TypeKind_Duration)
	return logger, nil
}

func (l *AppImpl) Close() error {
	sqlDB, err := l.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (l *AppImpl) GetAllModels() LogModels {
	return l.config.Models
}

func (l *AppImpl) GetTypeKind(type_ reflect.Type) (TypeKind, bool) {
	kind, ok := l.typeKinds[type_.String()]
	return kind, ok
}

func (l *AppImpl) SetTypeKind(type_ reflect.Type, kind TypeKind) {
	l.typeKinds[type_.String()] = kind
}

func (l *AppImpl) GetTypeKindString(type_ string) (TypeKind, bool) {
	kind, ok := l.typeKinds[type_]
	return kind, ok
}

func (l *AppImpl) SetTypeKindString(type_ string, kind TypeKind) {
	l.typeKinds[type_] = kind
}

func (l *AppImpl) GetApp() App {
	return l
}
