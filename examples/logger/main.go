package main

import (
	"log"
	"net/http"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/logarweb"
	"github.com/Lexographics/logar/logfilter"
	"github.com/Lexographics/logar/proxy"
	"github.com/Lexographics/logar/proxy/consolelogger"
	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	app, err := logar.New(
		logar.WithAppName("logger"),
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
	logger := app.GetLogger()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	e.Any("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("http://localhost:3000", "/logger", app)))
	e.Any("/logger", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/logger/")
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = app.PrepareContext(ctx, logar.Map{"user_id": 123})
			app.AddContextValue(ctx, "request_id", uuid.New().String())

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		ctx := c.Request().Context()

		logger.WithContext(ctx).Info("logs", logar.Map{"message": "Request made to / !"}, "request")
		logger.WithContext(ctx).Info("logs", "Hey! This is just a string message!", "request")
		logger.Info("logs", "This logs without a context, so context values will not be added to the message", "request")

		return c.String(http.StatusOK, "Welcome to the logger example!")
	})

	err = e.Start(":3000")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
