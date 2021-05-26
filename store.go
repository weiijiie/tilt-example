package main

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Store struct {
	client *redis.Client
}

func NewStore(addr, password string) Store {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return Store{
		client: client,
	}
}

func (s Store) Ping(c context.Context) (string, error) {
	return s.client.Ping(c).Result()
}

func (s Store) Get(c context.Context, key string) (string, error) {
	return s.client.Get(c, key).Result()
}

func (s Store) Set(c context.Context, key, value string) (string, error) {
	return s.client.Set(c, key, value, 0).Result()
}
