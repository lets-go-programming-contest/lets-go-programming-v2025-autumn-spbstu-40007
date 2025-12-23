package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBService_GetNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := New(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("db query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errors.New("connection failed"))

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
	})

	t.Run("rows scan error", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(nil)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("rows error after iteration", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errors.New("row error"))

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows error")
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := New(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Charlie").
			AddRow("Dave")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Charlie", "Dave"}, names)
	})

	t.Run("db query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errors.New("query failed"))

		names, err := service.GetUniqueNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
	})

	t.Run("rows scanning error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Name").
			RowError(0, errors.New("corrupt data"))

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows error")
	})
}
