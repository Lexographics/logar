package gormlogger

import (
	"context"
	"fmt"

	"sadk.dev/logar"
	"sadk.dev/logar/models"
)

type writer struct {
	lg       logar.App
	model    logar.Model
	category string
}

func newWriter(logger logar.App, model logar.Model, category string) *writer {
	return &writer{
		lg:       logger,
		model:    model,
		category: category,
	}
}

func (w *writer) Printf(ctx context.Context, severity models.Severity, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	w.lg.GetLogger().WithContext(ctx).Print(w.model, msg, w.category, severity)
}
