package queues

import (
	"context"
	"encoding/json"
	"notify/models"

	"github.com/go-redis/redis/v8"
)

func EnqueueNotification(ctx context.Context, notificationType string, notification models.Notification) error {
	// Enqueue the notification in Redis based on type
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	data, _ := json.Marshal(notification)
	return client.RPush(ctx, "queue:"+notificationType, data).Err()
}

// func DequeueNotification(notificationType string) (*models.Notification, error) {
// 	var notification models.Notification
// 	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
// 	data, err := client.LPop(ctx, "queue:"+notificationType).Result()
// 	if err == redis.Nil {
// 		// Queue is empty, so this is not an error
// 		return nil, nil
// 	} else if err != nil {
// 		// Return other errors (like connection issues)
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	// if err != nil {
// 	// 	return notification, err
// 	// }
// 	json.Unmarshal([]byte(data), &notification)
// 	return &notification, nil
// }

func DequeueNotification(ctx context.Context, redisClient *redis.Client, queueName string) (*models.Notification, error) {
	result, err := redisClient.BLPop(ctx, 0, queueName).Result()
	if err != nil {
		if err == redis.Nil {
			// Queue is empty
			return nil, nil
		}
		return nil, err
	}

	var notification models.Notification
	err = json.Unmarshal([]byte(result[1]), &notification)
	if err != nil {
		return nil, err
	}

	return &notification, nil
}
