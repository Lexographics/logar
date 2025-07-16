package logar

import (
	"net/http"
	"time"

	"gorm.io/gorm"
	"sadk.dev/logar/proxy"
)

type Config struct {
	AppName         string
	Database        gorm.Dialector
	RequireAuth     bool
	AuthFunc        AuthFunc
	Models          LogModels
	Proxies         []proxy.Proxy
	Actions         Actions
	AdminUsername   string
	AdminPassword   string
	DefaultLanguage Language
	WebPanelConfig  WebPanelConfig
	SSEEnabled      bool
}

type LogModel struct {
	DisplayName string `json:"displayName"`
	Identifier  Model  `json:"identifier"`
	Icon        string `json:"icon"` // FontAwesome icon name. default: "fa-solid fa-cube"
}
type LogModels []LogModel

type AuthFunc func(r *http.Request) bool
type ConfigOpt func(*Config)

type Actions []Action

type Action struct {
	Path        string
	Func        interface{}
	Description string
}

type WebPanelConfig struct {
	SessionDuration time.Duration
}

type WebPanelConfigOpt func(*WebPanelConfig)

func WithSessionDuration(duration time.Duration) WebPanelConfigOpt {
	return func(cfg *WebPanelConfig) {
		cfg.SessionDuration = duration
	}
}

func WithWebPanelConfig(opts ...WebPanelConfigOpt) ConfigOpt {
	return func(cfg *Config) {
		webPanelConfig := defaultWebPanelConfig
		for _, opt := range opts {
			opt(&webPanelConfig)
		}

		cfg.WebPanelConfig = webPanelConfig
	}
}

func WithAppName(appName string) ConfigOpt {
	return func(cfg *Config) {
		cfg.AppName = appName
	}
}

func WithDatabase(database gorm.Dialector) ConfigOpt {
	return func(cfg *Config) {
		cfg.Database = database
	}
}

func WithAuth(authFunc AuthFunc) ConfigOpt {
	return func(cfg *Config) {
		cfg.RequireAuth = true
		cfg.AuthFunc = authFunc
	}
}

func AddModel(displayName string, modelId Model, icon ...string) ConfigOpt {
	return func(cfg *Config) {
		ico := "fa-solid fa-cube"
		if len(icon) > 0 {
			ico = icon[0]
		}
		cfg.Models = append(cfg.Models, LogModel{
			DisplayName: displayName,
			Identifier:  modelId,
			Icon:        ico,
		})
	}
}

func SetModels(models LogModels) ConfigOpt {
	return func(cfg *Config) {
		cfg.Models = models
	}
}

func AddProxy(proxy proxy.Proxy) ConfigOpt {
	return func(cfg *Config) {
		cfg.Proxies = append(cfg.Proxies, proxy)
	}
}

func WithAction(path string, description string, action interface{}) ConfigOpt {
	return func(cfg *Config) {
		cfg.Actions = append(cfg.Actions, Action{
			Path:        path,
			Func:        action,
			Description: description,
		})
	}
}

func WithAdminCredentials(username, password string) ConfigOpt {
	return func(cfg *Config) {
		cfg.AdminUsername = username
		cfg.AdminPassword = password
	}
}

func WithDefaultLanguage(language Language) ConfigOpt {
	return func(cfg *Config) {
		cfg.DefaultLanguage = language
	}
}

func WithSSEEnabled(enabled bool) ConfigOpt {
	return func(cfg *Config) {
		cfg.SSEEnabled = enabled
	}
}

func Combine(opts ...ConfigOpt) ConfigOpt {
	return func(cfg *Config) {
		for _, opt := range opts {
			opt(cfg)
		}
	}
}

func If(condition bool, opts ...ConfigOpt) ConfigOpt {
	if condition {
		return Combine(opts...)
	}
	return func(cfg *Config) {}
}

func IfElse(condition bool, ifOpts ConfigOpt, elseOpts ...ConfigOpt) ConfigOpt {
	if condition {
		return ifOpts
	}
	return Combine(elseOpts...)
}
