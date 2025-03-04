package services

// UserConfiger defines the user's configs actions.
type UserConfiger interface {
	IdentifyUser(chatID int64) error
}
