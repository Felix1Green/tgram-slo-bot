package poll_storage

type Storage interface {
	IsCurrentUserVoted(pollID string, userID int64) (bool, error)
	SetUserVoted(pollID string, userId int64) error
	RemovePoll(pollKey string) error
	GetActivePollKeys() ([]string, error)
	GetPollInfoFromKey(key string) (*Poll, error)
	CreateNewPoll(chatID int64, pollID string, msgID int) error
	GetActiveChatPoll(chatId int64) (string, error)
	SetActiveChatPoll(chatId int64, pollId string) error
	SetUserRegistered(pollID string, userId int64) error
	IsCurrentUserRegistered(pollID string, userID int64) (bool, error)
}
