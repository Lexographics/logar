package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"sadk.dev/logar"
	logarweb "sadk.dev/logar-web"
	"sadk.dev/logar/logfilter"
	"sadk.dev/logar/proxy"
	"sadk.dev/logar/proxy/consolelogger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	appName       = flag.String("app-name", "minimal", "Application name")
	adminUsername = flag.String("admin-username", "admin", "Admin username")
	adminPassword = flag.String("admin-password", "admin", "Admin password")
	dbPath        = flag.String("db-path", "logs.db", "Database path")
	serverAddr    = flag.String("server-addr", ":3000", "Server address")
	apiURL        = flag.String("api-url", "http://localhost:3000", "API URL for web client")
	basePath      = flag.String("base-path", "", "Base path for web client")
)

func main() {
	flag.Parse()

	app, err := logar.New(
		logar.WithAppName(*appName),
		logar.WithAdminCredentials(*adminUsername, *adminPassword),
		logar.WithDatabase(sqlite.Open(*dbPath)),

		logar.AddModel("User Trace", "user-trace", "fa-solid fa-users"),
		logar.AddModel("Logs", "logs", "fa-solid fa-file-lines"),

		logar.AddProxy(proxy.NewProxy(
			consolelogger.New(),
			logfilter.NewFilter(),
		)),

		logar.WithAction("Server/Time", "Get current time", func() string {
			return time.Now().Format(time.RFC3339)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	e.GET("/manifest.json", func(c echo.Context) error {
		base := *basePath
		if base == "" {
			base = "/"
		}

		data := `{
  "name": "Logar Management",
  "short_name": "Logar",
  "description": "Logar is a tool for logging and monitoring your application.",
  "start_url": "` + base + `",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#0a111f",
  "orientation": "portrait-primary",
  "lang": "en-US",
  "icons": [
    {
      "src": "` + base + `/icons/icon-192x192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "` + base + `/icons/icon-512x512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}`
		var manifest map[string]any
		err := json.Unmarshal([]byte(data), &manifest)
		if err != nil {
			return c.JSON(500, map[string]any{"error": err.Error()})
		}
		return c.JSON(200, manifest)
	})
	e.Any("*", echo.WrapHandler(logarweb.ServeHTTP(*apiURL, *basePath, app)))

	err = e.Start(*serverAddr)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
