package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"ozon_test/internal/generator"
	storage2 "ozon_test/internal/storage"
)

type storage struct {
	links *pgx.Conn
}

func NewStorage(db *pgx.Conn) *storage {
	return &storage{links: db}
}

//func ensureIndexes(ctx context.Context, collection *mongo.Collection) {
//	//indexModels := []mongo.IndexModel{
//	//	{
//	//		Keys: bsonx.Doc{{Key: "_id", Value: bsonx.Int32(1)}},
//	//	},
//	//}
//	//opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
//	//
//	//_, err := collection.Indexes().CreateMany(ctx, indexModels, opts)
//	//if err != nil {
//	//	panic(fmt.Errorf("failed to ensure indexes %w", err))
//	//}
//}

func (s *storage) PutURL(ctx context.Context, url storage2.ShortedURL) (storage2.URLKey, error) {
	var id string
	for attempt := 0; attempt < 5; attempt++ {
		key := generator.GetRandomKey()
		err := s.links.
			QueryRow(ctx, "insert into postgres (id, url) values ($1, $2) returning id", string(key), string(url)).Scan(&id)
		if err != nil {
			return "", fmt.Errorf("something went wrong - %w", storage2.StorageError)
		}
		return storage2.URLKey(id), nil
	}
	return "", fmt.Errorf("too much attempts during inserting - %w", storage2.ErrCollision)
}

func (s *storage) GetURL(ctx context.Context, key storage2.URLKey) (storage2.ShortedURL, error) {
	var url string
	err := s.links.
		QueryRow(ctx, "select id, url from postgres where id = $1", string(key)).
		Scan(&url)
	if err != nil {
		return "", fmt.Errorf("something went wrong - %w", storage2.StorageError)
	}
	return storage2.ShortedURL(url), err
}
