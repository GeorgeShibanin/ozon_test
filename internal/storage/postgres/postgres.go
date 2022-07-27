package postgres

import (
	"context"
	"fmt"
	"github.com/GeorgeShibanin/ozon_test/internal/generator"
	storage2 "github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	RetriesCount = 5

	GetQuery    = `SELECT key, url FROM links WHERE key = $1`
	InsertQuery = `INSERT INTO links (key, url) values ($1, $2)`

	dsnTemplate = "postgres://%s:%s@%s:%v/%s"
)

type storage struct {
	links *pgx.Conn
}

func Init(ctx context.Context, host, user, db, password string, port uint16) (*storage, error) {
	links, err := pgx.Connect(ctx, fmt.Sprintf(dsnTemplate, user, password, "database", port, db))
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}

	return &storage{links: links}, nil
}

func (s *storage) PutURL(ctx context.Context, url storage2.URL) (storage2.ShortedURL, error) {
	link := &Link{}
	for attempt := 0; attempt < RetriesCount; attempt++ {
		key := generator.GetRandomKey()
		err := s.links.QueryRow(ctx, GetQuery, key).
			Scan(&link.Key, &link.URL)

		if err != nil {
			return "", errors.Wrap(err, "can't get link")
		}

		if link.Key != "" {
			continue
		}

		err2, _ := s.links.Exec(ctx, InsertQuery, key, url)
		if err2 != nil {
			return "", fmt.Errorf("something went wrong - %w", storage2.StorageError)
		}
		return key, nil
	}
	return "", fmt.Errorf("too much attempts during inserting - %w", storage2.ErrCollision)
}

func (s *storage) GetURL(ctx context.Context, key storage2.ShortedURL) (storage2.URL, error) {
	link := &Link{}
	err := s.links.QueryRow(ctx, GetQuery, key).
		Scan(link.Key, &link.URL)
	if err != nil {
		return "", fmt.Errorf("something went wrong - %w", storage2.StorageError)
	}
	return storage2.URL(link.URL), err
}
