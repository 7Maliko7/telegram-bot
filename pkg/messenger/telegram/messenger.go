package telegram

import tgAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Messenger struct {
	bot   *tgAPI.BotAPI
	UChan tgAPI.UpdatesChannel
}

func (m *Messenger) Send(msg tgAPI.MessageConfig) (tgAPI.Message, error) {
	msg.ParseMode = tgAPI.ModeHTML
	result, err := m.bot.Send(msg)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (m *Messenger) SendText(chatID int64, body string) (tgAPI.Message, error) {
	msg := tgAPI.NewMessage(chatID, body)

	result, err := m.Send(msg)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (m *Messenger) SendKeyboard(chatID int64, body string, kb Keyboard) (tgAPI.Message, error) {
	var keyboardRow [][]tgAPI.KeyboardButton
	for _, r := range kb {
		var row []tgAPI.KeyboardButton
		for _, b := range r {
			row = append(row, tgAPI.NewKeyboardButton(b))
		}

		keyboardRow = append(keyboardRow, tgAPI.NewKeyboardButtonRow(row...))
	}
	keyboard := tgAPI.NewReplyKeyboard(keyboardRow...)

	msg := tgAPI.NewMessage(chatID, body)
	msg.ReplyMarkup = keyboard

	result, err := m.Send(msg)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (m *Messenger) CloseKeyboard(chatID int64, body string) (tgAPI.Message, error) {
	msg := tgAPI.NewMessage(chatID, body)
	msg.ReplyMarkup = tgAPI.NewRemoveKeyboard(true)

	result, err := m.Send(msg)
	if err != nil {
		return result, err
	}

	return result, nil
}
