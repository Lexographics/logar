package main

import (
	"context"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sadk.dev/logar"
	logarweb "sadk.dev/logar-web"
	"sadk.dev/logar/gormlogger"
	"sadk.dev/logar/logfilter"
	"sadk.dev/logar/proxy"
	"sadk.dev/logar/proxy/consolelogger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app, err := logar.New(
		logar.WithAppName("gormlogger"),
		logar.WithAdminCredentials("admin", "admin"),
		logar.WithDatabase(sqlite.Open("logs.db")),

		logar.AddModel("User Trace", "user-trace", "fa-solid fa-users"),
		logar.AddModel("Logs", "logs", "fa-solid fa-file-lines"),

		logar.AddProxy(proxy.NewProxy(
			consolelogger.New(),
			logfilter.NewFilter(),
		)),
	)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	e.Any("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("http://localhost:3000", "/logger", app)))
	e.Any("/logger", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/logger/")
	})

	testDb, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		log.Fatal(err)
	}
	testDb.Logger = gormlogger.New(app, "logs", "db-log", 0)
	ctx := app.PrepareContext(context.Background(), logar.Map{"user_id": 1})
	testDb.AutoMigrate(&User{})

	// this will not print
	testDb.WithContext(ctx).Create(&User{Name: "John"})

	// this will print
	testDb.WithContext(ctx).Debug().Create(&User{Name: "Jane"})

	// this will fail
	testDb.WithContext(ctx).Create(&User{Name: "John"})

	err = e.Start(":3000")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

type User struct {
	gorm.Model
	Name string `gorm:"unique"`
}
