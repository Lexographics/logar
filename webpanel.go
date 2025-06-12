package logar

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"sadk.dev/logar/models"
)

type WebPanel interface {
	Common

	LoginUser(username, password string) (models.User, error)
	CreateUser(username, displayName, password string, isAdmin bool) (models.User, error)
	GetUser(id uint) (models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user models.User) error
	CreateSession(user models.User, device string) (string, error)
	DeleteSession(token string) error
	GetSession(token string) (*models.Session, error)
	GetActiveSessions(userID uint) ([]models.Session, error)
	GetDefaultLanguage() Language
	Auth(r *http.Request) bool
}

type WebPanelImpl struct {
	core *AppImpl
}

func (w *WebPanelImpl) GetApp() App {
	return w.core
}

func (w *WebPanelImpl) LoginUser(username, password string) (models.User, error) {
	if w.core.config.AdminUsername == username && w.core.config.AdminPassword == password {
		return models.User{
			Username:    username,
			DisplayName: "Admin",
			IsAdmin:     true,
		}, nil
	}

	var user models.User
	err := w.core.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return models.User{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return models.User{}, fmt.Errorf("invalid password")
	}

	err = w.core.db.Table("users").Where("id = ?", user.ID).Update("last_activity", time.Now()).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (w *WebPanelImpl) CreateUser(username, displayName, password string, isAdmin bool) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Username:    username,
		DisplayName: displayName,
		Password:    string(hashedPassword),
		IsAdmin:     isAdmin,
	}

	err = w.core.db.Create(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (w *WebPanelImpl) GetUser(id uint) (models.User, error) {
	if id == 0 {
		return models.User{
			ID:          0,
			Username:    w.core.config.AdminUsername,
			DisplayName: "Admin",
			IsAdmin:     true,
		}, nil
	}

	var user models.User
	err := w.core.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (w *WebPanelImpl) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := w.core.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	users = append(users, models.User{
		ID:          0,
		Username:    w.core.config.AdminUsername,
		DisplayName: "System",
		IsAdmin:     true,
	})
	return users, nil
}

func (w *WebPanelImpl) UpdateUser(user models.User) error {
	return w.core.db.Save(&user).Error
}

func (w *WebPanelImpl) CreateSession(user models.User, device string) (string, error) {
	session := models.Session{
		UserID:       user.ID,
		ExpiresAt:    time.Now().Add(w.core.config.WebPanelConfig.SessionDuration),
		Token:        uuid.New().String(),
		Device:       device,
		LastActivity: time.Now(),
	}

	err := w.core.db.Create(&session).Error
	if err != nil {
		return "", err
	}
	return session.Token, nil
}

func (w *WebPanelImpl) DeleteSession(token string) error {
	return w.core.db.Where("token = ?", token).Delete(&models.Session{}).Error
}

func (w *WebPanelImpl) GetSession(token string) (*models.Session, error) {
	var session models.Session
	err := w.core.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		w.DeleteSession(token)
		return nil, fmt.Errorf("session expired")
	}

	now := time.Now()
	session.LastActivity = now
	session.ExpiresAt = now.Add(w.core.config.WebPanelConfig.SessionDuration)
	err = w.core.db.Save(&session).Error
	if err != nil {
		return nil, err
	}
	err = w.core.db.Model(&models.User{}).Where("id = ?", session.UserID).Update("last_activity", now).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (w *WebPanelImpl) GetActiveSessions(userID uint) ([]models.Session, error) {
	var sessions []models.Session
	err := w.core.db.Where("user_id = ? and expires_at > ?", userID, time.Now()).Order("last_activity DESC").Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (w *WebPanelImpl) GetDefaultLanguage() Language {
	return w.core.config.DefaultLanguage
}

func (w *WebPanelImpl) Auth(r *http.Request) bool {
	if w.core.config.AuthFunc == nil {
		return true
	}

	return w.core.config.AuthFunc(r)
}
