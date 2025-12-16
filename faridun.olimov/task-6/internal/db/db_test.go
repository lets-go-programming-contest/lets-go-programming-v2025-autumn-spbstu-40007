package db_test

import (
	"errors"
	"testing"

	"task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errTestQuery = errors.New("query fail")
	errTestRows  = errors.New("rows error")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errTestQuery)

		_, err = service.GetNames()
		assert.ErrorIs(t, err, errTestQuery)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		// Передаем строку, которая не может быть корректно обработана (например, nil в NOT NULL поле)
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err = service.GetNames()
		assert.Error(t, err)
	})

	t.Run("rows_err", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			CloseError(errTestRows)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err = service.GetNames()
		assert.ErrorIs(t, err, errTestRows)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errTestQuery)

		_, err = service.GetUniqueNames()
		assert.ErrorIs(t, err, errTestQuery)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err = service.GetUniqueNames()
		assert.Error(t, err)
	})

	t.Run("rows_err", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			CloseError(errTestRows)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err = service.GetUniqueNames()
		assert.ErrorIs(t, err, errTestRows)
	})
}
