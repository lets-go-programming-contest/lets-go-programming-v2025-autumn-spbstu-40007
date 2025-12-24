package db

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// realDBForTests — адаптер для *sql.DB в тестах пакета db.
type realDBForTests struct {
	db *sql.DB
}

type nilRowsDB struct{}

func (n *nilRowsDB) Query(string, ...any) (*sql.Rows, error) {
	return nil, nil
}

func Test_queryStrings_NilRowsWithoutError(t *testing.T) {
	t.Parallel()

	svc := New(&nilRowsDB{})

	res, err := svc.queryStrings("SELECT name")

	require.Error(t, err)
	require.Nil(t, res)
	require.EqualError(t, err, "db returned nil rows without error")
}

func (r *realDBForTests) Query(query string, args ...any) (*sql.Rows, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return rows, fmt.Errorf("db.Query: %w", err)
	}
	return rows, nil
}

func Test_queryStrings_FinalRowsError_ScanKeyword(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.
		NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("scan failure"))

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	svc := New(&realDBForTests{sqlDB})

	res, err := svc.queryStrings("SELECT name")
	require.Error(t, err)
	require.EqualError(t, err, "rows scanning: scan failure")
	require.Nil(t, res)

	require.NoError(t, mock.ExpectationsWereMet())
}

func Test_queryStrings_FinalRowsError_NoScan(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.
		NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("rows failure"))

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	svc := New(&realDBForTests{sqlDB})

	res, err := svc.queryStrings("SELECT name")
	require.Error(t, err)
	require.EqualError(t, err, "rows error: rows failure")
	require.Nil(t, res)

	require.NoError(t, mock.ExpectationsWereMet())
}
