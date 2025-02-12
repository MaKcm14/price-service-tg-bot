package services

import (
	"context"
	"fmt"
	"log/slog"
)

type UserInteractor struct {
	logger *slog.Logger

	repo Repository
}

func NewUserInteractor(log *slog.Logger, repo Repository) UserInteractor {
	return UserInteractor{
		logger: log,
		repo:   repo,
	}
}

// IdentifyUser checks whether user is in the DB and add it if he isn't in it.
func (u UserInteractor) IdentifyUser(chatID int64) error {
	const op = "services.identification"

	if flagExist, err := u.repo.IsUserExists(context.Background(), chatID); err != nil {
		u.logger.Warn(fmt.Sprintf("error of the %v: %v", op, err))
		return fmt.Errorf("%w: %w", ErrDBInteraction, err)
	} else if !flagExist {
		u.repo.AddUser(context.Background(), chatID)
	}

	return nil
}
