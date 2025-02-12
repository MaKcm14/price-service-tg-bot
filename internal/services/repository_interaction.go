package services

import "context"

type (
	// RepoValidator defines the check-actions.
	RepoValidator interface {
		IsUserExists(ctx context.Context, chatID int64) (bool, error)
	}

	// RepoUpdater defines the modify-actions.
	RepoUpdater interface {
		AddUser(ctx context.Context, chatID int64) error
	}

	Repository interface {
		RepoValidator
		RepoUpdater
	}
)
