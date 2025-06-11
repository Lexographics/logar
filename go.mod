module sadk.dev/logar

go 1.23.3

require (
	github.com/expr-lang/expr v1.17.2
	github.com/google/uuid v1.6.0
	github.com/labstack/echo/v4 v4.13.3
	golang.org/x/crypto v0.36.0
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
	sadk.dev/logar-web v0.0.0-00010101000000-000000000000
)

replace sadk.dev/logar-web => ../logar-web

require (
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/time v0.8.0 // indirect
)

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/mileusna/useragent v1.3.5
	golang.org/x/text v0.23.0 // indirect
)
