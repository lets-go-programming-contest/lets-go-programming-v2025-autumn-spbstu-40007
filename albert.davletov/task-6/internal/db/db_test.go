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
	errConnectionFailed = errors.New("connection failed")
	errRowError         = errors.New("row error")
	errSyntaxError      = errors.New("syntax error")
	errDistinctRowError = errors.New("distinct row error")
)

func TestDB_GetNames(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer conn.Close()

		service := db.New(conn)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Bob").
			AddRow("Charlie")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Bob", "Charlie"}, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer conn.Close()

		service := db.New(conn)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errConnectionFailed)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query:")
		assert.Contains(t, err.Error(), errConnectionFailed.Error())
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer conn.Close()

		service := db.New(conn)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning:")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer conn.Close()

		service := db.New(conn)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errRowError)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error:")
		assert.Contains(t, err.Error(), errRowError.Error())
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer conn.Close()

		service := db.New(conn)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Bob").
			AddRow("Bob").
			AddRow("Charlie")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Bob", "Bob", "Charlie"}, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer conn.Close()

		service := db.New(conn)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errSyntaxError)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query:")
		assert.Contains(t, err.Error(), errSyntaxError.Error())
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer conn.Close()

		service := db.New(conn)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning:")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer conn.Close()

		service := db.New(conn)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Bob").
			RowError(0, errDistinctRowError)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error:")
		assert.Contains(t, err.Error(), errDistinctRowError.Error())
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
