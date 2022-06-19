package games_storage

type Storage interface {
	GetChatGames(chatID int64) ([]string, error)
	CreateNewChatGame(chatID int64, game string) error
}
