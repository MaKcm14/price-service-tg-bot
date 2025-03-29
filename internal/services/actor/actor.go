package actor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
)

// APIActor defines the logic of api-interactions.
type APIActor struct {
	logger *slog.Logger
	api    services.ApiInteractor
	repo   services.Repository
}

func NewAPI(log *slog.Logger, repo services.Repository, api services.ApiInteractor) APIActor {
	return APIActor{
		api:    api,
		logger: log,
		repo:   repo,
	}
}

// SendTrackedProducts defines the logic of sending the tracked-products requests.
func (a APIActor) SendTrackedProducts(ctx context.Context) {
	const op = "services.send-tracked-products"

	time.Sleep(time.Second * 90)

	for {
		timeStart := time.Now()
		res, err := a.repo.GetUsersTrackedProducts(ctx)

		if err != nil {
			a.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
			time.Sleep(time.Hour * 24)
			continue
		}

		for chatID, request := range res {
			err = a.api.SendAsyncBestPriceRequest(request, map[string]string{
				"ChatID": fmt.Sprint(chatID),
			})

			if err != nil {
				a.logger.Warn(fmt.Sprintf("error of the %s: %s", op, err))
			}

			time.Sleep(time.Minute * 1)
		}

		for time.Since(timeStart) < time.Hour*24 {
			continue
		}
	}
}
