package telegrambot

import (
	"github.com/Lexographics/logar/proxy"
	"github.com/Lexographics/logar/proxy/telegramlogger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot interface {
	ProxyTo(chatId int64) proxy.ProxyTarget
	Send(message string, chatId int64) error
}

type telegramBot struct {
	bot *tgbotapi.BotAPI
}

func (t *telegramBot) ProxyTo(chatId int64) proxy.ProxyTarget {
	return telegramlogger.New(t, chatId)
}

func (f *telegramBot) Send(message string, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := f.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func New(apiKey string) (TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}

	return &telegramBot{
		bot: bot,
	}, nil
}
