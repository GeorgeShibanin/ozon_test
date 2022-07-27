package postgres

import (
	"context"
	"fmt"
	storage2 "github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	GetIdQuery    = `SELECT id, url FROM links WHERE id = $1`
	GetByUrlQuery = `SELECT id, url FROM links WHERE url = $1`
	InsertQuery   = `INSERT INTO links (id, url) values ($1, $2)`

	dsnTemplate = "postgres://%s:%s@%s:%v/%s"
)

type Storage struct {
	Conn postgresInterface
}

type postgresInterface interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func Init(ctx context.Context, host, user, db, password string, port uint16) (*Storage, error) {
	conn, err := pgx.Connect(ctx, fmt.Sprintf(dsnTemplate, user, password, host, port, db))
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}

	return &Storage{Conn: conn}, nil
}

func (s *Storage) PutURL(ctx context.Context, key storage2.ShortedURL, url storage2.URL) (storage2.ShortedURL, error) {
	tx, err := s.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	link := &Link{}
	err = tx.QueryRow(ctx, GetIdQuery, key).Scan(&link.Key, &link.URL)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", errors.Wrap(err, "can't get by link")
	}

	if link.Key != "" {
		return "", storage2.ErrAlreadyExist
	}

	err = tx.QueryRow(ctx, GetByUrlQuery, url).Scan(&link.Key, &link.URL)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", errors.Wrap(err, "can't by url")
	}

	if link.URL != "" {
		return storage2.ShortedURL(link.Key), nil
	}

	tag, err := tx.Exec(ctx, InsertQuery, key, url)
	if err != nil {
		return "", errors.Wrap(err, "can't insert link")
	}

	if tag.RowsAffected() != 1 {
		return "", errors.Wrap(err, fmt.Sprintf("unexpected rows affected value: %v", tag.RowsAffected()))
	}

	return key, nil
}

func (s *Storage) GetURL(ctx context.Context, key storage2.ShortedURL) (storage2.URL, error) {
	link := &Link{}
	err := s.Conn.QueryRow(ctx, GetIdQuery, key).
		Scan(&link.Key, &link.URL)
	if err != nil {
		return "", fmt.Errorf("something went wrong - %w", storage2.StorageError)
	}
	return storage2.URL(link.URL), err
}
