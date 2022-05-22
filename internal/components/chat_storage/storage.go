package chat_storage

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Storage interface {
	GetChatUsers(chatID int64) ([]*tgbotapi.User, error)
	RegistrateNewUser(chatID int64, user *tgbotapi.User) error
}
