package gormlogger

import (
	"fmt"

	"sadk.dev/logar"
	"sadk.dev/logar/models"
)

type writer struct {
	lg       logar.App
	severity models.Severity
	model    logar.Model
	category string
}

func newWriter(logger logar.App, model logar.Model, category string) *writer {
	return &writer{
		lg:       logger,
		severity: models.Severity_Info,
		model:    model,
		category: category,
	}
}

func (w *writer) SetSeverity(severity models.Severity) {
	w.severity = severity
}

func (w *writer) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf("\n"+""+format+"\n", args...)

	w.lg.GetLogger().Print(w.model, msg, w.category, w.severity)
}
