package in_memory

import (
	"context"
	"sync"

	"github.com/GeorgeShibanin/ozon_test/internal/storage"
)

type inMemoryStore struct {
	mutex sync.RWMutex
	store map[storage.ShortedURL]storage.URL
}

func Init() *inMemoryStore {
	return &inMemoryStore{
		//реализация потокобезопасности
		mutex: sync.RWMutex{},
		store: make(map[storage.ShortedURL]storage.URL),
	}
}

func (s *inMemoryStore) PutURL(ctx context.Context, key storage.ShortedURL, url storage.URL) (storage.ShortedURL, error) {
	var key1 storage.ShortedURL
	//проверяем наличие ссылки в хранилище
	for k, v := range s.store {
		if url == v {
			key1 = k
			break
		}
	}
	if key1 == "" {
		//закрываем хранилище для конкурентной записи
		s.mutex.Lock()
		//откроем обратно после завершения функции
		defer s.mutex.Unlock()

		_, ok := s.store[key]
		if ok {
			return "", storage.ErrAlreadyExist
		}

		s.store[key] = url
		return key, nil
	} else {
		return key1, nil
	}

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
