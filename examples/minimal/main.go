package main

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/gormlogger"
	"github.com/Lexographics/logar/logarweb"
	"github.com/Lexographics/logar/logfilter"
	"github.com/Lexographics/logar/proxy"
	"github.com/Lexographics/logar/proxy/consolelogger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomType int32

type AnotherCustomType float64

func main() {

	logger, err := logar.New(
		logar.WithAppName("minimal"),
		logar.WithDatabase("logs.db"),

		logar.AddModel("User Trace", "user-trace", "fa-solid fa-users"),
		logar.AddModel("Logs", "logs", "fa-solid fa-file-lines"),
		logar.WithMasterCredentials("username", "password"),

		logar.WithAction("Server/Test", "Test action", func() (string, int, string) {
			return "test value 1", 123, "test value 3"
		}),
		logar.WithAction("Server/Time", "Get current time", func() string {
			return time.Now().Format(time.RFC3339)
		}),
		logar.WithAction("Server/Ping", "Ping the server", func(duration time.Duration, timestamp time.Time, test string, test2 CustomType) string {
			return "duration: " + fmt.Sprint(duration) + " timestamp: " + fmt.Sprint(timestamp.String()) + " test: " + test + " test2: " + fmt.Sprint(test2)
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

	logger.SetTypeKind(reflect.TypeOf(CustomType(0)), logar.TypeKind_Int)
	logger.SetTypeKind(reflect.TypeOf(AnotherCustomType(0)), logar.TypeKind_Float)

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

	app.Any("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("/logger", logger)))
	app.Any("/logger", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/logger/")
	})

	logger.Info("logs", logar.Map{"message": "App Started"}, "app-start")

	err = app.Start(":3000")
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

type User struct {
	gorm.Model
	Username string
	Email    string
}
