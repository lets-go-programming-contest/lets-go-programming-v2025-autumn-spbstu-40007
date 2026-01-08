package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDB_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := New(db)

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
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := New(db)

		expectedErr := errors.New("connection failed")
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(expectedErr)

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db query:")
		assert.Contains(t, err.Error(), "connection failed")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := New(db)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(123)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning:")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errors.New("row error"))

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows error:")
		assert.Contains(t, err.Error(), "row error")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := New(db)

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
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := New(db)

		expectedErr := errors.New("syntax error")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(expectedErr)

		names, err := service.GetUniqueNames()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db query:")
		assert.Contains(t, err.Error(), "syntax error")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := New(db)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(456)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning:")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		service := New(db)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Bob").
			RowError(0, errors.New("distinct row error"))

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows error:")
		assert.Contains(t, err.Error(), "distinct row error")
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
