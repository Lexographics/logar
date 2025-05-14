package logar

import (
	"context"
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
	GetLogger() Logger
	GetActionManager() ActionManager
	GetWebPanel() WebPanel
	GetAnalytics() Analytics
	// GetFeatureFlags() FeatureFlags

	Close() error
	GetAllModels() LogModels
	SetTypeKind(type_ reflect.Type, kind TypeKind)
	SetTypeKindString(type_ string, kind TypeKind)
	GetTypeKind(type_ reflect.Type) (TypeKind, bool)
	GetTypeKindString(type_ string) (TypeKind, bool)

	PrepareContext(parent context.Context, values Map) context.Context
	GetContextValues(ctx context.Context) (Map, bool)
	GetFromContext(ctx context.Context, key string) (any, bool)
	AddContextValue(ctx context.Context, key string, value any) App
}

type AppImpl struct {
	logger        Logger
	actionManager ActionManager
	webPanel      WebPanel
	analytics     Analytics
	// featureFlags  FeatureFlags

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
		Logger: logger.Discard,
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
		// &FeatureFlag{},
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

	logger.logger = &LoggerImpl{core: logger}
	logger.actionManager = &ActionManagerImpl{core: logger}
	logger.webPanel = &WebPanelImpl{core: logger}
	logger.analytics = &AnalyticsImpl{core: logger}
	// logger.featureFlags = &FeatureFlagsImpl{core: logger}

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

func (l *AppImpl) GetLogger() Logger {
	return l.logger
}

func (l *AppImpl) GetActionManager() ActionManager {
	return l.actionManager
}

func (l *AppImpl) GetWebPanel() WebPanel {
	return l.webPanel
}

func (l *AppImpl) GetAnalytics() Analytics {
	return l.analytics
}

// func (l *AppImpl) GetFeatureFlags() FeatureFlags {
// 	return l.featureFlags
// }

func (l *AppImpl) PrepareContext(parent context.Context, values Map) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	value := Map{}
	for k, v := range values {
		value[k] = v
	}
	return context.WithValue(parent, logarContextKey, &value)
}

func (l *AppImpl) AddContextValue(ctx context.Context, key string, value any) App {
	if ctx == nil {
		return l
	}
	values, ok := l.GetContextValues(ctx)
	if !ok {
		return l
	}
	values[key] = value
	return l
}

func (l *AppImpl) GetContextValues(ctx context.Context) (Map, bool) {
	if ctx == nil {
		return nil, false
	}
	value, ok := ctx.Value(logarContextKey).(*Map)
	if !ok {
		return nil, false
	}
	return *value, true
}

func (l *AppImpl) GetFromContext(ctx context.Context, key string) (any, bool) {
	if ctx == nil {
		return nil, false
	}
	values, ok := l.GetContextValues(ctx)
	if !ok {
		return nil, false
	}
	value, ok := values[key]
	return value, ok
}
