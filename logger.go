package logar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/Lexographics/logar/options/config"
	"github.com/Lexographics/logar/proxy"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Logger struct {
	db      *gorm.DB
	config  config.Config
	proxies []proxy.Proxy
	actions config.ActionMap
}

type Map map[string]any

func New(opts ...config.ConfigOpt) (*Logger, error) {
	cfg := config.Config{
		AppName:     "logger",
		Database:    "logs.db",
		AutoMigrate: true,
		RequireAuth: false,
		AuthFunc:    nil,
		Models:      config.LogModels{},
		Proxies:     []proxy.Proxy{},
		Actions:     config.ActionMap{},
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
		actions: cfg.Actions,
	}, nil
}

func (l *Logger) Close() error {
	sqlDB, err := l.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (l *Logger) Auth(r *http.Request) bool {
	if l.config.AuthFunc == nil {
		return true
	}

	return l.config.AuthFunc(r)
}

func (l *Logger) GetAllModels() config.LogModels {
	return l.config.Models
}

func (l *Logger) IsAuthCredentialsCorrect(username, password string) bool {
	if l.config.MasterUsername == "" || l.config.MasterPassword == "" {
		return false
	}

	return username == l.config.MasterUsername && password == l.config.MasterPassword
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

func (l *Logger) InvokeAction(path string, args ...any) ([]any, error) {
	action, ok := l.GetActionDetails(path)
	if !ok {
		return nil, fmt.Errorf("action '%s' not found", path)
	}

	actionFunc := reflect.ValueOf(action.Func)
	if actionFunc.Kind() != reflect.Func {
		return nil, fmt.Errorf("path '%s' does not point to a function", path)
	}

	actionType := actionFunc.Type()
	if actionType.NumIn() != len(args) && !(actionType.IsVariadic() && len(args) >= actionType.NumIn()-1) {
		return nil, fmt.Errorf("path '%s' expects %d arguments, got %d", path, actionType.NumIn(), len(args))
	}

	inArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		inArgs[i] = reflect.ValueOf(arg)
	}

	out := actionFunc.Call(inArgs)

	result := make([]any, len(out))
	for i, val := range out {
		result[i] = val.Interface()
	}

	return result, nil
}

func (l *Logger) GetActionArgTypes(path string) ([]reflect.Type, error) {
	action, ok := l.GetActionDetails(path)
	if !ok {
		return nil, fmt.Errorf("action '%s' not found", path)
	}

	actionFunc := reflect.ValueOf(action.Func)
	if actionFunc.Kind() != reflect.Func {
		return nil, fmt.Errorf("path '%s' does not point to a function", path)
	}

	actionType := actionFunc.Type()
	numArgs := actionType.NumIn()
	argTypes := make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		argTypes[i] = actionType.In(i)
	}

	return argTypes, nil
}

func (l *Logger) GetActionsMap() config.ActionMap {
	return l.actions
}

func (l *Logger) GetAllActions() []string {
	actions := []string{}
	for _, action := range l.actions {
		actions = append(actions, action.Path)
	}
	return actions
}

func (l *Logger) GetActionDetails(path string) (config.Action, bool) {
	for _, action := range l.actions {
		if action.Path == path {
			return action, true
		}
	}
	return config.Action{}, false
}
