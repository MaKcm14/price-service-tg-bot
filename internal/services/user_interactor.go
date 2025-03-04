package services

import (
	"context"
	"fmt"
	"log/slog"
)

// UserInteractor defines the logic of processing the some user's actions.
type UserInteractor struct {
	logger *slog.Logger
	repo   Repository
}

func NewUserInteractor(log *slog.Logger, repo Repository) UserInteractor {
	return UserInteractor{
		logger: log,
		repo:   repo,
	}
}

// IdentifyUser checks whether user is in the DB and add it if he isn't in it.
func (u UserInteractor) IdentifyUser(tgID int64) error {
	const op = "services.identification"

	if flagExist, err := u.repo.IsUserExists(context.Background(), tgID); err != nil {
		u.logger.Error(fmt.Sprintf("error of the %v: %v", op, err))
		return fmt.Errorf("error of the %s: %w: %w", op, ErrDBInteraction, err)
	} else if !flagExist {
		u.repo.AddUser(context.Background(), tgID)
	}

	return nil
}
