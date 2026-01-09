package db_test

import (
	"database/sql"
	"testing"

	db "task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errRowsError = sqlmock.ErrCancelled
	errDBError   = sqlmock.ErrCancelled
)

func TestNew(t *testing.T) {
	t.Parallel()

	var d *sql.DB
	svc := db.New(d)

	require.Equal(t, d, svc.DB)
}

func TestGetNames_OK(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, res)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_QueryError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(errDBError)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetNames()

	require.Error(t, err)
	require.Nil(t, res)
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetNames()

	require.Error(t, err)
	require.Nil(t, res)
}

func TestGetNames_RowsError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errRowsError)

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetNames()

	require.Error(t, err)
	require.Nil(t, res)
}

func TestGetNames_RowsErr_AfterIteration(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.
		NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errRowsError)

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})

	res, err := svc.GetNames()

	require.Error(t, err)
	require.EqualError(t, err, "rows error: "+errRowsError.Error())
	require.Nil(t, res)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_Empty(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetNames()

	require.NoError(t, err)
	assert.Empty(t, res)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_OK(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetUniqueNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, res)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnError(errDBError)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, res)
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, res)
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errRowsError)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, res)
}

func TestGetUniqueNames_Empty(t *testing.T) {
	t.Parallel()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	svc := db.New(&realDB{sqlDB})
	res, err := svc.GetUniqueNames()

	require.NoError(t, err)
	assert.Empty(t, res)
	require.NoError(t, mock.ExpectationsWereMet())
}
