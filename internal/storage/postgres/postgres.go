package postgres

import (
	"context"
	"fmt"

	"github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	GetIdQuery    = `SELECT id, url FROM links WHERE id = $1`
	GetByUrlQuery = `SELECT id, url FROM links WHERE url = $1`
	InsertQuery   = `INSERT INTO links (id, url) values ($1, $2)`

	dsnTemplate = "postgres://%s:%s@%s:%v/%s"
)

type StoragePostgres struct {
	conn postgresInterface
}

type postgresInterface interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func initConnection(conn postgresInterface) *StoragePostgres {
	return &StoragePostgres{conn: conn}
}

func Init(ctx context.Context, host, user, db, password string, port uint16) (*StoragePostgres, error) {
	//подключение к базе через переменные окружения
	conn, err := pgx.Connect(ctx, fmt.Sprintf(dsnTemplate, user, password, host, port, db))
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}

	return initConnection(conn), nil
}

func (s *StoragePostgres) PutURL(ctx context.Context, key storage.ShortedURL, url storage.URL) (storage.ShortedURL, error) {
	//объявляем транзакцию
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
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
	//далее в транзакции посылаем раличные запросы, после чего коммитим либо откатываемся если была ошибка
	link := &Link{}
	//проверяем если ли поступивший url уже в базе
	err = tx.QueryRow(ctx, GetByUrlQuery, url).Scan(&link.Key, &link.URL)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", errors.Wrap(err, "can't get by url")
	}
	if link.URL != "" {
		return storage.ShortedURL(link.Key), nil
	}

	//проверяем если ли сгенерировнный ключ уже в базе
	err = tx.QueryRow(ctx, GetIdQuery, key).Scan(&link.Key, &link.URL)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", errors.Wrap(err, "can't get by link")
	}

	if link.Key != "" {
		return "", storage.ErrAlreadyExist
	}

	//вставляем в базу новое значение
	tag, err := tx.Exec(ctx, InsertQuery, key, url)
	if err != nil {
		return "", errors.Wrap(err, "can't insert link")
	}

	if tag.RowsAffected() != 1 {
		return "", errors.Wrap(err, fmt.Sprintf("unexpected rows affected value: %v", tag.RowsAffected()))
	}

	return key, nil
}

func (s *StoragePostgres) GetURL(ctx context.Context, key storage.ShortedURL) (storage.URL, error) {
	link := &Link{}
	//получаем из базы значение по ключу
	err := s.conn.QueryRow(ctx, GetIdQuery, key).
		Scan(&link.Key, &link.URL)
	if err != nil {
		return "", fmt.Errorf("something went wrong - %w", storage.StorageError)
	}
	return storage.URL(link.URL), err
}
