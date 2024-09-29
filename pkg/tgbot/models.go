package tgbot

import "birthdays/repo"

// Структура сообщения
type Message struct {
	Chat struct {
		ID int64 `json:"id"`
	} `json:"chat"`
	Text string `json:"text"`
}

// Структура для обработки обновлений от Telegram
type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

// Структура для обработки callback_query
type CallbackQuery struct {
	ID   string `json:"id"`
	Data string `json:"data"`
	From struct {
		ID int64 `json:"id"`
	} `json:"from"`
	Message struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		MessageID int `json:"message_id"`
	} `json:"message"`
	Birthday repo.Birthday `json:"birthday"`
}

// Структура запроса на отправку сообщения
type SendMessageRequest struct {
	ChatID      int64                `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup *ReplyKeyboardMarkup `json:"reply_markup,omitempty"`
	ParseMode   string               `json:"parse_mode,omitempty"`
}

// Структура для клавиатуры бота
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}
