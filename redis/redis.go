package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConn struct {
	client    *redis.Client
	topicName string
}

func New() (*RedisConn, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	return &RedisConn{
		client:    client,
		topicName: "messages",
	}, nil
}

func (r *RedisConn) SubscribeAndWriteToChann(ch chan []byte) {
	su := r.client.Subscribe(context.TODO(), r.topicName)

	_, err := su.Receive(context.Background())
	if err != nil {
		fmt.Println("Error occured in subscribe", err)
		return
	}

	redisChan := su.Channel()

	for {
		data, ok := <-redisChan
		if !ok {
			fmt.Println("Error occured while listening for messages", err)
			return
		}

		ch <- []byte(data.Payload)
	}
}

func (r *RedisConn) Publish(data []byte) error {
	return r.client.Publish(context.Background(), r.topicName, string(data)).Err()
}
