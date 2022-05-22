package poll_storage

type Storage interface {
	IsCurrentUserVoted(chatID int64, pollID string, userID int64) (bool, error)
	SetUserVoted(chatID int64, pollID string, userId int64) error
	CreateNewPoll(chatID int64, pollID string) error
}
