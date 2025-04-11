package main

import (
	"net/http"
	"strings"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/gormlogger"
	"github.com/Lexographics/logar/internal/logarweb"
	"github.com/Lexographics/logar/internal/logfilter"
	"github.com/Lexographics/logar/internal/options/config"
	"github.com/Lexographics/logar/proxy"
	"github.com/Lexographics/logar/proxy/consolelogger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var name string = "Example application"

func main() {
	needAuth := false

	logger, err := logar.New(
		config.WithAppName(name),
		config.WithDatabase("logs.db"),
		config.AddModel("System Logs", "system-logs"),
		config.AddModel("User Trace", "user-trace"),

		config.AddProxy(proxy.NewProxy(
			consolelogger.New(),
			logfilter.NewFilter(
				logfilter.Not(
					logfilter.IsCategory("db-log"),
				),
			),
		)),

		config.If(needAuth,
			config.WithAuth(func(r *http.Request) bool {
				username, password, ok := r.BasicAuth()
				if ok && username == "admin" && password == "password" {
					return true
				}

				return false
			}),
		),
	)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{
		Logger: gormlogger.New(logger, "user-trace", "db-log", 1),
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	logger.Info("system-logs", "User registered https://example.com", "user-register-telegram")

	app := echo.New()
	app.Use(middleware.CORS())
	app.Use(middleware.Recover())
	app.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
		RequestIDHandler: func(c echo.Context, requestid string) {
			c.Set("requestid", requestid)
		},
	}))

	app.GET("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("/logger", logger)))

	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/logger") {
				return next(c)
			}

			requestid := c.Get("requestid").(string)
			userid := 4
			c.Set("userid", userid)

			logger.Trace("user-trace", map[string]interface{}{
				"requestid": requestid,
				"url":       c.Request().Method + " " + c.Request().URL.Path,
				"user_id":   userid,
			}, "request")

			err := next(c)
			if err != nil {
				logger.Error("user-trace", err, "request")
				return err
			}

			status := c.Response().Status

			logger.Trace("user-trace", map[string]interface{}{
				"requestid": requestid,
				"status":    status,
				"user_id":   userid,
			}, "request")

			return nil
		}
	})

	app.GET("/", func(c echo.Context) error {
		errorText := c.QueryParam("error")

		if errorText != "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "An error occurred",
				"error":   errorText,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello, Logar!",
		})
	})

	logger.Info("system-logs", "App Started https://example.com", "app-start")
	logger.Info("system-logs", "App Started http://example.com", "app-start")

	logger.Info("system-logs", logar.Map{"hello": "this"}, "app-start")

	err = app.Start(":3000")
	if err != nil {
		panic(err)
	}
}

type User struct {
	gorm.Model
	Username string
	Email    string
}
