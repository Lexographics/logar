package main

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/gormlogger"
	"github.com/Lexographics/logar/logarweb"
	"github.com/Lexographics/logar/logfilter"
	"github.com/Lexographics/logar/options/config"
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
		config.AddModel("Test Logs", "test-logs"),
		config.AddModel("Test1", "test1"),
		config.AddModel("Test2", "test2"),
		config.AddModel("Test3", "test3"),
		config.AddModel("Test4", "test4"),
		// config.AddModel("Test5", "test5"),
		// config.AddModel("Test6", "test6"),
		// config.AddModel("Test7", "test7"),
		// config.AddModel("Test6", "test6"),
		// config.AddModel("Test8", "test8"),
		// config.AddModel("Test9", "test9"),
		// config.AddModel("Test10", "test10"),
		config.AddModel("All Logs", "__all__"),

		config.WithMasterCredentials("username", "password"),

		// Example Actions
		config.WithAction("Server/Test", "Test action", func() (string, int, string) {
			return "test value 1", 123, "test value 3"
		}),
		config.WithAction("Server/Time", "Get current time", func() string {
			return time.Now().Format(time.RFC3339)
		}),
		config.WithAction("Server/Ping", "Ping the server", func() string {
			return "pong"
		}),
		config.WithAction("Math/Add", "Add two numbers", func(a, b float64) float64 {
			return a + b
		}),
		config.WithAction("Math/Sub", "Subtract two numbers", func(a, b float64) float64 {
			return a - b
		}),
		config.WithAction("Math/Mul", "Multiply two numbers", func(a, b float64) float64 {
			return a * b
		}),
		config.WithAction("Math/Div", "Divide two numbers", func(a, b float64) float64 {
			return a / b
		}),
		config.WithAction("Math/Advanced/Pow", "Power of two numbers", func(a, b float64) float64 {
			return math.Pow(a, b)
		}),
		config.WithAction("Math/Advanced/Sqrt", "Square root of a number", func(a float64) float64 {
			return math.Sqrt(a)
		}),

		config.WithAction("Greet", "Greet a user", func(name string) string {
			return fmt.Sprintf("Hello, %s!", name)
		}),
		config.WithAction("Concat", "Concatenate strings", func(args []any) string {
			str := ""
			for _, arg := range args {
				str += fmt.Sprintf("%v", arg)
			}
			return str
		}),

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

	res, err := logger.InvokeAction("Math/Div", 4.0, 3.0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("add res: %v\n", res)

	res, err = logger.InvokeAction("Server/Ping")
	if err != nil {
		panic(err)
	}

	fmt.Printf("ping res: %v\n", res)

	res, err = logger.InvokeAction("Greet", "John")
	if err != nil {
		panic(err)
	}

	fmt.Printf("greet res: %v\n", res)

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

	app.Any("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("/logger", logger)))
	app.Any("/logger", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/logger/")
	})

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

			start := time.Now()
			err := next(c)
			duration := time.Since(start).String()

			status := c.Response().Status

			if err != nil {
				fmt.Printf("err: %v\n", err)
				logger.Error("user-trace", map[string]interface{}{
					"requestid": requestid,
					"status":    status,
					"user_id":   userid,
					"error":     err,
					"duration":  duration,
				}, "request")
				return err
			}

			logger.Trace("user-trace", map[string]interface{}{
				"requestid": requestid,
				"status":    status,
				"user_id":   userid,
				"duration":  duration,
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

	logger.Info("system-logs", logar.Map{"message": "App Started"}, "app-start")

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
