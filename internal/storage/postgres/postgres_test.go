package postgres

import (
	"context"
	"reflect"
	"testing"

	"github.com/GeorgeShibanin/ozon_test/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
)

func TestGetURL(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mock.Close(context.Background())

	link := &Link{Key: "my-new-link", URL: "https://www.vk.com"}
	row := pgxmock.NewRows([]string{"id", "url"}).AddRow(link.Key, link.URL)
	mock.
		ExpectQuery("SELECT id, url FROM links WHERE").
		WithArgs(storage.ShortedURL(link.Key)).
		WillReturnRows(row)

	repo := initConnection(mock)

	res, err := repo.GetURL(context.Background(), storage.ShortedURL(link.Key))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
		return
	}
	if !reflect.DeepEqual(storage.URL(link.URL), res) {
		t.Errorf("results not match, want %v, have %v", link.URL, res)
		return
	}
}

func TestPutURL(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mock.Close(context.Background())

	link := &Link{Key: "my-new-link", URL: "https://www.vk.com"}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, url FROM links WHERE").
		WithArgs(storage.URL(link.URL)).
		WillReturnError(pgx.ErrNoRows)
	mock.ExpectQuery("SELECT id, url FROM links WHERE").
		WithArgs(storage.ShortedURL(link.Key)).
		WillReturnError(pgx.ErrNoRows)
	mock.ExpectExec("INSERT INTO links").
		WithArgs(storage.ShortedURL(link.Key), storage.URL(link.URL)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	repo := initConnection(mock)
	res, err := repo.PutURL(context.Background(), storage.ShortedURL(link.Key), storage.URL(link.URL))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
		return
	}
	if !reflect.DeepEqual(storage.ShortedURL(link.Key), res) {
		t.Errorf("results not match, want %v, have %v", storage.ShortedURL(link.URL), res)
		return
	}
}
