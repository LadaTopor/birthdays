package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Name     string `yaml:"POSTGRES_DB"`
	User     string `yaml:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	Port     string `yaml:"PORT"`
	Host     string `yaml:"HOST"`
	BotToken string `yaml:"botToken"`
	ChatId   int64  `yaml:"chatId"`
}

func LoadConfig(logger *logrus.Logger) *Config {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Fatalf("Ошибка чтения файла YAML: %v", err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		logger.Fatalf("Ошибка разбора YAML: %v", err)
	}

	return &config
}
