package db_test

import (
	"testing"

	db "task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestQueryStrings_FinalRowsError_ScanKeyword(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.
		NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errScanError)

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	service := db.New(&realDB{sqlDB})

	names, err := service.GetNames()

	require.Error(t, err)
	require.EqualError(t, err, "rows scanning: "+errScanError.Error())
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}
