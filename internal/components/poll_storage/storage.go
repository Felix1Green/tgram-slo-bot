package poll_storage

type Storage interface {
	IsCurrentUserVoted(chatID int64, pollID string, userID int64) (bool, error)
	SetUserVoted(chatID int64, pollID string, userId int64) error
	RemovePoll(pollKey string) error
	GetActivePollKeys() ([]string, error)
	GetPollInfoFromKey(key string) (*Poll, error)
	CreateNewPoll(chatID int64, pollID string) error
}
