package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	db "github.com/ksuah/task-6/internal/db"
)

var (
	errDatabaseDown  = errors.New("database connection failed")
	errRowProcessing = errors.New("row processing error")
)

func TestDBService_GetNames_OK(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	mockRows := sqlmock.NewRows([]string{"name"}).
		AddRow("ksuah").
		AddRow("pavel")

	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(mockRows)

	dbService := db.New(dbConn)

	result, err := dbService.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"ksuah", "pavel"}, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	mock.ExpectQuery("^SELECT name FROM users$").WillReturnError(errDatabaseDown)

	dbService := db.New(dbConn)

	result, err := dbService.GetNames()
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "db query:")
	require.ErrorIs(t, err, errDatabaseDown)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	mockRows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(mockRows)

	dbService := db.New(dbConn)

	result, err := dbService.GetNames()
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows scanning:")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_RowsError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	mockRows := sqlmock.NewRows([]string{"name"}).
		AddRow("ksuah").
		RowError(0, errRowProcessing)

	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(mockRows)

	dbService := db.New(dbConn)

	result, err := dbService.GetNames()
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows error:")
	require.ErrorIs(t, err, errRowProcessing)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_OK(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	mockRows := sqlmock.NewRows([]string{"name"}).
		AddRow("ksuah").
		AddRow("pavel")

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(mockRows)

	dbService := db.New(dbConn)

	result, err := dbService.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"ksuah", "pavel"}, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnError(errDatabaseDown)

	dbService := db.New(dbConn)

	result, err := dbService.GetUniqueNames()
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "db query:")
	require.ErrorIs(t, err, errDatabaseDown)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	mockRows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(mockRows)

	dbService := db.New(dbConn)

	result, err := dbService.GetUniqueNames()
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows scanning:")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	mockRows := sqlmock.NewRows([]string{"name"}).
		AddRow("ksuah").
		RowError(0, errRowProcessing)

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(mockRows)

	dbService := db.New(dbConn)

	result, err := dbService.GetUniqueNames()
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows error:")
	require.ErrorIs(t, err, errRowProcessing)

	require.NoError(t, mock.ExpectationsWereMet())
}
