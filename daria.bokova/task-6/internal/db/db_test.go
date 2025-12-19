package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()
	t.Run("success - multiple rows", func(t *testing.T) {
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

		service := New(db)
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
	})

	t.Parallel()

	t.Run("success - empty result", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"})
		mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

		service := New(db)
		names, err := service.GetNames()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(names) != 0 {
			t.Errorf("expected empty slice, got %v", names)
		}
	})

	t.Parallel()

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		mock.ExpectQuery("^SELECT name FROM users$").
			WillReturnError(errors.New("syntax error"))

		service := New(db)
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
	})

	t.Parallel()

	t.Run("scan error on null value", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		// NULL значение вызовет ошибку при сканировании в string
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

		service := New(db)
		names, err := service.GetNames()

		// NULL в string должно вызвать ошибку
		if err == nil {
			t.Error("expected error for NULL value, got nil")
		}

		if names != nil {
			t.Errorf("expected nil names for NULL value, got %v", names)
		}
	})

	t.Parallel()

	t.Run("rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			CloseError(errors.New("rows iteration failed"))
		mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

		service := New(db)
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
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()
	t.Run("success - with duplicates", func(t *testing.T) {
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

		service := New(db)
		names, err := service.GetUniqueNames()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(names) != 3 {
			t.Errorf("expected 3 names, got %d", len(names))
		}
	})

	t.Parallel()

	t.Run("success - single row", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("SingleUser")
		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

		service := New(db)
		names, err := service.GetUniqueNames()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(names) != 1 || names[0] != "SingleUser" {
			t.Errorf("expected [SingleUser], got %v", names)
		}
	})

	t.Parallel()

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
			WillReturnError(errors.New("connection failed"))

		service := New(db)
		names, err := service.GetUniqueNames()

		if err == nil {
			t.Error("expected error, got nil")
		}
		if names != nil {
			t.Errorf("expected nil names, got %v", names)
		}
	})

	t.Parallel()

	t.Run("scan error on null", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

		service := New(db)
		names, err := service.GetUniqueNames()

		if err == nil {
			t.Error("expected error for NULL, got nil")
		}

		if names != nil {
			t.Errorf("expected nil names for NULL, got %v", names)
		}
	})

	t.Parallel()

	t.Run("rows error for unique", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			CloseError(errors.New("rows error"))
		mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

		service := New(db)
		names, err := service.GetUniqueNames()

		if err == nil {
			t.Error("expected rows error, got nil")
		}

		if names != nil {
			t.Errorf("expected nil names, got %v", names)
		}
	})
}

func TestDBService_New(t *testing.T) {
	t.Parallel()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	service := New(db)
	if service.DB != db {
		t.Errorf("expected DB to be %v, got %v", db, service.DB)
	}
}
