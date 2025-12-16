package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		service := New(db)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		service := New(db)
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("fail"))

		_, err = service.GetNames()
		assert.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		service := New(db)
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err = service.GetNames()
		assert.Error(t, err)
	})

	t.Run("rows_error_at_end", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		service := New(db)
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			CloseError(errors.New("rows error"))

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err = service.GetNames()
		assert.Error(t, err)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		service := New(db)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Unique")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Unique"}, res)
	})

	t.Run("query_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		service := New(db)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("fail"))

		_, err = service.GetUniqueNames()
		assert.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		service := New(db)
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err = service.GetUniqueNames()
		assert.Error(t, err)
	})

	t.Run("rows_error_unique", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		service := New(db)
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("UniqueUser").
			CloseError(errors.New("rows error"))

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err = service.GetUniqueNames()
		assert.Error(t, err)
	})
}
