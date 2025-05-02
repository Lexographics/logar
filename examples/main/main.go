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
		logar.WithAppName(name),
		logar.WithDatabase("logs.db"),
		logar.WithDefaultLanguage(logar.English),
		logar.AddModel("System Logs", "system-logs"),
		logar.AddModel("User Trace", "user-trace"),
		logar.AddModel("Test Logs", "test-logs"),
		logar.AddModel("Test1", "test1"),
		logar.AddModel("Test2", "test2"),
		logar.AddModel("Test3", "test3"),
		logar.AddModel("Test4", "test4"),
		// logar.AddModel("Test5", "test5"),
		// logar.AddModel("Test6", "test6"),
		// logar.AddModel("Test7", "test7"),
		// logar.AddModel("Test6", "test6"),
		// logar.AddModel("Test8", "test8"),
		// logar.AddModel("Test9", "test9"),
		// logar.AddModel("Test10", "test10"),
		logar.AddModel("All Logs", "__all__"),

		logar.WithMasterCredentials("username", "password"),

		// Example Actions
		logar.WithAction("Server/Test", "Test action", func() (string, int, string) {
			return "test value 1", 123, "test value 3"
		}),
		logar.WithAction("Server/Time", "Get current time", func() string {
			return time.Now().Format(time.RFC3339)
		}),
		logar.WithAction("Server/Ping", "Ping the server", func() string {
			return "pong"
		}),
		logar.WithAction("Math/Add", "Add two numbers", func(a, b float64) float64 {
			return a + b
		}),
		logar.WithAction("Math/Sub", "Subtract two numbers", func(a, b float64) float64 {
			return a - b
		}),
		logar.WithAction("Math/Mul", "Multiply two numbers", func(a, b float64) float64 {
			return a * b
		}),
		logar.WithAction("Math/Div", "Divide two numbers", func(a, b float64) float64 {
			return a / b
		}),
		logar.WithAction("Math/Advanced/Pow", "Power of two numbers", func(a, b float64) float64 {
			return math.Pow(a, b)
		}),
		logar.WithAction("Math/Advanced/Sqrt", "Square root of a number", func(a float64) float64 {
			return math.Sqrt(a)
		}),

		logar.WithAction("Greet", "Greet a user", func(name string) string {
			return fmt.Sprintf("Hello, %s!", name)
		}),
		logar.WithAction("Concat", "Concatenate strings", func(args []any) string {
			str := ""
			for _, arg := range args {
				str += fmt.Sprintf("%v", arg)
			}
			return str
		}),

		logar.AddProxy(proxy.NewProxy(
			consolelogger.New(),
			logfilter.NewFilter(
				logfilter.Not(
					logfilter.IsCategory("db-log"),
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

	app.Any("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("http://localhost:3000", "/logger", logger)))
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
			"status": "ok",
		})
	})

	app.GET("/timer", func(c echo.Context) error {
		timer := logger.NewTimer()
		defer timer.Log("user-trace", "Handler done", "handler")

		time.Sleep(time.Second * 1)

		timer.Log("user-trace", "Do some stuff", "handler")

		time.Sleep(time.Millisecond * 1500)

		timer.Log("user-trace", "Do some more stuff", "handler")

		return c.String(http.StatusOK, "Timer test!")
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
