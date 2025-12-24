package db_test

import (
	"database/sql"
	"testing"

	db "task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var d *sql.DB
	svc := db.New(d)
	require.Equal(t, d, svc.DB)
}

func TestGetNames_OK(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, res)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_QueryError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(errDBError)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetNames()

	require.Error(t, err)
	require.Nil(t, res)
}

func TestGetNames_ScanError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetNames()

	require.Error(t, err)
	require.Nil(t, res)
}

func TestGetNames_RowsError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errRowsError)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetNames()

	require.Error(t, err)
	require.Nil(t, res)
}
