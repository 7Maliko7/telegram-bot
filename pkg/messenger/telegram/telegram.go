package telegram

import (
	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	bot       *tgAPI.BotAPI
	messenger *Messenger
	handlers  map[string]func(tgAPI.Update) bool
	scenario  func(tgAPI.Update) bool
}

func New(token string) (*Telegram, error) {
	bot, err := tgAPI.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = false

	return &Telegram{
		bot:       bot,
		messenger: nil,
		handlers:  make(map[string]func(tgAPI.Update) bool),
	}, nil
}

func (m *Telegram) RegisterHandler(cmd string, h func(tgAPI.Update) bool) {
	m.handlers[cmd] = h
}

func (m *Telegram) RegisterScenario(h func(tgAPI.Update) bool) {
	m.scenario = h
}

func (m *Telegram) MakeMessenger() *Messenger {
	u := tgAPI.NewUpdate(0)
	u.Timeout = 60
	msgr := &Messenger{
		bot:   m.bot,
		UChan: m.bot.GetUpdatesChan(u),
	}

	m.messenger = msgr
	return msgr
}

func (m *Telegram) Listen() {
	for update := range m.messenger.UChan {
		if update.Message != nil {
			if handler, ok := m.handlers[update.Message.Command()]; ok {
				go handler(update)
			} else if handler, ok := m.handlers[update.Message.Text]; ok {
				go handler(update)
			} else {
				go m.scenario(update)
			}
		}
	}
}
