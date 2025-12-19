package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	mydb "task-6/internal/db"
)

var (
	errSyntax        = errors.New("syntax error")
	errConnection    = errors.New("connection failed")
	errRowsIteration = errors.New("rows iteration failed")
	errRows          = errors.New("rows error")
)

// Тесты для GetNames разделены на отдельные функции

func TestDBService_GetNames_SuccessMultipleRows(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		AddRow("Charlie")
	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := mydb.New(db)
	names, err := service.GetNames()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(names) != 3 {
		t.Errorf("expected 3 names, got %d", len(names))
	}

	if names[0] != "Alice" || names[2] != "Charlie" {
		t.Errorf("names mismatch: %v", names)
	}
}

func TestDBService_GetNames_SuccessEmptyResult(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := mydb.New(db)
	names, err := service.GetNames()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(names) != 0 {
		t.Errorf("expected empty slice, got %v", names)
	}
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT name FROM users$").
		WillReturnError(errSyntax)

	service := mydb.New(db)
	names, err := service.GetNames()

	if err == nil {
		t.Error("expected error, got nil")
	}

	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}

	if err.Error() != "db query: syntax error" {
		t.Errorf("error message mismatch: %v", err)
	}
}

func TestDBService_GetNames_ScanErrorOnNullValue(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := mydb.New(db)
	names, err := service.GetNames()

	if err == nil {
		t.Error("expected error for NULL value, got nil")
	}

	if names != nil {
		t.Errorf("expected nil names for NULL value, got %v", names)
	}
}

func TestDBService_GetNames_RowsError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		CloseError(errRowsIteration)
	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := mydb.New(db)
	names, err := service.GetNames()

	if err == nil {
		t.Error("expected rows error, got nil")
	}

	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}

	if err.Error() != "rows error: rows iteration failed" {
		t.Errorf("error message mismatch: %v", err)
	}
}

// Тесты для GetUniqueNames тоже разделим

func TestDBService_GetUniqueNames_SuccessWithDuplicates(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		AddRow("Alice")
	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := mydb.New(db)
	names, err := service.GetUniqueNames()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(names) != 3 {
		t.Errorf("expected 3 names, got %d", len(names))
	}
}

func TestDBService_GetUniqueNames_SuccessSingleRow(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("SingleUser")
	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := mydb.New(db)
	names, err := service.GetUniqueNames()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(names) != 1 || names[0] != "SingleUser" {
		t.Errorf("expected [SingleUser], got %v", names)
	}
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
		WillReturnError(errConnection)

	service := mydb.New(db)
	names, err := service.GetUniqueNames()

	if err == nil {
		t.Error("expected error, got nil")
	}

	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
}

func TestDBService_GetUniqueNames_ScanErrorOnNull(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := mydb.New(db)
	names, err := service.GetUniqueNames()

	if err == nil {
		t.Error("expected error for NULL, got nil")
	}

	if names != nil {
		t.Errorf("expected nil names for NULL, got %v", names)
	}
}

func TestDBService_GetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		CloseError(errRows)
	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := mydb.New(db)
	names, err := service.GetUniqueNames()

	if err == nil {
		t.Error("expected rows error, got nil")
	}

	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
}

func TestDBService_New(t *testing.T) {
	t.Parallel()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := mydb.New(db)
	if service.DB != db {
		t.Errorf("expected DB to be %v, got %v", db, service.DB)
	}
}
