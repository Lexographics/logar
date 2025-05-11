package main

import (
	"log"
	"net/http"
	"time"

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
		logar.WithAppName("minimal"),
		logar.WithAdminCredentials("admin", "admin"),

		logar.AddModel("Logs", "logs", "fa-solid fa-file-lines"),

		logar.AddProxy(proxy.NewProxy(
			consolelogger.New(),
			logfilter.NewFilter(
				logfilter.Not(
					logfilter.MessageContains(`"iteration":5`),
				),
			),
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

	ticker := time.NewTicker(time.Second * 2)
	iteration := 0
	go func() {
		for range ticker.C {
			app.Info("logs", logar.Map{
				"iteration": iteration + 1,
				"message":   "every fifth iteration will be ignored",
			}, "app-log")
			iteration++
			iteration = iteration % 10
		}
	}()

	err = e.Start(":3000")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
