package config

import (
	"os"

	"github.com/Daniil-Sakharov/HockeyProject/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *Config

// Config глобальная конфигурация приложения
type Config struct {
	Logger   LoggerConfig
	Postgres PostgresConfig
	Telegram TelegramConfig
	FHSPB    FHSPBConfig
}

// Load загружает конфигурацию из .env файла и переменных окружения
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	telegramCfg, err := env.NewTelegramConfig()
	if err != nil {
		return err
	}

	fhspbCfg, err := env.NewFHSPBConfig()
	if err != nil {
		return err
	}

	appConfig = &Config{
		Logger:   loggerCfg,
		Postgres: postgresCfg,
		Telegram: telegramCfg,
		FHSPB:    fhspbCfg,
	}

	return nil
}

// AppConfig возвращает глобальную конфигурацию
func AppConfig() *Config {
	return appConfig
}
