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
		t.Fatalf("error opening mock: %s", err)
	}
	defer db.Close()

	service := New(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").AddRow("bob")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"alice", "bob"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("db error"))

		names, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").RowError(0, errors.New("row error"))
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, names)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening mock: %s", err)
	}
	defer db.Close()

	service := New(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("alice")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"alice"}, names)
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("db error"))

		names, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Nil(t, names)
	})

	t.Run("rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").RowError(0, errors.New("row error"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		assert.Error(t, err)
		assert.Nil(t, names)
	})
}
