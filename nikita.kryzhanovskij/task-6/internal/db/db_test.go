package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbpkg "nikita.kryzhanovskij/task-6/internal/db"
)

var (
	errConnectionFailed = errors.New("connection failed")
	errRowError         = errors.New("row error")
	errQueryFailed      = errors.New("query failed")
	errCorruptData      = errors.New("corrupt data")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = db.Close()
		})

		service := dbpkg.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("db query error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = db.Close()
		})

		service := dbpkg.New(db)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errConnectionFailed)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("rows scan error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = db.Close()
		})

		service := dbpkg.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(nil)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("rows error after iteration", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = db.Close()
		})

		service := dbpkg.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errRowError)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = db.Close()
		})

		service := dbpkg.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Charlie").
			AddRow("Dave")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Charlie", "Dave"}, names)
	})

	t.Run("db query error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = db.Close()
		})

		service := dbpkg.New(db)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errQueryFailed)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("rows scanning error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = db.Close()
		})

		service := dbpkg.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(nil)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = db.Close()
		})

		service := dbpkg.New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Name").
			RowError(0, errCorruptData)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})
}
