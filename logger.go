package logar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/Lexographics/logar/proxy"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TypeKind string

const (
	TypeKind_Text     TypeKind = "text"
	TypeKind_Int      TypeKind = "int"
	TypeKind_Float    TypeKind = "float"
	TypeKind_Bool     TypeKind = "bool"
	TypeKind_Time     TypeKind = "time"
	TypeKind_Duration TypeKind = "duration"
)

type Logger struct {
	db        *gorm.DB
	config    Config
	proxies   []proxy.Proxy
	actions   Actions
	typeKinds map[string]TypeKind
}

type Map map[string]any

func New(opts ...ConfigOpt) (*Logger, error) {
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
	)
	if err != nil {
		return nil, err
	}

	// Delete expired sessions
	err = db.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
	if err != nil {
		return nil, err
	}

	logger := &Logger{
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

func (l *Logger) GetAllModels() LogModels {
	return l.config.Models
}

func (l *Logger) LoginUser(username, password string) (models.User, error) {
	if l.config.MasterUsername == username && l.config.MasterPassword == password {
		return models.User{
			Username:    username,
			DisplayName: "Master",
			IsAdmin:     true,
		}, nil
	}

	var user models.User
	err := l.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return models.User{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return models.User{}, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (l *Logger) CreateUser(username, displayName, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username:    username,
		DisplayName: displayName,
		Password:    string(hashedPassword),
	}

	return l.db.Create(&user).Error
}
func (l *Logger) CreateSession(user models.User, device string) (string, error) {
	session := models.Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
		Token:     uuid.New().String(),
		Device:    device,
	}

	err := l.db.Create(&session).Error
	if err != nil {
		return "", err
	}
	return session.Token, nil
}

func (l *Logger) DeleteSession(token string) error {
	return l.db.Where("token = ?", token).Delete(&models.Session{}).Error
}

func (l *Logger) GetSession(token string) (*models.Session, error) {
	var session models.Session
	err := l.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		l.DeleteSession(token)
		return nil, fmt.Errorf("session expired")
	}

	session.LastActivity = time.Now()
	session.ExpiresAt = time.Now().Add(time.Hour * 24)
	err = l.db.Save(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (l *Logger) GetActiveSessions(userID uint) ([]models.Session, error) {
	var sessions []models.Session
	err := l.db.Where("user_id = ? and expires_at > ?", userID, time.Now()).Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (l *Logger) GetTypeKind(type_ reflect.Type) (TypeKind, bool) {
	kind, ok := l.typeKinds[type_.String()]
	return kind, ok
}

func (l *Logger) SetTypeKind(type_ reflect.Type, kind TypeKind) {
	l.typeKinds[type_.String()] = kind
}

func (l *Logger) GetTypeKindString(type_ string) (TypeKind, bool) {
	kind, ok := l.typeKinds[type_]
	return kind, ok
}

func (l *Logger) SetTypeKindString(type_ string, kind TypeKind) {
	l.typeKinds[type_] = kind
}

func (l *Logger) GetDefaultLanguage() Language {
	return l.config.DefaultLanguage
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

func (l *Logger) GetActionsMap() Actions {
	return l.actions
}

func (l *Logger) GetAllActions() []string {
	actions := []string{}
	for _, action := range l.actions {
		actions = append(actions, action.Path)
	}
	return actions
}

func (l *Logger) GetActionDetails(path string) (Action, bool) {
	for _, action := range l.actions {
		if action.Path == path {
			return action, true
		}
	}
	return Action{}, false
}

func (l *Logger) AddAction(action Action) {
	for i, existingAction := range l.actions {
		if existingAction.Path == action.Path {
			l.actions[i] = action
			return
		}
	}
	l.actions = append(l.actions, action)
}

func (l *Logger) RemoveAction(path string) {
	for i, action := range l.actions {
		if action.Path == path {
			l.actions = append(l.actions[:i], l.actions[i+1:]...)
			return
		}
	}
}
