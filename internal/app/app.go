package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	conf "github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/config"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/controller/tgbot"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/api"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/kafka"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/postgres"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
)

// Service defines the configured application that is ready to be started.
type Service struct {
	logger      *slog.Logger
	bot         *tgbot.TgBot
	mainLogFile *os.File
	dbConn      postgres.PostgreSQLRepo
}

func NewService() *Service {
	log, file := mustSetLogger()

	log.Info("main application's configuring begun")

	config := conf.NewSettings(log,
		conf.ConfigBotToken("BOT_TOKEN"),
		conf.ConfigDSN("DSN"),
		conf.ConfigSocket("PRICE_SERVICE_SOCKET"),
		conf.ConfigBrokers("BROKERS"),
	)

	bot, dbConn := mustSetBot(log, config)

	return &Service{
		logger:      log,
		bot:         bot,
		dbConn:      dbConn,
		mainLogFile: file,
	}
}

// Run starts the application and every connected service.
func (s *Service) Run() {
	defer s.mainLogFile.Close()
	defer s.bot.Close()
	defer s.dbConn.Close()
	defer s.logger.Info("the bot was fully STOPPED")

	s.logger.Info("starting the bot and the other services")
	s.bot.Run()
}

func mustSetLogger() (*slog.Logger, *os.File) {
	date := strings.Split(time.Now().String()[:19], " ")

	mainLogFile, err := os.Create(fmt.Sprintf("../../logs/tg-bot-main-logs_%s___%s.txt",
		date[0], strings.Join(strings.Split(date[1], ":"), "-")))

	if err != nil {
		panic(fmt.Sprintf("error of creating the main-log-file: %v", err))
	}

	return slog.New(slog.NewTextHandler(mainLogFile, &slog.HandlerOptions{Level: slog.LevelInfo})), mainLogFile
}

func mustSetBot(log *slog.Logger, config conf.Settings) (*tgbot.TgBot, postgres.PostgreSQLRepo) {
	log.Info("connecting to the DB begun")

	dbConn, err := postgres.New(context.Background(), config.DSN, log)

	if err != nil {
		log.Error(fmt.Sprintf("error starting the DB: %v", err))
		panic(err)
	}

	products := make(chan *tgbot.TrackedProduct)

	reader, err := kafka.NewConsumer(config.Brokers, log, products)

	if err != nil {
		log.Error(fmt.Sprintf("error of Kafka: %v", err))
		panic(err)
	}

	bot, err := tgbot.New(config.TgBotToken, log,
		services.NewUserInteractor(log, dbConn),
		api.NewPriceServiceApi(config.PriceServiceSocket, log),
		dbConn, products, reader)

	if err != nil {
		log.Error(fmt.Sprintf("error of configuring the bot: %v", err))
		panic(err)
	}

	return bot, dbConn
}
