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

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		wantNames   []string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - multiple rows",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Charlie")
				m.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)
			},
			wantNames: []string{"Alice", "Bob", "Charlie"},
			wantErr:   false,
		},
		{
			name: "success - empty result",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				m.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)
			},
			wantNames: []string{},
			wantErr:   false,
		},
		{
			name: "query error",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("^SELECT name FROM users$").
					WillReturnError(errSyntax)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "syntax error",
		},
		{
			name: "scan error on null value",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				m.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "rows error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					CloseError(errRowsIteration)
				m.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "rows iteration failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			tt.setupMock(mock)
			service := mydb.New(db)
			names, err := service.GetNames()

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errContains != "" && err != nil {
					if err.Error() != "db query: "+tt.errContains &&
						err.Error() != "rows error: "+tt.errContains {
						t.Errorf("error should contain %q, got %v", tt.errContains, err)
					}
				}
				if names != nil {
					t.Errorf("expected nil names, got %v", names)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(names) != len(tt.wantNames) {
				t.Errorf("expected %d names, got %d", len(tt.wantNames), len(names))
			}

			for i, want := range tt.wantNames {
				if i < len(names) && names[i] != want {
					t.Errorf("names[%d] = %q, want %q", i, names[i], want)
				}
			}
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		wantNames   []string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - with duplicates",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Alice")
				m.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)
			},
			wantNames: []string{"Alice", "Bob", "Alice"},
			wantErr:   false,
		},
		{
			name: "success - single row",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("SingleUser")
				m.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)
			},
			wantNames: []string{"SingleUser"},
			wantErr:   false,
		},
		{
			name: "query error",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("^SELECT DISTINCT name FROM users$").
					WillReturnError(errConnection)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "connection failed",
		},
		{
			name: "scan error on null",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				m.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)
			},
			wantNames: nil,
			wantErr:   true,
		},
		{
			name: "rows error for unique",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					CloseError(errRows)
				m.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)
			},
			wantNames:   nil,
			wantErr:     true,
			errContains: "rows error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			tt.setupMock(mock)
			service := mydb.New(db)
			names, err := service.GetUniqueNames()

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errContains != "" && err != nil {
					if err.Error() != "db query: "+tt.errContains &&
						err.Error() != "rows error: "+tt.errContains {
						t.Errorf("error should contain %q, got %v", tt.errContains, err)
					}
				}
				if names != nil {
					t.Errorf("expected nil names, got %v", names)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(names) != len(tt.wantNames) {
				t.Errorf("expected %d names, got %d", len(tt.wantNames), len(names))
			}

			for i, want := range tt.wantNames {
				if i < len(names) && names[i] != want {
					t.Errorf("names[%d] = %q, want %q", i, names[i], want)
				}
			}
		})
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
