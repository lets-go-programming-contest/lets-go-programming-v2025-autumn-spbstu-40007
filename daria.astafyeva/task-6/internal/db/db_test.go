package db_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/itsdasha/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	queryErr = errors.New("query failed")
	scanErr  = errors.New("scan failed")
	rowsErr  = errors.New("rows iteration failed")
)

func TestDBService_New(t *testing.T) {
	t.Parallel()

	dbMock, _, _ := sqlmock.New()
	defer dbMock.Close()

	service := db.New(dbMock)

	assert.NotNil(t, service)
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Len(t, names, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnError(queryErr)

		service := db.New(dbMock)
		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows scanning")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, rowsErr)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Alice")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		service := db.New(dbMock)
		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.Len(t, names, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnError(queryErr)

		service := db.New(dbMock)
		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
