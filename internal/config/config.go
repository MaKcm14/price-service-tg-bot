package config

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
)

// Settings defines the main application's settings.
type Settings struct {
	TgBotToken         string
	DSN                string
	PriceServiceSocket string
	Brokers            []string
}

func NewSettings(log *slog.Logger, opts ...ConfigOpt) Settings {
	var set Settings

	if err := godotenv.Load("../../.env"); err != nil {
		errRecord := fmt.Sprintf("error of loading the .env file: %v", err)
		log.Error(errRecord)
		panic(errRecord)
	}

	for _, opt := range opts {
		err := opt(&set)

		if err != nil {
			log.Error(err.Error())
			panic(err)
		}
	}

	return set
}
