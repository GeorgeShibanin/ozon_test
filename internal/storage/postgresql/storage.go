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
	for attempt := 0; attempt < 5; attempt++ {
		key := storage2.URLKey(generator.GetRandomKey())
		err := s.links.
			QueryRow(context.Background(), "select id, url from links where url = $1", url).
			Scan(&key, &url)
		if err != nil {
			continue
			//return "l", fmt.Errorf("something went wrong - %w", storage2.StorageError)
		}
		return key, nil
	}
	return "", fmt.Errorf("too much attempts during inserting - %w", storage2.ErrCollision)
}

func (s *storage) GetURL(ctx context.Context, key storage2.URLKey) (storage2.ShortedURL, error) {

	return "link", nil
}
