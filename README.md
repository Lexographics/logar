# Logar

A lightweight, flexible logging library for Go applications.

## Features

- Simple, intuitive API
- Multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL, TRACE)
- Output to console, file, or custom writers
- Context-aware logging

## Installation

```bash
go get github.com/Lexographics/logar
```

## Quick Start

```go
package main

import (
  "github.com/Lexographics/logar"
)

func main() {
  // Create a new logger
  logger, err := logar.New(
    logar.AddModel("System Logs", "system-logs"),
    logar.AddModel("User Activity", "user-activity"),

    // Add a console proxy to output all logs to terminal
    logar.AddProxy(proxy.NewProxy(
      consolelogger.New(),
      proxy.NewFilter(),
    )),
  )

  // Basic logging
  logger.Info("system-logs", "App Started. No errors", "startup")
}
```

## Configuration

```go
// Configure logger with options
logger := logar.New(
  
  logar.WithAppName("app name")
  
  logar.WithDatabase("logs.db")
  
  logar.WithAuth(func(r *http.Request) bool {
    username, password, ok := r.BasicAuth()
    if ok && username == "admin" && password == "password" {
      return true
    }
    return false
  }),
  
  // Add model with display name and model identifier
  logar.AddModel("System Logs", "system-logs"),

  // Forwards logs to the proxy based on given filters
  logar.AddProxy(proxy.NewProxy(
    consolelogger.New(),
    proxy.NewFilter(),
  )),

  // Conditional options
  logar.If(env == "test", 
    logar.WithDatabase("test.db")
    ...
  )

  logar.IfElse(env == "test",
    logar.WithDatabase("test.db"),
    logar.WithDatabase("prod.db")
    ...
  )
)
```

## Advanced Usage

See the [documentation](https://godoc.org/github.com/Lexographics/logar) for more examples and advanced usage.

## License

MIT

## Contributing

If you'd like to contribute to the project, please open an issue on GitHub to discuss your ideas or report bugs.

**Warning:** This project does not currently adhere strictly to semver, and API compatibility is not guaranteed between minor and patch updates.

## Known Issues

- Infinite pagination on web panel does not include filters in requests
