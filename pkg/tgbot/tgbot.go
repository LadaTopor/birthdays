package tgbot

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"birthdays/cmd/config"
	"birthdays/repo"
)

type TgBot struct {
	log   *logrus.Logger
	repo  *repo.Repo
	token string
	//chatId int64
}

func NewTgBot(config *config.Config, logger *logrus.Logger, db *sql.DB) *TgBot {
	return &TgBot{
		log:   logger,
		repo:  repo.NewRepo(db),
		token: config.BotToken,
		//chatId: config.ChatId,
	}
}

var userStates = make(map[int64]string)

const telegramAPI = "https://api.telegram.org/bot"

func (tg *TgBot) Start() {
	offset := 0

	// Бесконечный цикл получения обновлений
	for {
		updates, err := tg.getUpdates(offset)
		if err != nil {
			tg.log.Printf("Error getting updates: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Обработка каждого обновления
		for _, update := range updates {
			offset = update.UpdateID + 1
			// Если это callback_query (нажатие на кнопку)
			//if update.CallbackQuery != nil {
			//	tg.log.Printf("Received callback query: %s", update.CallbackQuery.Data)
			//	tg.handleCallbackQuery(*update.CallbackQuery)
			//}

			// Если это сообщение от пользователя и оно содержит команду /start
			if update.Message != nil {
				if update.Message.Text == "/start" {

					chatID := update.Message.Chat.ID
					tg.log.Printf("Received /start from chat ID %d", chatID)
					// Отправляем кнопки
					err = tg.sendButtonsMessage(chatID)
					if err != nil {
						tg.log.Printf("Failed to send buttons: %v", err)
					}
				} else {
					tg.handleMessage(update.Message)
				}

			}
		}

		time.Sleep(1 * time.Second)
	}
}

//// Функция отправки POST-запроса на сервер
//func (tg *TgBot) sendRequestToServer(url string) error {
//	reqBody := []byte(`{"status": "request from bot"}`)
//
//	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
//	if err != nil {
//		return fmt.Errorf("failed to send request to server: %v", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("server returned non-200 status: %d", resp.StatusCode)
//	}
//
//	body, _ := ioutil.ReadAll(resp.Body)
//	tg.log.Printf("Response from server: %s", string(body))
//
//	return nil
//}

// Функция отправки сообщения с кнопками
func (tg *TgBot) sendButtonsMessage(chatID int64) error {
	path := fmt.Sprintf("%s%s/SendMessage", telegramAPI, tg.token)

	// Создаем клавиатуру
	keyboard := ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{
				{Text: "Список"},
				{Text: "Добавить"},
				{Text: "Удалить"},
			},
		},
		ResizeKeyboard:  true,  // Делает клавиатуру адаптивной под размер экрана
		OneTimeKeyboard: false, // Клавиатура останется на экране, пока не будет заменена
	}

	// Создаем запрос на отправку сообщения с кнопками
	reqBody, err := json.Marshal(SendMessageRequest{
		ChatID:      chatID,
		Text:        "Выбери действие:",
		ReplyMarkup: &keyboard,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Отправляем запрос
	resp, err := http.Post(path, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status: %s", resp.Status)
	}

	return nil
}

//
//// DELETE Функция обработки callback_query (нажатий на кнопки)
//func (tg *TgBot) handleCallbackQuery(query CallbackQuery) {
//	//var serverURL string
//
//	// Определяем, на какую кнопку нажал пользователь
//	if query.Data == "get_birthdays" {
//		list, err := tg.repo.GetListBirthdays()
//		if err != nil {
//			tg.log.Error(err)
//			return
//		}
//		tg.SendMessage(formatBirthdays(list))
//	} else if query.Data == "add_birthday" {
//		// Отправляем сообщение с инструкцией по вводу имени и даты рождения
//		tg.SendMessage("Введите Имя и дату рождения в формате дд.мм.гггг")
//
//		// Устанавливаем состояние для пользователя, что бот ожидает имя и дату
//		userStates[query.From.ID] = "waiting_for_name_and_birthdate"
//		return
//	} else if query.Data == "del_birthday" {
//		//serverURL = "http://localhost:1323/del/birthday"
//	}
//
//	// Отправляем запрос на сервер
//	//err := tg.sendRequestToServer(serverURL)
//	//if err != nil {
//	//	tg.log.Printf("Error sending request to server: %v", err)
//	//	// Отправляем сообщение об ошибке пользователю
//	//	tg.SendMessage("Ошибка при отправке запроса.")
//	//} else {
//	//	// Отправляем успешное сообщение пользователю
//	//	tg.SendMessage("Запрос отправлен успешно.")
//	//}
//
//	// Ответ на callback_query (Telegram требует ответ)
//	tg.answerCallbackQuery(query.ID)
//}

// DELETE Функция для ответа на callback_query
//func (tg *TgBot) answerCallbackQuery(queryID string) {
//	path := fmt.Sprintf("%s%s/answerCallbackQuery", telegramAPI, tg.token)
//
//	// Отправляем пустой ответ на callback_query
//	reqBody := fmt.Sprintf(`{"callback_query_id": "%s"}`, queryID)
//
//	resp, err := http.Post(path, "application/json", bytes.NewBuffer([]byte(reqBody)))
//	if err != nil {
//		tg.log.Printf("Failed to answer callback query: %v", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		tg.log.Printf("Failed to answer callback query, status: %s", resp.Status)
//	}
//}

func (tg *TgBot) SendMessage(chatId int64, message string) {
	path := fmt.Sprintf("%s%s/SendMessage", telegramAPI, tg.token)

	tg.log.Println(message)
	reqBody, err := json.Marshal(SendMessageRequest{
		ChatID:    chatId,
		Text:      message,
		ParseMode: "Markdown",
	})
	if err != nil {
		tg.log.Errorf("failed to marshal request: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	//r, _ := httputil.DumpRequest(req, true)
	//fmt.Println(string(r))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		tg.log.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r, _ := io.ReadAll(resp.Body)
		tg.log.Errorf("failed to send message, status: %s", r)
	}

	return
}

func (tg *TgBot) handleMessage(message *Message) {
	chatID := message.Chat.ID
	userState := userStates[chatID]

	switch strings.ToLower(message.Text) {
	case "отмена":
		userStates[chatID] = ""

		tg.SendMessage(chatID, "Действие отменено.")
		return
	case "список":
		list, err := tg.repo.GetListBirthdays(chatID)
		if err != nil {
			tg.log.Error(err)
			return
		}
		tg.SendMessage(chatID, formatBirthdays(list))
		return
	case "добавить":
		tg.SendMessage(chatID, "Введите имя, фамилию и дату рождения в формате дд.мм.гггг\n(или \"Отмена\")")

		// Устанавливаем состояние для пользователя, что бот ожидает имя и дату
		userStates[chatID] = "waiting_for_name_and_birthdate"
		return
	case "удалить":
		tg.SendMessage(chatID, "Введите номер, который хотите удалить")

		// Устанавливаем состояние для пользователя, что бот ожидает имя и дату
		userStates[chatID] = "waiting_for_del_birthdate"
		return
	}

	// Если бот ожидает имя и дату рождения
	switch userState {
	case "waiting_for_name_and_birthdate":
		// Проверяем формат даты рождения (дд.мм.гггг)
		parts := strings.Split(message.Text, " ")
		if len(parts) != 3 {
			tg.SendMessage(chatID, "Неправильный формат. Должно быть: \"Фамилия Имя дд.мм.гггг\"")
			return
		}

		name := parts[0]
		surname := parts[1]
		birthdate := parts[2]

		fullname := fmt.Sprintf("%s %s", name, surname)

		// Пример валидации даты (более сложную валидацию можно добавить)
		if len(birthdate) < 8 {
			tg.SendMessage(chatID, "Неправильный формат даты. Должно быть: \"Фамилия Имя дд.мм.гггг\"")
			return
		} else if len(birthdate) == 9 {
			birthdate = "0" + birthdate
		}

		// Сохраняем данные пользователя (здесь можно отправить их на сервер или в базу данных)
		tg.log.Printf("Пользователь ввёл: Имя - %s, Дата рождения - %s", fullname, birthdate)
		err := tg.repo.AddNewBirthday(fullname, birthdate, chatID)
		if err != nil {
			tg.log.Error(err)
			return
		}

		// Отправляем сообщение с подтверждением
		tg.SendMessage(chatID, fmt.Sprintf("Сохранено: Имя - %s, Дата рождения - %s", fullname, birthdate))

		// Сбрасываем состояние пользователя
		userStates[chatID] = ""
		return
	case "waiting_for_del_birthdate":
		id, err := strconv.Atoi(message.Text)
		if err != nil {
			tg.SendMessage(chatID, "Введен неверный Id!")
			return
		}
		err = tg.repo.DeleteBirthday(id, chatID)
		if err != nil {
			tg.log.Error(err)
			tg.SendMessage(chatID, "Такого id не существует!")
			return
		}
		tg.SendMessage(chatID, fmt.Sprintf("№%d удален", id))
		userStates[chatID] = ""
		return
	}

	// Если состояние не определено, отправляем стандартное сообщение
	tg.SendMessage(chatID, "Используйте /start для начала.")
}

func (tg *TgBot) getUpdates(offset int) ([]Update, error) {
	path := fmt.Sprintf("%s%s/getUpdates?offset=%d", telegramAPI, tg.token, offset)

	resp, err := http.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !result.OK {
		return nil, fmt.Errorf("telegram API error")
	}

	return result.Result, nil
}

func formatBirthdays(birthdays []repo.Birthday) string {
	result := "Сохранные дни рождения:\n"
	for _, item := range birthdays {
		result += fmt.Sprintf("№%d %s %s\n", item.Id, item.FullName, item.Birthdate)
	}

	return result
}
