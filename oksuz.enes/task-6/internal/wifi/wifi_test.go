package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBService_GetNames(t *testing.T) {
	tests := []struct {
		name          string
		mockBehavior  func(mock sqlmock.Sqlmock)
		expectedNames []string
		expectError   bool
	}{
		{
			name: "Success",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("John").AddRow("Jane")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"John", "Jane"},
			expectError:   false,
		},
		{
			name: "QueryError",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("db error"))
			},
			expectedNames: nil,
			expectError:   true,
		},
		{
			name: "ScanError",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
		},
		{
			name: "RowsError",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("John").RowError(0, errors.New("row error"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockBehavior(mock)

			service := New(db)
			got, err := service.GetNames()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNames, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	tests := []struct {
		name          string
		mockBehavior  func(mock sqlmock.Sqlmock)
		expectedNames []string
		expectError   bool
	}{
		{
			name: "Success",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: []string{"Alice"},
			expectError:   false,
		},
		{
			name: "QueryError",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("fail"))
			},
			expectedNames: nil,
			expectError:   true,
		},
		{
			name: "ScanError",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
		},
		{
			name: "RowsError",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Bob").RowError(0, errors.New("err"))
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockBehavior(mock)

			service := New(db)
			got, err := service.GetUniqueNames()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNames, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
