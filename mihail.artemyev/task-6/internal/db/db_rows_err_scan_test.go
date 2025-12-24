package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestQueryStrings_FinalRowsError_ScanKeyword(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("scan failure"))

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	service := New(db)

	res, err := service.queryStrings("SELECT name")

	require.Error(t, err)
	require.EqualError(t, err, "rows scanning: scan failure")
	require.Nil(t, res)

	require.NoError(t, mock.ExpectationsWereMet())
}
