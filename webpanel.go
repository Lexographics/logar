package logar

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type WebPanel interface {
	Common
	GetWebPanel() WebPanel

	LoginUser(username, password string) (models.User, error)
	CreateUser(username, displayName, password string, isAdmin bool) (models.User, error)
	GetUser(id uint) (models.User, error)
	GetDefaultLanguage() Language
	Auth(r *http.Request) bool
}

func (l *AppImpl) GetWebPanel() WebPanel {
	return l
}

func (l *AppImpl) LoginUser(username, password string) (models.User, error) {
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

	err = l.db.Table("users").Where("id = ?", user.ID).Update("last_activity", time.Now()).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (l *AppImpl) CreateUser(username, displayName, password string, isAdmin bool) (models.User, error) {
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

	err = l.db.Create(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (l *AppImpl) GetUser(id uint) (models.User, error) {
	if id == 0 {
		return models.User{
			ID:          0,
			Username:    l.config.MasterUsername,
			DisplayName: "Master",
			IsAdmin:     true,
		}, nil
	}

	var user models.User
	err := l.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (l *AppImpl) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := l.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (l *AppImpl) UpdateUser(user models.User) error {
	return l.db.Save(&user).Error
}

func (l *AppImpl) CreateSession(user models.User, device string) (string, error) {
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

func (l *AppImpl) DeleteSession(token string) error {
	return l.db.Where("token = ?", token).Delete(&models.Session{}).Error
}

func (l *AppImpl) GetSession(token string) (*models.Session, error) {
	var session models.Session
	err := l.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		l.DeleteSession(token)
		return nil, fmt.Errorf("session expired")
	}

	now := time.Now()
	session.LastActivity = now
	session.ExpiresAt = now.Add(time.Hour * 24)
	err = l.db.Save(&session).Error
	if err != nil {
		return nil, err
	}
	err = l.db.Model(&models.User{}).Where("id = ?", session.UserID).Update("last_activity", now).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (l *AppImpl) GetActiveSessions(userID uint) ([]models.Session, error) {
	var sessions []models.Session
	err := l.db.Where("user_id = ? and expires_at > ?", userID, time.Now()).Order("last_activity DESC").Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (l *AppImpl) GetDefaultLanguage() Language {
	return l.config.DefaultLanguage
}

func (l *AppImpl) Auth(r *http.Request) bool {
	if l.config.AuthFunc == nil {
		return true
	}

	return l.config.AuthFunc(r)
}
