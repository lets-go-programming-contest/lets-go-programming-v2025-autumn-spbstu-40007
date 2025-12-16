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
	errTestQuery = errors.New("fail")
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
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errTestQuery)

		_, err = service.GetNames()
		assert.Error(t, err)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		_, err = service.GetNames()
		assert.Error(t, err)
	})

	t.Run("rows_error_at_end", func(t *testing.T) {
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
		assert.Error(t, err)
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
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Unique")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		res, err := service.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Unique"}, res)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errTestQuery)

		_, err = service.GetUniqueNames()
		assert.Error(t, err)
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

	t.Run("rows_error_unique", func(t *testing.T) {
		t.Parallel()
		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("UniqueUser").
			CloseError(errTestRows)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		_, err = service.GetUniqueNames()
		assert.Error(t, err)
	})
}
