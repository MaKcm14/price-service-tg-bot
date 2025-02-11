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
}

func NewSettings(log *slog.Logger) Settings {
	if err := godotenv.Load("../../.env"); err != nil {
		errRecord := fmt.Sprintf("error of loading the .env file: %v", err)
		log.Error(errRecord)
		panic(errRecord)
	}

	token := os.Getenv("BOT_TOKEN")

	if len(token) == 0 {
		errRecord := "error of BOT_TOKEN var: check this var is set or its value is not empty and try again"
		log.Error(errRecord)
		panic(errRecord)
	}

	return Settings{
		TgBotToken: token,
	}
}
