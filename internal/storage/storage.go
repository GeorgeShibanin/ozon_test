package storage

import (
	"context"
	"errors"
	"fmt"
)

var (
	StorageError    = errors.New("storage")
	ErrCollision    = fmt.Errorf("%w.collision", StorageError)
	ErrAlreadyExist = errors.New("key already exist")
	ErrNotFound     = fmt.Errorf("%w.not_found", StorageError)
)

type ShortedURL string
type URL string

type Storage interface {
	PutURL(ctx context.Context, url URL) (ShortedURL, error)
	GetURL(ctx context.Context, key ShortedURL) (URL, error)
}
