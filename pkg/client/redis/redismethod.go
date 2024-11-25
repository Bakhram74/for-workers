package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ShamilKhal/shgo/internal/entity"
	"github.com/redis/go-redis/v9"
)

// FetchContactList of the user. It includes all the messages sent to and received by contact
// It will return a sorted list by last activity with a contact
func FetchContactList(id string) ([]entity.ContactList, error) {
	zRangeArg := redis.ZRangeArgs{
		Key:   contactListZKey(id), //"contacts:" + username
		Start: 0,
		Stop:  -1,
		Rev:   true,
	}
	// redis-cli
	// SYNTAX: ZRANGE key from_index to_index REV WITHSCORES
	// ZRANGE contacts:id 0 -1 REV WITHSCORES
	res, err := redisClient.client.ZRangeArgsWithScores(context.Background(), zRangeArg).Result()

	if err != nil {
		log.Println("error while fetching contact list. username: ",id, err)
		return nil, err
	}

	contactList := DeserialiseContactList(res)

	return contactList, nil
}

// UpdateContactList add contact to username's contact list
// if not present or update its timestamp as last contacted
func updateContactList(username, contact string) error {
	zs := redis.Z{Score: float64(time.Now().Unix()), Member: contact}

	// redis-cli SCORE is always float or int
	// SYNTAX: ZADD key SCORE MEMBER
	// ZADD contacts:username 1661360942123 contact
	err := redisClient.client.ZAdd(context.Background(),
		contactListZKey(username), //"contacts:" + username
		zs,
	).Err()
	//TODO error
	if err != nil {
		log.Println("error while updating contact list. username: ",
			username, "contact:", contact, err)
		return err
	}

	return nil
}

func CreateChat(chat entity.Chat) (string, error) {

	chatKey := chatKey() //fmt.Sprintf("chat#%d", time.Now().UnixMilli())

	by, _ := json.Marshal(chat)

	// redis-cli
	// SYNTAX: JSON.SET key $ json_in_string
	// JSON.SET chat#1661360942123 $ '{"from": "sun", "to":"earth","message":"good morning!"}'
	res, err := redisClient.client.Do(
		context.Background(),
		"JSON.SET",
		chatKey,
		"$",
		string(by),
	).Result()
	//TODO error
	if err != nil {
		log.Println("error while setting chat json", err)
		return "", err
	}

	log.Println("chat successfully set", res)

	// add contacts to both user's contact list
	err = updateContactList(chat.From, chat.To)
	if err != nil {
		log.Println("error while updating contact list of", chat.From) //TODO error
	}

	err = updateContactList(chat.To, chat.From)
	if err != nil {
		log.Println("error while updating contact list of", chat.To)
	}

	return chatKey, nil
}

func FetchChatBetween(user1, user2, fromTS, toTS string) ([]entity.Chat, error) {
	// Build the search query
	query := fmt.Sprintf(
		"@from:{%s|%s} @to:{%s|%s} @timestamp:[%s %s]",
		user1, user2, user1, user2, fromTS, toTS,
	)
	// 	// redis-cli
	// SYNTAX: FT.SEARCH index query
	// FT.SEARCH idx#chats '@from:{user2|user1} @to:{user1|user2} @timestamp:[0 +inf] SORTBY timestamp DESC'
	res, err := redisClient.client.Do(context.Background(),
		"FT.SEARCH",
		chatIndex(),
		query,
		"SORTBY", "timestamp", "DESC",
	).Result()

	if err != nil {
		return nil, err
	}
	// Check the type of the response
	responseMap, ok := res.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", res)
	}
	// Check total results
	totalResults, ok := responseMap["total_results"].(int64)
	if !ok || totalResults == 0 {
		return []entity.Chat{}, nil // No results found
	}
	// Extract results
	results, ok := responseMap["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected results type: %T", responseMap["results"])
	}
	var chats []entity.Chat
	// Iterate over each result
	for _, result := range results {
		resultMap, ok := result.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected result entry type: %T", result)
		}
		// Extract extra attributes from the result
		extraAttributes, ok := resultMap["extra_attributes"].(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected extra_attributes type: %T", resultMap["extra_attributes"])
		}
		// Extract chat data
		chatDataJSON, ok := extraAttributes["$"].(string) // Assuming chat data is stored as a JSON string
		if !ok {
			return nil, fmt.Errorf("unexpected chat data type: %T", extraAttributes["$"])
		}
		var chatData entity.Chat
		if err := json.Unmarshal([]byte(chatDataJSON), &chatData); err != nil {
			return nil, fmt.Errorf("error unmarshalling chat data: %w", err)
		}
		// Append to chats slice
		chats = append(chats, chatData)
	}

	return chats, nil
}

func CreateFetchChatBetweenIndex() {
	res, err := redisClient.client.Do(context.Background(),
		"FT.CREATE",
		chatIndex(), //"idx#chats"
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA", "$.from", "AS", "from", "TAG",
		"$.to", "AS", "to", "TAG",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE",
	).Result()

	fmt.Println(res, "CreateFetchChatBetweenIndex", err)
}

func SetPin(ctx context.Context, id string, value entity.PinValue) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = redisClient.client.Set(ctx, id, val, redisClient.config.Redis.ExpiredAt).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetPin(ctx context.Context, id string) (entity.PinValue, error) {
	var v entity.PinValue

	value := redisClient.client.Get(ctx, id)

	err := json.Unmarshal([]byte(value.Val()), &v)
	if err != nil {
		return entity.PinValue{}, errors.New("redis: key not found")
	}
	return v, nil
}

func DeletePinValue(ctx context.Context, id string) error {
	err := redisClient.client.Del(ctx, id).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetLimit(ctx context.Context, key string) error {

	var maxLimit int64 = 3
	var counter int64

	counter, err := redisClient.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			err := redisClient.client.Set(ctx, key, 1, redisClient.config.Redis.ExpiredAt).Err()
			if err != nil {
				return err
			} else {
				return nil
			}
		}
		return err
	}
	if counter >= maxLimit {
		return fmt.Errorf("SetLimit: limit has reached try after %v", redisClient.client.TTL(ctx, key))
	}
	err = redisClient.client.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
