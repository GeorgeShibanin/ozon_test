package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"os"
	"ozon_test/internal/handlers"
	"ozon_test/internal/storage"
	"ozon_test/internal/storage/postgresql"
	"time"
)

func NewServer() *http.Server {
	r := mux.NewRouter()
	handler := &handlers.HTTPHandler{}
	storageType := os.Getenv("STORAGE_MODE")

	if storageType == "postgres" {
		conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())
		postgresStorage := postgresql.NewStorage(conn)
		handler = &handlers.HTTPHandler{
			Storage: postgresStorage,
		}
	} else if storageType == "inmemory" {
		handler = &handlers.HTTPHandler{
			StorageInMemory: make(map[storage.URLKey]storage.ShortedURL),
		}
	} else if storageType == "cached" {
	}

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
