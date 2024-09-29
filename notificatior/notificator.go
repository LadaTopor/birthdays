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
	_, err := nf.cron.AddFunc("0 8 * * *", func() {
		birthdays, err := nf.db.GetListBirthdays()
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
	deltaDate := time.Now().AddDate(0, 0, 5)
	fmt.Println(deltaDate)

	for _, item := range birthdays {
		date, err := time.Parse("02.01.2006", fmt.Sprintf("%s.%d", item.Birthdate, time.Now().Year()))
		if err != nil {
			nf.log.Error(err)
		}
		fmt.Println(date)

		if deltaDate.After(date) || date.Equal(deltaDate) {
			nf.tgBot.SendMessage(fmt.Sprintf("У *%s* скоро день рождения! *%d %s*", item.FullName, date.Day(), date.Month().String()))
		}
	}
}
