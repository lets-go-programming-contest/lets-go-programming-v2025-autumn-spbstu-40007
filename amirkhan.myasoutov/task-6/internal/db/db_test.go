package db_test

import (
	"context"
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ami0-0/task-6/internal/db"
)

var (
	errDB  = errors.New("connection error")
	errRow = errors.New("rows iteration error")
)

func TestDBService_GetNames(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
		names, err := service.GetNames(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})
	
	t.Run("empty result", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
		names, err := service.GetNames(ctx)
		require.NoError(t, err)
		assert.Empty(t, names)
	})

	t.Run("query error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDB)
		_, err = service.GetNames(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query")
	})

	t.Run("scan error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
		_, err = service.GetNames(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("rows error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRow)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
		_, err = service.GetNames(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	ctx := context.Background()
	
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		names, err := service.GetUniqueNames(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})
	
	t.Run("empty result", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		names, err := service.GetUniqueNames(ctx)
		require.NoError(t, err)
		assert.Empty(t, names)
	})

	t.Run("query error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errDB)
		_, err = service.GetUniqueNames(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query")
	})

	t.Run("scan error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		_, err = service.GetUniqueNames(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("rows error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() { _ = mockDB.Close() }()
		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRow)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		_, err = service.GetUniqueNames(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
	})
}