package db_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	db "task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errDatabaseError      = errors.New("database error")
	errScanError          = errors.New("scan error")
	errRowsError          = errors.New("rows error")
	errDBError            = errors.New("db error")
	errMockUnexpectedType = errors.New("mock returned unexpected type")
	errMockNilRows        = errors.New("mock returned nil rows without error")
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	ret := m.Called(append([]any{query}, args...)...)

	if ret.Get(0) == nil {
		if err := ret.Error(1); err != nil {
			return nil, fmt.Errorf("mock: %w", err)
		}

		return nil, errMockNilRows
	}

	got := ret.Get(0)

	rows, ok := got.(*sql.Rows)
	if !ok {
		return nil, errMockUnexpectedType
	}

	if err := ret.Error(1); err != nil {
		return rows, fmt.Errorf("mock: %w", err)
	}

	return rows, nil
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		AddRow("Charlie")

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(&realDB{dbConn})
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob", "Charlie"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB := new(MockDatabase)
	mockDB.On("Query", "SELECT name FROM users", mock.Anything).Return(nil, errDatabaseError)

	service := db.DBService{DB: mockDB}
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")
	mockDB.AssertExpectations(t)
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errScanError)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(&realDB{dbConn})
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "scan error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_RowsError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errRowsError)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(&realDB{dbConn})
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_Empty(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(&realDB{dbConn})
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(&realDB{dbConn})
	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB := new(MockDatabase)
	mockDB.On("Query", "SELECT DISTINCT name FROM users", mock.Anything).Return(nil, errDBError)

	service := db.DBService{DB: mockDB}
	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")
	mockDB.AssertExpectations(t)
}

func TestGetUniqueNames_ScanError(t *testing.T) {
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

func TestGetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errRowsError)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(&realDB{dbConn})
	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_Empty(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(&realDB{dbConn})
	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Equal(t, []string{}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

type realDB struct {
	db *sql.DB
}

func (r *realDB) Query(query string, args ...any) (*sql.Rows, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return rows, fmt.Errorf("db.Query: %w", err)
	}

	return rows, nil
}
