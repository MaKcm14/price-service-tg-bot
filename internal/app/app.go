package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	conf "github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/config"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/controller/tgbot"
)

// Service defines the configured application that is ready to be started.
type Service struct {
	logger *slog.Logger
	bot    tgbot.TgBot
}

func NewService() Service {
	date := strings.Split(time.Now().String()[:19], " ")

	mainLogFile, err := os.Create(fmt.Sprintf("../../logs/tg-bot-main-logs_%s___%s.txt",
		date[0], strings.Join(strings.Split(date[1], ":"), "-")))

	if err != nil {
		panic(fmt.Sprintf("error of creating the main-log-file: %v", err))
	}

	log := slog.New(slog.NewTextHandler(mainLogFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("main application's configuring begun")

	config := conf.NewSettings(log)

	bot, err := tgbot.New(config.TgBotToken, log)

	if err != nil {
		log.Error(fmt.Sprintf("error of starting the bot: %v", err))
		panic(err)
	}

	return Service{
		logger: log,
		bot:    bot,
	}
}

// Run starts the application and every connected service.
func (s Service) Run() {
	defer s.logger.Info("the bot was fully STOPPED")

	s.logger.Info("starting the bot and the other services")
	s.bot.Run()
}
