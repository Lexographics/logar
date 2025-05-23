package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/logarweb"
	"github.com/Lexographics/logar/logfilter"
	"github.com/Lexographics/logar/proxy"
	"github.com/Lexographics/logar/proxy/consolelogger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	app, err := logar.New(
		logar.WithAppName("featureflags"),
		logar.WithAdminCredentials("admin", "admin"),

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
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = app.PrepareContext(ctx, logar.Map{"user_id": 123})

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		ctx := c.Request().Context()
		enabled, err := app.GetFeatureFlags().HasFeatureFlag(ctx, "experimental-feature")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error checking feature flag: "+err.Error())
		}

		return c.String(http.StatusOK, fmt.Sprintf("Experimental feature flag enabled: %v", enabled))
	})

	e.Any("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("http://localhost:3000", "/logger", app)))
	e.Any("/logger", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/logger/")
	})

	err = e.Start(":3000")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
