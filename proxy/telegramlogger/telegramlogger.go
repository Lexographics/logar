package telegramlogger

import (
	"fmt"
	"strings"

	"github.com/Lexographics/logar/models"
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
	return l.bot.Send(fmt.Sprintf("["+strings.ToUpper(log.Severity.String())+"] %s", log.Message), l.chatId)
}
