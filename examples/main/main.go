package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/proxy"
	"github.com/Lexographics/logar/proxy/consolelogger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type RequestIDType struct{}

var RequestID RequestIDType

type UserIDType struct{}

var UserID UserIDType

var name string = "Example application"

func main() {
	needAuth := false

	logger, err := logar.New(
		logar.WithAppName(name),
		logar.WithDatabase("logs.db"),
		logar.AddModel("System Logs", "system-logs"),
		logar.AddModel("User Trace", "user-trace"),

		logar.AddProxy(proxy.NewProxy(
			consolelogger.New(),
			proxy.NewFilter(
				proxy.Not(
					proxy.IsCategory("db-laog"),
				),
			),
		)),

		logar.If(needAuth,
			logar.WithAuth(func(r *http.Request) bool {
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

	db, err := gorm.Open(sqlite.Open("app.db"))
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Get("/logger/*", adaptor.HTTPHandler(logger.ServeHTTP()))

	app.Use(requestid.New())
	app.Use(func(c *fiber.Ctx) error {
		requestid := c.Locals("requestid").(string)
		c.SetUserContext(context.WithValue(c.UserContext(), RequestID, requestid))

		userid := 4
		c.SetUserContext(context.WithValue(c.UserContext(), UserID, userid))

		logger.Trace("user-trace", fiber.Map{
			"requestid": requestid,
			"url":       c.Method() + " " + c.Path(),
			"user_id":   userid,
		}, "request")

		return c.Next()
	})

	afterLogger := func(c *fiber.Ctx) error {
		reqid := c.UserContext().Value(RequestID).(string)
		userid := c.UserContext().Value(UserID).(int)

		status := c.Response().StatusCode()
		body := c.Response().Body()

		var bodyData map[string]any
		err := json.Unmarshal(body, &bodyData)
		if err != nil {
			bodyData = nil
		}

		logger.Trace("user-trace", fiber.Map{
			"requestid": reqid,
			"status":    status,
			"user_id":   userid,
			"body":      bodyData,
		}, "request")
		return nil
	}

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			logger.Error("system-logs", err, "request")
			return nil
		}

		afterLogger(c)
		return nil
	})

	app.Get("/", func(c *fiber.Ctx) error {
		errorText := c.Query("error", "")

		if errorText != "" {
			c.Status(400).JSON(fiber.Map{
				"message": "An error occurred",
				"error":   errorText,
			})
			return c.Next()
		}

		c.Status(200).JSON(fiber.Map{
			"message": "Hello, Logar!",
		})
		return nil
	})

	logger.Info("system-logs", "App Started", "app-start")
	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

type User struct {
	gorm.Model
	Username string
	Email    string
}
