package service

import (
	"strings"

	"github.com/ShamilKhal/shgo/internal/entity"
	"github.com/ShamilKhal/shgo/pkg/client/redis"
	"github.com/ShamilKhal/shgo/pkg/utils"
)

type ChatService struct {
	redis *redis.Redis
}

func NewChatService(redis *redis.Redis) *ChatService {
	return &ChatService{
		redis: redis,
	}
}

func (s *ChatService) GetContactList(id string) ([]entity.ContactList, error) {

	id = strings.ReplaceAll(id, "-", "")

	contactList, err := redis.FetchContactList(id)
	if err != nil {
		return []entity.ContactList{}, err
	}

	for i, v := range contactList {
		contactList[i] = entity.ContactList{
			Username:     utils.RecoverUUID(v.Username),
			LastActivity: v.LastActivity,
		}
	}

	return contactList, nil
}

func (s *ChatService) GetChatHistory(user1, user2, fromTS, toTS string) ([]entity.Chat, error) {

	user1 = strings.ReplaceAll(user1, "-", "")
	user2 = strings.ReplaceAll(user2, "-", "")

	chats, err := redis.FetchChatBetween(user1, user2, fromTS, toTS)
	if err != nil {
		return []entity.Chat{}, err
	}

	for i, v := range chats {
		chats[i] = entity.Chat{
			ID:        v.ID,
			From:      utils.RecoverUUID(v.From),
			To:        utils.RecoverUUID(v.To),
			Msg:       v.Msg,
			Timestamp: v.Timestamp,
		}
	}

	return chats, nil
}

func (s *ChatService) CreateChat(c entity.Chat) (string, error) {
	modified := entity.Chat{
		ID:        c.ID,
		From:      strings.ReplaceAll(c.From, "-", ""),
		To:        strings.ReplaceAll(c.To, "-", ""),
		Msg:       c.Msg,
		Timestamp: c.Timestamp,
	}
	return redis.CreateChat(modified)
}
