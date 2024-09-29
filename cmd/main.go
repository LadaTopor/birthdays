package main

import (
	"birthdays/cmd/config"
	"birthdays/notificatior"
	"birthdays/pkg/tgbot"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	conf := config.LoadConfig(logger)
	db := NewPostgresConnection(conf, logger)
	tgBot := tgbot.NewTgBot(conf, logger, db)

	notif := notificatior.NewNotificator(logger, tgBot, db)
	notif.Go()

	e := echo.New()

	e.Logger.Fatal(e.Start(":8080"))
}
