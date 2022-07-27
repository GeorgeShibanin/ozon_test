package in_memory

import (
	"context"
	"github.com/GeorgeShibanin/ozon_test/internal/generator"
	"sync"

	"github.com/GeorgeShibanin/ozon_test/internal/storage"
)

const RetriesCount = 5

type inMemoryStore struct {
	mutex sync.RWMutex
	store map[storage.ShortedURL]storage.URL
}

func Init() *inMemoryStore {
	return &inMemoryStore{
		mutex: sync.RWMutex{},
		store: make(map[storage.ShortedURL]storage.URL),
	}
}

func (s *inMemoryStore) PutURL(ctx context.Context, url storage.URL) (storage.ShortedURL, error) {
	var key storage.ShortedURL
	for k, v := range s.store {
		if url == v {
			key = k
			break
		}
	}
	if key == "" {
		//generate unique key
		for i := 0; i < RetriesCount; i++ {
			key = generator.GetRandomKey()
			if _, ok := s.store[key]; ok {
				return "", storage.ErrAlreadyExist
			}
		}
		s.mutex.Lock()
		s.store[key] = url
		s.mutex.Unlock()
	}
	return key, nil
}

func (s *inMemoryStore) GetURL(ctx context.Context, key storage.ShortedURL) (storage.URL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.store[key]
	if !ok {
		return "", storage.ErrNotFound
	}

	return value, nil
}
