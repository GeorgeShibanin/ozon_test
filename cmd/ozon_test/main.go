package main

import (
	"context"
	"github.com/GeorgeShibanin/ozon_test/internal/ratelimit"
	"github.com/GeorgeShibanin/ozon_test/internal/storage/rediscachedstorage"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GeorgeShibanin/ozon_test/internal/config"
	"github.com/GeorgeShibanin/ozon_test/internal/handlers"
	"github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/GeorgeShibanin/ozon_test/internal/storage/in_memory"
	"github.com/GeorgeShibanin/ozon_test/internal/storage/postgres"
)

type ConnectionType string

const (
	ConnectionTypePostgres ConnectionType = "postgres"
	ConnectionTypeRedis    ConnectionType = "redis"
	ConnectionTypeInMemory ConnectionType = "in_memory"
)

func main() {
	srv := NewServer()
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func NewServer() *http.Server {
	r := mux.NewRouter()

	var store storage.Storage
	var err error

	storageMode := ConnectionType(os.Getenv("STORAGE_MODE"))

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Redis_URL,
	})
	rateLimitFactory := ratelimit.NewFactory(redisClient)

	switch storageMode {
	case ConnectionTypePostgres:
		store = initPostgres()
	case ConnectionTypeRedis:
		store = initPostgres()
		store, err = rediscachedstorage.Init(redisClient, store)
		if err != nil {
			log.Fatalf("can't init postgres connection: %s", err.Error())
		}
	case ConnectionTypeInMemory:
		store = in_memory.Init()
	default:
		store = initPostgres()
	}

	handler := handlers.NewHTTPHandler(store, rateLimitFactory)
	r.HandleFunc("/{shortUrl:\\w{10}}", handler.HandleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/urls", handler.HandlePostUrl)

	return &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func initPostgres() *postgres.StoragePostgres {
	store, err := postgres.Init(
		context.Background(),
		config.PostgresHost,
		config.PostgresUser,
		config.PostgresDB,
		config.PostgresPassword,
		config.PostgresPort,
	)
	if err != nil {
		log.Fatalf("can't init postgres connection: %s", err.Error())
	}
	return store
}
