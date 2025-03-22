package config

import (
	"net/http"

	"github.com/Lexographics/logar/proxy"
)

type Config struct {
	AppName     string
	Database    string
	AutoMigrate bool
	RequireAuth bool
	AuthFunc    AuthFunc
	Models      LogModels
	Proxies     []proxy.Proxy
}

type LogModel struct {
	DisplayName string
	ModelId     string
}
type LogModels []LogModel

type AuthFunc func(r *http.Request) bool
type ConfigOpt func(*Config)

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
			ModelId:     modelId,
		})
	}
}

func AddProxy(proxy proxy.Proxy) ConfigOpt {
	return func(cfg *Config) {
		cfg.Proxies = append(cfg.Proxies, proxy)
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
