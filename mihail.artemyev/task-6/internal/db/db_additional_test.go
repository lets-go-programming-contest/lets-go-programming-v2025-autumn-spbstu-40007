package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestRealDB_QueryAdapter(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	r := &realDB{db}
	gotRows, err := r.Query("SELECT name FROM users")
	require.NoError(t, err)
	defer gotRows.Close()

	var name string
	require.True(t, gotRows.Next())
	require.NoError(t, gotRows.Scan(&name))
	assert.Equal(t, "Alice", name)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestNewDBService(t *testing.T) {
	mockDB := new(MockDatabase)
	svc := New(mockDB)
	assert.Equal(t, mockDB, svc.DB)
}

func TestGetNames_FinalRowsError_NotScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("some rows error"))
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := New(&realDB{db})
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_FinalRowsError_Scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("scan error"))
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := New(&realDB{db})
	names, err := service.GetUniqueNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows scanning")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRealDB_QueryErrorPropagates(t *testing.T) {
	mockDB := new(MockDatabase)
	mockDB.On("Query", "SELECT name FROM users").Return((*sql.Rows)(nil), errors.New("db fail"))

	service := New(mockDB)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")
	mockDB.AssertExpectations(t)
}
