package main

import (
	"birthdays/cmd/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewPostgresConnection(config *config.Config, logger *logrus.Logger) *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Name)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		logger.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	return db
}
