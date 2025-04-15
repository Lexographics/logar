package config

import (
	"net/http"

	"github.com/Lexographics/logar/proxy"
)

type Config struct {
	AppName        string
	Database       string
	AutoMigrate    bool
	RequireAuth    bool
	AuthFunc       AuthFunc
	Models         LogModels
	Proxies        []proxy.Proxy
	Actions        ActionMap
	MasterUsername string
	MasterPassword string
}

type LogModel struct {
	DisplayName string `json:"displayName"`
	Identifier  string `json:"identifier"`
}
type LogModels []LogModel

type AuthFunc func(r *http.Request) bool
type ConfigOpt func(*Config)

type ActionMap []Action

type Action struct {
	Path        string
	Func        interface{}
	Description string
}

func WithAppName(appName string) ConfigOpt {
	return func(cfg *Config) {
		cfg.AppName = appName
	}
}

func WithDatabase(database string) ConfigOpt {
	return func(cfg *Config) {
		cfg.Database = database
	}
}

func EnableAutoMigrate(cfg *Config) {
	cfg.AutoMigrate = true
}

func DisableAutoMigrate(cfg *Config) {
	cfg.AutoMigrate = false
}

func WithAuth(authFunc AuthFunc) ConfigOpt {
	return func(cfg *Config) {
		cfg.RequireAuth = true
		cfg.AuthFunc = authFunc
	}
}

func AddModel(displayName, modelId string) ConfigOpt {
	return func(cfg *Config) {
		cfg.Models = append(cfg.Models, LogModel{
			DisplayName: displayName,
			Identifier:  modelId,
		})
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

func WithMasterCredentials(username, password string) ConfigOpt {
	return func(cfg *Config) {
		cfg.MasterUsername = username
		cfg.MasterPassword = password
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
