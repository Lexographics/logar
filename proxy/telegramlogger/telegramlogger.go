package telegramlogger

import (
	"strings"
	"time"

	"github.com/Lexographics/logar/internal/domain/models"
)

type MessageSender interface {
	Send(message string, chatId int64) error
}

func New(bot MessageSender, chatId int64) *telegramLogger {
	return &telegramLogger{
		bot:    bot,
		chatId: chatId,
	}
}

type telegramLogger struct {
	bot    MessageSender
	chatId int64
}

func (l *telegramLogger) Send(log models.Log, rawMessage string) error {
	return l.bot.Send("["+strings.ToUpper(log.Severity.String())+"] "+log.CreatedAt.Format(time.RFC3339)+" "+log.Message, l.chatId)
}
