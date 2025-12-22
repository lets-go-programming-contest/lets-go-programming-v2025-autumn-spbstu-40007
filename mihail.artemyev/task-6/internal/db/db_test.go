package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockDatabase для тестирования
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	ret := m.Called(query, args)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*sql.Rows), ret.Error(1)
}

// Test GetNames - успешное получение имён
func TestGetNames_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		AddRow("Charlie")

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := New(&realDB{db})
	names, err := service.GetNames()

	assert.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob", "Charlie"}, names)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test GetNames - ошибка при выполнении запроса
func TestGetNames_QueryError(t *testing.T) {
	mockDB := new(MockDatabase)
	mockDB.On("Query", "SELECT name FROM users").Return(nil, errors.New("database error"))

	service := DBService{DB: mockDB}
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")
	mockDB.AssertExpectations(t)
}

// Test GetNames - ошибка при сканировании строк
func TestGetNames_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("scan error"))

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := New(&realDB{db})
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows scanning")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test GetNames - ошибка в результатах
func TestGetNames_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(1, errors.New("rows error"))

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := New(&realDB{db})
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test GetNames - пустой результат
func TestGetNames_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := New(&realDB{db})
	names, err := service.GetNames()

	assert.NoError(t, err)
	assert.Equal(t, []string{}, names)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test GetUniqueNames - успешное получение уникальных имён
func TestGetUniqueNames_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := New(&realDB{db})
	names, err := service.GetUniqueNames()

	assert.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob"}, names)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test GetUniqueNames - ошибка при выполнении запроса
func TestGetUniqueNames_QueryError(t *testing.T) {
	mockDB := new(MockDatabase)
	mockDB.On("Query", "SELECT DISTINCT name FROM users").Return(nil, errors.New("db error"))

	service := DBService{DB: mockDB}
	names, err := service.GetUniqueNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")
	mockDB.AssertExpectations(t)
}

// Test GetUniqueNames - ошибка при сканировании
func TestGetUniqueNames_ScanError(t *testing.T) {
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
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test GetUniqueNames - ошибка в результатах
func TestGetUniqueNames_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(1, errors.New("rows error"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := New(&realDB{db})
	names, err := service.GetUniqueNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test GetUniqueNames - пустой результат
func TestGetUniqueNames_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := New(&realDB{db})
	names, err := service.GetUniqueNames()

	assert.NoError(t, err)
	assert.Equal(t, []string{}, names)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Адаптер для работы с *sql.DB
type realDB struct {
	db *sql.DB
}

func (r *realDB) Query(query string, args ...any) (*sql.Rows, error) {
	return r.db.Query(query, args...)
}
