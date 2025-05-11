package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/logarweb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mileusna/useragent"
)

func main() {
	app, err := logar.New(
		logar.WithAppName("analytics"),
		logar.WithAdminCredentials("admin", "admin"),
	)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/logger") || strings.HasSuffix(c.Request().URL.Path, "/sse") {
				return next(c)
			}

			res := c.Response()

			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			statusCode := http.StatusOK
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				if ok {
					statusCode = he.Code
				} else {
					statusCode = http.StatusInternalServerError
				}
			} else {
				statusCode = c.Response().Status
			}

			userID := "123"
			contentLength, _ := strconv.ParseInt(c.Request().Header.Get(echo.HeaderContentLength), 10, 64)
			ua := useragent.Parse(c.Request().UserAgent())

			app.RegisterRequest(logar.RequestLog{
				Timestamp:  time.Now(),
				VisitorID:  userID,
				Instance:   "eu-east-1",
				Path:       c.Request().URL.Path,
				Latency:    duration,
				StatusCode: statusCode,
				UserAgent:  c.Request().UserAgent(),
				OS:         ua.OS,
				Browser:    ua.Name,
				Referer:    c.Request().Referer(),
				BytesSent:  res.Size,
				BytesRecv:  contentLength,
			})

			return err
		}
	})

	e.Any("/logger/*", echo.WrapHandler(logarweb.ServeHTTP("http://localhost:3000", "/logger", app)))
	e.Any("/logger", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/logger/")
	})

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/page/1")
	})
	for i := 1; i <= 10; i++ {
		e.GET(fmt.Sprintf("/page/%d", i), func(c echo.Context) error {
			time.Sleep(time.Duration(10+rand.Int()%190) * time.Millisecond)

			return c.HTML(http.StatusOK, fmt.Sprintf(`
<html>
<body>
	<h1>Page %d</h1>
	
	<p> Other Pages: </p>
	<ul>
		<li><a href="/page/1">Page 1</a></li>
		<li><a href="/page/2">Page 2</a></li>
		<li><a href="/page/3">Page 3</a></li>
		<li><a href="/page/4">Page 4</a></li>
		<li><a href="/page/5">Page 5</a></li>
		<li><a href="/page/6">Page 6</a></li>
		<li><a href="/page/7">Page 7</a></li>
		<li><a href="/page/8">Page 8</a></li>
		<li><a href="/page/9">Page 9</a></li>
		<li><a href="/page/10">Page 10</a></li>
	</ul>
</body>
</html>
`, i))
		})
	}

	err = e.Start(":3000")
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
