package db_test

import (
	"errors"
	"testing"

	"github.com/tntkatz/task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errDB  = errors.New("db error")
	errRow = errors.New("row iteration error")
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, _ := sqlmock.New()

	defer mockDB.Close()

	service := db.New(mockDB)
	assert.NotNil(t, service)
	assert.Equal(t, mockDB, service.DB)
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()

		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()

		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDB)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()

		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()

		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errRow)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows error")
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()

		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()

		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errDB)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()

		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()

		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errRow)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, names)
	})
}
