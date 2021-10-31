package services

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(address, password string, db int) (*RedisService, error) {
	client, err := newRedisClient(address, password, db)
	if err != nil {
		return nil, err
	}

	return &RedisService{client: client}, nil
}

func (service RedisService) Set(key string, data map[string]interface{}) error {
	err := service.client.HSet(context.Background(), key, data).Err()
	if err != nil {
		return err
	}
	return nil
}

func (service RedisService) Get(key string) (map[string]string, error) {
	result, err := service.client.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func newRedisClient(address, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: address, Password: password, DB: db})
	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}
	return client, nil
}
