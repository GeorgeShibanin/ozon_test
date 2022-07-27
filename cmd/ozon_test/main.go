package main

import (
	"context"
	"fmt"
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

func NewServer() *http.Server {
	r := mux.NewRouter()

	var store storage.Storage
	var err error
	storageMode := ConnectionType(os.Getenv("STORAGE_MODE"))
	fmt.Println(storageMode)
	handler := &handlers.HTTPHandler{}
	switch storageMode {
	case ConnectionTypePostgres:
		store, err = postgres.Init(
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
		handler = handlers.NewHTTPHandler(store)
	case ConnectionTypeRedis:
		store, err = postgres.Init(
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
		redisClient := redis.NewClient(&redis.Options{
			Addr: config.Redis_URL,
		})
		store, err = rediscachedstorage.Init(redisClient, store)
		if err != nil {
			log.Fatalf("can't init postgres connection: %s", err.Error())
		}
		rateLimitFactory := ratelimit.NewFactory(redisClient)
		handler = handlers.NewHTTPHandlerCached(store, rateLimitFactory)
	case ConnectionTypeInMemory:
		store = in_memory.Init()
		handler = handlers.NewHTTPHandler(store)
	default:
		store, err = postgres.Init(
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
		redisClient := redis.NewClient(&redis.Options{
			Addr: config.Redis_URL,
		})
		store, err = rediscachedstorage.Init(redisClient, store)
		if err != nil {
			log.Fatalf("can't init postgres connection: %s", err.Error())
		}
		rateLimitFactory := ratelimit.NewFactory(redisClient)
		handler = handlers.NewHTTPHandlerCached(store, rateLimitFactory)
	}

	r.HandleFunc("/{shortUrl:\\w{10}}", handler.HandleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/urls", handler.HandlePostUrl)

	return &http.Server{
		Handler:      r,
		Addr:         ":8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func main() {
	srv := NewServer()
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
