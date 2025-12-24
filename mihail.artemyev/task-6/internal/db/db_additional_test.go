package db_test

import (
	"database/sql"
	"errors"
	"testing"

	db "task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errSomeRows = errors.New("some rows error")
	errDBFail   = errors.New("db fail")
)

func TestRealDB_QueryAdapter(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	r := &realDB{dbConn}
	gotRows, err := r.Query("SELECT name FROM users")
	require.NoError(t, err)

	defer gotRows.Close()

	var name string

	require.True(t, gotRows.Next())
	require.NoError(t, gotRows.Scan(&name))
	assert.Equal(t, "Alice", name)

	require.NoError(t, gotRows.Err())

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestNewDBService(t *testing.T) {
	t.Parallel()

	mockDB := new(MockDatabase)
	svc := db.New(mockDB)
	assert.Equal(t, mockDB, svc.DB)
}

func TestGetNames_FinalRowsError_NotScan(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errSomeRows)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(&realDB{dbConn})
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_FinalRowsError_Scan(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errScanError)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(&realDB{dbConn})
	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows scanning")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRealDB_QueryErrorPropagates(t *testing.T) {
	t.Parallel()

	mockDB := new(MockDatabase)
	mockDB.On("Query", "SELECT name FROM users").Return((*sql.Rows)(nil), errDBFail)

	service := db.New(mockDB)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")
	mockDB.AssertExpectations(t)
}
