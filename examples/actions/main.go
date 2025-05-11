package main

import (
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/logarweb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomIntType int32

type CustomFloatType float64

type CustomDuration time.Duration

func main() {

	app, err := logar.New(
		logar.WithAppName("actions"),
		logar.WithAdminCredentials("admin", "admin"),

		logar.WithAction("Server/Echo", "Echo the server", func(message string) string {
			return message
		}),
		logar.WithAction("Server/Time", "Get time with offset", func(offset CustomDuration) string {
			return time.Now().Add(time.Duration(offset)).Format(time.RFC3339)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	app.SetTypeKind(reflect.TypeOf(CustomIntType(0)), logar.TypeKind_Int)
	app.SetTypeKind(reflect.TypeOf(CustomFloatType(0)), logar.TypeKind_Float)
	app.SetTypeKind(reflect.TypeOf(CustomDuration(0)), logar.TypeKind_Duration)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	e.Any("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("http://localhost:3000", "/logger", app)))
	e.Any("/logger", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/logger/")
	})

	err = e.Start(":3000")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
