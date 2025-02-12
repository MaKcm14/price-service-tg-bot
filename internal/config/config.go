package conf

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Settings defines the main application's settings.
type Settings struct {
	TgBotToken string
	DSN        string
}

func NewSettings(log *slog.Logger) Settings {
	const errVar = "check this var is set or its value is not empty and try again"

	if err := godotenv.Load("../../.env"); err != nil {
		errRecord := fmt.Sprintf("error of loading the .env file: %v", err)
		log.Error(errRecord)
		panic(errRecord)
	}

	token := os.Getenv("BOT_TOKEN")

	if len(token) == 0 {
		log.Error(fmt.Sprintf("error of BOT_TOKEN var: %v", errVar))
		panic(fmt.Sprintf("error of BOT_TOKEN var: %v", errVar))
	}

	dsn := os.Getenv("DSN")

	if len(dsn) == 0 {
		log.Error(fmt.Sprintf("error of DSN var: %v", errVar))
		panic(fmt.Sprintf("error of DSN var: %v", errVar))
	}

	return Settings{
		TgBotToken: token,
		DSN:        dsn,
	}
}
