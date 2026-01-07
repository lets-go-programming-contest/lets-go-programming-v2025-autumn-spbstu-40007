package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"task-6/internal/db"
)

func prepareDB(t *testing.T) (db.DBService, sqlmock.Sqlmock) {
	t.Helper()

	conn, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() { _ = conn.Close() })

	return db.New(conn), mock
}

func TestDBService_GetNames(t *testing.T) {
	service, mock := prepareDB(t)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errors.New("query failure"))

		_, err := service.GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errors.New("row problem"))

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		_, err := service.GetNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	service, mock := prepareDB(t)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Charlie")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Charlie"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errors.New("query failure"))

		_, err := service.GetUniqueNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(nil)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errors.New("row problem"))

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
