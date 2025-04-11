package gormlogger

import (
	"fmt"

	"github.com/Lexographics/logar"
	"github.com/Lexographics/logar/internal/domain/models"
)

type writer struct {
	lg       *logar.Logger
	severity models.Severity
	model    string
	category string
}

func newWriter(logger *logar.Logger, model, category string) *writer {
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

	w.lg.Print(w.model, msg, w.category, w.severity)
}
