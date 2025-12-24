package db_test

import (
	"testing"

	db "task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestQueryStrings_FinalRowsError_NoScan(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.
		NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errRowsError)

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	service := db.New(&realDB{sqlDB})

	result, err := service.GetNames()

	require.Error(t, err)
	require.EqualError(t, err, "rows error: "+errRowsError.Error())
	require.Nil(t, result)

	require.NoError(t, mock.ExpectationsWereMet())
}
