package rediscachedstorage

import (
	"context"
	"fmt"
	"github.com/GeorgeShibanin/ozon_test/internal/storage"
	storage2 "github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Storage struct {
	conn   storage.Storage
	client *redis.Client
}

func Init(redisClient *redis.Client, persistentStorage storage.Storage) (*Storage, error) {
	return &Storage{
		conn:   persistentStorage,
		client: redisClient,
	}, nil
}

func (s *Storage) PutURL(ctx context.Context, key storage2.ShortedURL, url storage2.URL) (storage2.ShortedURL, error) {
	urlPut, err := s.conn.PutURL(ctx, key, url)
	if err != nil {
		return urlPut, err
	}
	err = s.client.Set(ctx, "surl:"+string(key), string(urlPut), time.Hour).Err()
	if err != nil {
		log.Printf("Failed to insert key %s into cache due to an error: %s\n", key, err)
	}
	return urlPut, nil
}

func (s *Storage) GetURL(ctx context.Context, key storage2.ShortedURL) (storage2.URL, error) {
	get := s.client.Get(ctx, string(key))
	switch url, err := get.Result(); {
	case err == redis.Nil:
		// continue execution
	case err != nil:
		return "", fmt.Errorf("%w: failed to get value from redis due to error %s", storage.StorageError, err)
	default:
		log.Printf("Successfully obtained url from cache for key %s", key)
		return storage2.URL(url), nil
	}
	log.Printf("Loading post by key %s from persistent storage", key)
	urlGet, err := s.conn.GetURL(ctx, key)
	if err != nil {
		return urlGet, err
	}
	err = s.client.Set(ctx, "surl:"+string(key), string(urlGet), time.Hour).Err()
	if err != nil {
		log.Printf("Failed to insert key %s into cache due to an error: %s\n", key, err)
	}
	return urlGet, nil
}
