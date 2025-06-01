package main

import (
	"log"
	"net/http"

	"sadk.dev/logar"
	logarweb "sadk.dev/logar-web"
	"sadk.dev/logar/logfilter"
	"sadk.dev/logar/proxy"
	"sadk.dev/logar/proxy/consolelogger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	logar.SetTablePrefix("minimal_")
	app, err := logar.New(
		logar.WithAppName("minimal"),
		logar.WithAdminCredentials("admin", "admin"),
		logar.WithDatabase(sqlite.Open("minimal.db")),

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

	err = e.Start(":3000")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
