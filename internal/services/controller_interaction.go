package services

type (
	UserActions interface {
		IdentifyUser(chatID int64) error
	}
)
