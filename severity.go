package logar

import "github.com/Lexographics/logar/internal/domain/models"

func Severity(severity int) models.Severity {
	return models.Severity(severity)
}

const (
	Trace   models.Severity = models.Severity_Trace
	Log     models.Severity = models.Severity_Log
	Info    models.Severity = models.Severity_Info
	Warning models.Severity = models.Severity_Warning
	Error   models.Severity = models.Severity_Error
	Fatal   models.Severity = models.Severity_Fatal
)
