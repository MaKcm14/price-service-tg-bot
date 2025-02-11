package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	conf "github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/config"
)

// Service defines the configured application that is ready to be started.
type Service struct {
	logger *slog.Logger
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

	config := conf.NewSettings()

	_ = config

	return Service{
		logger: log,
	}
}

func (s Service) Run() {

}
