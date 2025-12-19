package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDatabaseNameOperations(t *testing.T) {

	t.Run("RetrieveAllUserNames", func(t *testing.T) {
		database, mockObj, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Mock creation failed: %v", err)
		}
		defer database.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("John Doe").
			AddRow("Jane Smith")
		mockObj.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)

		dbService := New(database)
		resultNames, err := dbService.GetNames()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(resultNames) != 2 {
			t.Fatalf("Expected 2 names, got %d", len(resultNames))
		}

		if resultNames[0] != "John Doe" || resultNames[1] != "Jane Smith" {
			t.Errorf("Name mismatch. Got: %v", resultNames)
		}

		if err := mockObj.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("HandleQueryFailure", func(t *testing.T) {
		database, mockObj, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Mock creation failed: %v", err)
		}
		defer database.Close()

		mockObj.ExpectQuery(`SELECT name FROM users`).
			WillReturnError(sqlmock.ErrCancelled)

		dbService := New(database)
		resultNames, err := dbService.GetNames()

		if err == nil {
			t.Error("Expected error but got none")
		}

		if resultNames != nil {
			t.Error("Expected nil result on query failure")
		}

		if err := mockObj.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("HandleRowScanError", func(t *testing.T) {
		database, mockObj, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Mock creation failed: %v", err)
		}
		defer database.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mockObj.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)

		dbService := New(database)
		resultNames, err := dbService.GetNames()

		if err == nil {
			t.Error("Expected scan error but got none")
		}

		if resultNames != nil {
			t.Error("Expected nil result on scan error")
		}

		if err := mockObj.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}

func TestUniqueNameOperations(t *testing.T) {

	t.Run("RetrieveDistinctNames", func(t *testing.T) {
		database, mockObj, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Mock creation failed: %v", err)
		}
		defer database.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob").
			AddRow("Alice")
		mockObj.ExpectQuery(`SELECT DISTINCT name FROM users`).WillReturnRows(rows)

		dbService := New(database)
		uniqueNames, err := dbService.GetUniqueNames()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(uniqueNames) != 3 {
			t.Fatalf("Expected 3 names, got %d", len(uniqueNames))
		}

		if err := mockObj.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("HandleDistinctQueryFailure", func(t *testing.T) {
		database, mockObj, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Mock creation failed: %v", err)
		}
		defer database.Close()

		mockObj.ExpectQuery(`SELECT DISTINCT name FROM users`).
			WillReturnError(sqlmock.ErrCancelled)

		dbService := New(database)
		uniqueNames, err := dbService.GetUniqueNames()

		if err == nil {
			t.Error("Expected error but got none")
		}

		if uniqueNames != nil {
			t.Error("Expected nil result on query failure")
		}

		if err := mockObj.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})

	t.Run("TestEmptyResultSet", func(t *testing.T) {
		database, mockObj, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Mock creation failed: %v", err)
		}
		defer database.Close()

		rows := sqlmock.NewRows([]string{"name"})
		mockObj.ExpectQuery(`SELECT DISTINCT name FROM users`).WillReturnRows(rows)

		dbService := New(database)
		uniqueNames, err := dbService.GetUniqueNames()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(uniqueNames) != 0 {
			t.Errorf("Expected empty slice, got %v", uniqueNames)
		}

		if err := mockObj.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}
	})
}
