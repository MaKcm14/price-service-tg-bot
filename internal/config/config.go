package conf

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Settings defines the main application's settings.
type Settings struct {
	TgBotToken string
}

func NewSettings() Settings {
	if err := godotenv.Load("../../.env"); err != nil {
		panic(fmt.Sprintf("error of loading the .env file: %v", err))
	}

	token := os.Getenv("BOT_TOKEN")

	if len(token) == 0 {
		panic("error of BOT_TOKEN var: check this var is set or its value is not empty and try again")
	}

	return Settings{
		TgBotToken: token,
	}
}
