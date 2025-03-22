package consolelogger

import (
	"fmt"
	"strings"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
	"github.com/Lexographics/logar/proxy"
)

type consoleLogger struct {
}

func New() proxy.ProxyTarget {
	return &consoleLogger{}
}

func (l *consoleLogger) Send(log models.Log, rawMesage string) error {
	fmt.Print("["+strings.ToUpper(log.Severity.String())+"] ", log.CreatedAt.Format(time.DateTime), " ", rawMesage, "\n")
	return nil
}
