package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sadk.dev/logar"
	logarweb "sadk.dev/logar-web"
	"sadk.dev/logar/gormlogger"
	"sadk.dev/logar/logfilter"
	"sadk.dev/logar/proxy"
	"sadk.dev/logar/proxy/consolelogger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var name string = "Example application"

func main() {
	app, err := logar.New(
		logar.WithAppName(name),
		logar.WithDatabase(sqlite.Open("logs.db")),
		logar.WithDefaultLanguage(logar.English),

		logar.WithAdminCredentials("admin", "password"),

		logar.AddModel("User Logs", "user-trace", "fa-solid fa-user"),
		logar.AddModel("System Logs", "system-logs", "fa-solid fa-gear"),
		logar.AddModel("All Logs", "__all__", "fa-solid fa-file-lines"),

		// Example Actions
		logar.WithAction("Server/Status", "Get server status", func() map[string]interface{} {
			return map[string]interface{}{
				"status":    "healthy",
				"timestamp": time.Now().Format(time.RFC3339),
				"uptime":    time.Since(time.Now().Add(-time.Second * 10)).String(),
			}
		}),
		logar.WithAction("Server/Info", "Get server information", func() map[string]interface{} {
			return map[string]interface{}{
				"name":      name,
				"version":   "1.0.0",
				"startTime": time.Now().Add(-time.Second * 10).Format(time.RFC3339),
			}
		}),

		logar.WithAction("Notification/Send", "Send a notification to all users", func(message string) string {
			return "Notification sent: " + message
		}),

		logar.AddProxy(proxy.NewProxy(
			consolelogger.New(),
			logfilter.NewFilter(
				logfilter.Not(
					logfilter.IsCategory("db-log"),
				),
			),
		)),
	)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{
		Logger: gormlogger.New(app, "user-trace", "db-log", 1),
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
		RequestIDHandler: func(c echo.Context, requestid string) {
			c.Set("requestid", requestid)
		},
	}))

	e.Any("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("http://localhost:3000", "/logger", app)))
	e.Any("/logger", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/logger/")
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/logger") {
				return next(c)
			}

			requestid := c.Get("requestid").(string)
			userid := 4
			c.Set("userid", userid)

			app.GetLogger().Trace("user-trace", map[string]interface{}{
				"requestid": requestid,
				"url":       c.Request().Method + " " + c.Request().URL.Path,
				"user_id":   userid,
			}, "request")

			start := time.Now()
			err := next(c)
			duration := time.Since(start).String()

			status := c.Response().Status

			if err != nil {
				fmt.Printf("err: %v\n", err)
				app.GetLogger().Error("user-trace", map[string]interface{}{
					"requestid": requestid,
					"status":    status,
					"user_id":   userid,
					"error":     err,
					"duration":  duration,
				}, "request")
				return err
			}

			app.GetLogger().Trace("user-trace", map[string]interface{}{
				"requestid": requestid,
				"status":    status,
				"user_id":   userid,
				"duration":  duration,
			}, "request")

			return nil
		}
	})

	e.GET("/", func(c echo.Context) error {
		errorText := c.QueryParam("error")

		if errorText != "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "An error occurred",
				"error":   errorText,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
		})
	})

	e.GET("/timer", func(c echo.Context) error {
		timer := app.GetLogger().NewTimer()
		defer timer.Log("user-trace", "Handler done", "handler")

		time.Sleep(time.Second * 1)

		timer.Log("user-trace", "Do some stuff", "handler")

		time.Sleep(time.Millisecond * 1500)

		timer.Log("user-trace", "Do some more stuff", "handler")

		return c.String(http.StatusOK, "Timer test!")
	})

	app.GetLogger().Info("system-logs", logar.Map{"message": "App Started"}, "app-start")

	err = e.Start(":3000")
	if err != nil {
		panic(err)
	}
}

type User struct {
	gorm.Model
	Username string
	Email    string
}
