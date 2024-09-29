package notificatior

import (
	"birthdays/pkg/tgbot"
	"birthdays/repo"
	"database/sql"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"time"
)

type Notificator struct {
	db    *repo.Repo
	cron  *cron.Cron
	tgBot *tgbot.TgBot
	log   *logrus.Logger
}

func NewNotificator(logger *logrus.Logger, tgBot *tgbot.TgBot, db *sql.DB) *Notificator {
	return &Notificator{
		log:   logger,
		cron:  cron.New(),
		tgBot: tgBot,
		db:    repo.NewRepo(db),
	}
}

func (nf *Notificator) Go() {
	go nf.tgBot.Start()

	// Добавляем задачу, которая будет выполняться каждый день в 9 утра
	_, err := nf.cron.AddFunc("0 17 * * *", func() {
		birthdays, err := nf.db.GetAllOfListBirthdays()
		if err != nil {

			nf.log.Printf("Ошибка при получении пользователей: %v", err)
			return
		}
		nf.checkBirthdays(birthdays)
	})

	if err != nil {
		nf.log.Fatalf("Ошибка при добавлении cron-задачи: %v", err)
	}

	// Запуск cron планировщика
	nf.cron.Start()

}

func (nf *Notificator) checkBirthdays(birthdays []repo.Birthday) {
	now := time.Now().Truncate(24 * time.Hour)
	deltaDate := now.AddDate(0, 0, 6)
	location, _ := time.LoadLocation("Europe/Moscow")

	for _, item := range birthdays {
		date, err := time.Parse("02.01.2006", fmt.Sprintf("%s.%d", item.Birthdate[:len(item.Birthdate)-5], time.Now().Year()))
		if err != nil {
			nf.log.Error(err)
		}

		date = time.Date(date.Year(), date.Month(), date.Day(), 3, 0, 0, 0, location)
		if time.Now().Truncate(24 * time.Hour).Equal(date) {
			nf.tgBot.SendMessage(item.ChatId, fmt.Sprintf("У *%s* cегодня день рождения! (*%s*)\nНе забудь поздравить.", item.FullName, item.Birthdate))
		} else if date.After(now) && date.Before(deltaDate) {
			nf.tgBot.SendMessage(item.ChatId, fmt.Sprintf("У *%s* скоро день рождения! *%s*", item.FullName, item.Birthdate))
		}
	}
}
