package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"task-6/internal/db"
)

func createService(t *testing.T) (db.DBService, sqlmock.Sqlmock) {
	t.Helper()

	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open mock db: %v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	return db.New(conn), mock
}

func TestDBService_GetNames_AllCases(t *testing.T) {
	service, mock := createService(t)

	t.Run("happy path", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alex").
			AddRow("Maria")

		mock.ExpectQuery("^SELECT name FROM users$").
			WillReturnRows(rows)

		result, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Alex", "Maria"}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query fails", func(t *testing.T) {
		mock.ExpectQuery("^SELECT name FROM users$").
			WillReturnError(errors.New("db offline"))

		_, err := service.GetNames()
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan fails", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery("^SELECT name FROM users$").
			WillReturnRows(rows)

		_, err := service.GetNames()
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alex").
			RowError(0, errors.New("broken row"))

		mock.ExpectQuery("^SELECT name FROM users$").
			WillReturnRows(rows)

		_, err := service.GetNames()
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames_AllCases(t *testing.T) {
	service, mock := createService(t)

	t.Run("happy path", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alex").
			AddRow("Kate")

		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
			WillReturnRows(rows)

		result, err := service.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Alex", "Kate"}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query fails", func(t *testing.T) {
		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
			WillReturnError(errors.New("db offline"))

		_, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan fails", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
			WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alex").
			RowError(0, errors.New("broken row"))

		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
			WillReturnRows(rows)

		_, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
