package redis

import (
	"github.com/ShamilKhal/shgo/internal/entity"

	"github.com/redis/go-redis/v9"
)

func DeserialiseContactList(contacts []redis.Z) []entity.ContactList {
	contactList := make([]entity.ContactList, 0, len(contacts))

	for _, contact := range contacts {
		contactList = append(contactList, entity.ContactList{
			Username:     contact.Member.(string),
			LastActivity: int64(contact.Score),
		})
	}

	return contactList
}
