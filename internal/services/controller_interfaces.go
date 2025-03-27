package services

import "context"

type (
	// UserConfiger defines the user's configs actions.
	UserConfiger interface {
		IdentifyUser(chatID int64) error
	}

	// Actor defines the interactions with the different actors.
	Actor interface {
		SendTrackedProducts(ctx context.Context)
	}
)
