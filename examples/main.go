package main

import (
	"flag"
	"net/http"

	"github.com/Lexographics/logar"
)

var name string = "Example application"

func main() {
	number := flag.Int("number", 0, "A number")
	flag.Parse()

	logger, err := logar.New(
		logar.WithAppName(name),
		logar.AddModel("System Logs", "system-logs"),
	)
	if err != nil {
		panic(err)
	}

	if *number == 0 {
		logger.Warn("system-logs", "No number provided", "config")
	} else if *number == 1 {
		logger.Error("system-logs", "Number 1 is not allowed", "config")
	}

	logger.Info("system-logs", "Configuration loaded", "config")

	mux := logger.ServeHTTP()
	http.ListenAndServe(":8080", mux)
}
