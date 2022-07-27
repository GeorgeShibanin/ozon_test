package handlers

import (
	"context"
	"fmt"
	storage2 "github.com/GeorgeShibanin/ozon_test/internal/storage"

	"github.com/GeorgeShibanin/ozon_test/internal/storage/postgres"
	"github.com/pashagolub/pgxmock"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	db, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close(context.Background())

	sl := &postgres.Link{Key: "my-new-link", URL: "https://www.vk.com"}
	repo := &postgres.Storage{Conn: db}
	row := pgxmock.NewRows([]string{"id", "url"}).AddRow(sl.Key, sl.URL)
	// Basic test
	db.
		ExpectQuery("SELECT id, url FROM links WHERE").
		WithArgs(sl.Key, sl.URL).
		WillReturnRows(row)

	link, err := repo.GetURL(context.Background(), storage2.ShortedURL(sl.Key))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if err := db.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
		return
	}
	if !reflect.DeepEqual(link, sl.URL) {
		t.Errorf("results not match, want %v, have %v", sl.URL, link)
		return
	}

	// Test with error
	db.
		ExpectQuery("SELECT id, url FROM links WHERE").
		WithArgs(sl.Key, sl.URL).
		WillReturnError(fmt.Errorf("db_error"))

	link, err = repo.GetURL(context.Background(), storage2.ShortedURL(sl.Key))
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := db.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
		return
	}
	if !reflect.DeepEqual(link, sl.URL) {
		t.Errorf("results not match, want %v, have %v", sl.URL, link)
		return
	}
}

//func TestSet(t *testing.T) {
//	db, err := pgxmock.NewConn()
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	defer db.Close(context.Background())
//
//	sl := &postgres.Link{Key: "my-new-link", URL: "https://www.vk.com"}
//	repo := &postgres.Storage{Conn: db}
//	// Basic test
//	db.
//		ExpectExec("INSERT INTO links (`id`, `url`) VALUES").
//		WithArgs(sl.Key, sl.URL).
//		WillReturnResult(pgxmock.NewResult("INSERT", 1))
//
//	id, err := repo.PutURL(context.Background(), storage2.ShortedURL(sl.Key), storage2.URL(sl.URL))
//	if err != nil {
//		t.Errorf(err.Error())
//		return
//	}
//	if id != 1 {
//		t.Errorf("bad id: want %v, got %v", 1, id)
//		return
//	}
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("unmet expectation error: %s", err)
//		return
//	}
//
//	// Test with err
//	mock.
//		ExpectExec("INSERT INTO links (`id`, `link`) VALUES").
//		WithArgs(sl.ID, sl.Link).
//		WillReturnError(fmt.Errorf("bad query"))
//
//	id, err = repo.Set(sl)
//	if err == nil {
//		t.Errorf("expected error, got nil")
//		return
//	}
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("unmet expectation error: %s", err)
//		return
//	}
//}
