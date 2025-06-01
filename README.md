# Logar

A lightweight, flexible management library for Go applications, providing logging, analytics, server actions, and more.

## Features

- Simple, intuitive API for all modules
- **Logging**:
  - Multiple log levels (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)
  - Output to console, file, or custom writers via proxies
  - Context-aware logging
- **Server Actions**:
  - Define and trigger custom server-side functions remotely with strongly-typed parameters.
- **Analytics**:
  - Track service metrics such as request latency, error rates, distribution of OS, browser and referers, and bandwidth.
- **Feature Flags**:
  - Toggle and control application features through runtime-configurable flags with conditional logic.
- **Web UI**:
  - View and manage logs through a web interface
  - Execute server actions remotely
  - View analytics dashboards
  - Manage feature flags and conditions
- Context-aware operations

## Installation

```bash
go get sadk.dev/logar
```

## Quick Start

Logar provides a suite of tools. Here's a basic logging example:

```go
package main

import (
  "sadk.dev/logar"
  // Import other necessary logar sub-packages for specific features
  // e.g., "sadk.dev/logar/logarweb" for the web UI
)

func main() {
  // Create a new Logar application instance
  // The New() function accepts various configuration options for logging,
  // analytics, actions, web panel, etc.
  app, err := logar.New(
    logar.WithAppName("My Awesome App"), // Sets the application name
    logar.AddModel("System Logs", "system-logs"),
    logar.AddModel("User Activity", "user-activity"),

    // Add a console proxy to output all logs to terminal
    // logar.AddProxy(proxy.NewProxy(
    //   consolelogger.New(),
    //   proxy.NewFilter(),
    // )),
    // For the web UI, actions, and analytics, additional setup is required.
    // See the examples directory for detailed demonstrations.
  )
  if err != nil {
    // Handle error
  }

  // Basic logging
  app.GetLogger().Info("system-logs", "App Started. No errors", "startup")

  // To explore server actions, analytics, feature flags, and the web UI,
  // please refer to the detailed examples in the 'examples/' directory
  // in the project repository.
}
```

## Configuration

Logar is configured using functional options with the `logar.New()` constructor.

```go
// Configure Logar with various options
app, err := logar.New(
  logar.WithAppName("My App"),
  logar.WithDatabase("logs.db"), // For persisting logs and analytics data
  logar.WithAdminCredentials("admin", "securepassword"), // For web UI authentication

  // Logging specific configurations
  logar.AddModel("System Events", "system-events"),

  // Action specific configurations
  logar.WithAction("Server/Echo", "Description of action", func(message string) string {
    // Action logic here
    return "Result: " + message
  }),

  // Analytics can be automatically collected via middleware
  // Or events can be registered manually:
  // app.GetAnalytics().RegisterEvent(...)


  // Forwards logs to the proxy based on given filters (example)
  // logar.AddProxy(proxy.NewProxy(
  //   consolelogger.New(),
  //   proxy.NewFilter(),
  // )),

  // Conditional options
  logar.If(env == "test",
    logar.WithDatabase("test.db"),
  ),

  logar.IfElse(env == "prod",
    logar.WithDatabase("prod.db"),
    logar.WithDatabase("dev.db"),
  )
)
if err != nil {
  // Handle error
}

// Access different modules:
logger := app.GetLogger()
analyticsEngine := app.GetAnalytics()
actionManager := app.GetActionManager()
featureFlagManager := app.GetFeatureFlags()

// Start the web server (typically using net/http or a framework like Echo)
// and integrate logarweb.ServeHTTP for the Logar web panel.
// e.g., e.Any("/logar/*", echo.WrapHandler(logarweb.ServeHTTP("http://localhost:3000", "/logar", app)))
```

## Advanced Usage & Examples

For detailed examples covering logging, the web UI, server actions, analytics, feature flags, and integrations with web frameworks like Echo, see the `examples/` directory in the GitHub repository.

See the [documentation](https://godoc.org/sadk.dev/logar) for more API details.

## License

MIT

## Contributing

If you'd like to contribute to the project, please open an issue on GitHub to discuss your ideas or report bugs.

**Warning:** This project does not currently adhere strictly to semver, and API compatibility is not guaranteed between minor and patch updates.

## Known Issues

- ...
