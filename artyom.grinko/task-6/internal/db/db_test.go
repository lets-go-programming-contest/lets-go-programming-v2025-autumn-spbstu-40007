package db_test

import (
	"database/sql"
	"fmt"
	"testing"

	"task-6/internal/db"
	"task-6/internal/functionals"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var singersNames = []string{
	"Édith Piaf",
	"Françoise Hardy",
	"Mylène Farmer",
	"Édith Piaf",
	"Joe Dassin",
}

func TestGetNames(t *testing.T) {
	t.Run("queryFailure", func(t *testing.T) {
		t.Parallel()

		database, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer database.Close()
		mock.
			ExpectQuery("SELECT name FROM users").
			WillReturnError(sql.ErrNoRows)

		databaseService := db.New(database)
		names, err := databaseService.GetNames()
		require.Empty(t, names)
		require.NotNil(t, err)
	})

	t.Run("scanFailure", func(t *testing.T) {
		t.Parallel()

		database, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer database.Close()
		rows := mock.
			NewRows([]string{"name"}).
			AddRow(nil)
		mock.
			ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows).
			WillReturnError(nil)

		databaseService := db.New(database)
		names, err := databaseService.GetNames()
		require.Empty(t, names)
		require.NotNil(t, err)
	})

	t.Run("rowsFailure", func(t *testing.T) {
		t.Parallel()

		database, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer database.Close()
		rows := mock.
			NewRows([]string{"name"})
		functionals.Iter(func(name string) {
			rows.AddRow(name)
		}, singersNames)
		rows.RowError(0, fmt.Errorf("some erreur"))
		mock.
			ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		databaseService := db.New(database)
		names, err := databaseService.GetNames()
		require.Empty(t, names)
		require.NotNil(t, err)
	})

	t.Run("querySuccess", func(t *testing.T) {
		t.Parallel()

		database, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer database.Close()
		rows := mock.
			NewRows([]string{"name"})
		functionals.Iter(func(name string) {
			rows.AddRow(name)
		}, singersNames)

		mock.
			ExpectQuery("SELECT name FROM users").
			WillReturnRows(rows)

		databaseService := db.New(database)
		names, err := databaseService.GetNames()
		require.Equal(t, names, singersNames)
		require.Nil(t, err)
	})
}

func TestUniqueNames(t *testing.T) {
	t.Run("queryFailure", func(t *testing.T) {
		t.Parallel()

		database, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer database.Close()
		mock.
			ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(sql.ErrNoRows)

		databaseService := db.New(database)
		names, err := databaseService.GetUniqueNames()
		require.Empty(t, names)
		require.NotNil(t, err)
	})

	t.Run("scanFailure", func(t *testing.T) {
		t.Parallel()

		database, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer database.Close()
		rows := mock.
			NewRows([]string{"name"}).
			AddRow(nil)
		mock.
			ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows).
			WillReturnError(nil)

		databaseService := db.New(database)
		names, err := databaseService.GetUniqueNames()
		require.Empty(t, names)
		require.NotNil(t, err)
	})

	t.Run("rowsFailure", func(t *testing.T) {
		t.Parallel()

		database, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer database.Close()
		rows := mock.
			NewRows([]string{"name"})
		functionals.Iter(func(name string) {
			rows.AddRow(name)
		}, singersNames)
		rows.RowError(0, fmt.Errorf("some erreur"))
		mock.
			ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(rows)

		databaseService := db.New(database)
		names, err := databaseService.GetUniqueNames()
		require.Empty(t, names)
		require.NotNil(t, err)
	})

	t.Run("querySuccess", func(t *testing.T) {
		t.Parallel()

		database, mock, err := sqlmock.New()
		require.Nil(t, err)
		defer database.Close()
		rows := mock.
			NewRows([]string{"name"})
		functionals.Iter(func(name string) {
			rows.AddRow(name)
		}, singersNames)
		expectedRows := sqlmock.NewRows([]string{"name"})
		functionals.Iter(func(name string) {
			expectedRows.AddRow(name)
		}, functionals.Unique(singersNames))

		mock.
			ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(expectedRows)

		databaseService := db.New(database)
		names, err := databaseService.GetUniqueNames()
		require.Equal(t, names, functionals.Unique(singersNames))
		require.Nil(t, err)
	})
}
