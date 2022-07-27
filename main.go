package main

import (
	"context"
	"fmt"
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
	switch storageMode {
	case ConnectionTypePostgres:
		store, err = postgres.Init(
			context.Background(),
			config.PostgresHost,
			config.PostgresUser,
			config.PostgresDB,
			"", // По-хорошему бы читать как переменную окружения
			config.PostgresPort,
		)
		if err != nil {
			log.Fatalf("can't init postgres connection: %s", err.Error())
		}
	case ConnectionTypeRedis:
	case ConnectionTypeInMemory:
		store = in_memory.Init()
	}

	handler := handlers.NewHTTPHandler(store)
	r.HandleFunc("/{shortUrl:\\w{10}}", handler.HandleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/urls", handler.HandlePostUrl)

	return &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func main() {
	srv := NewServer()
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
