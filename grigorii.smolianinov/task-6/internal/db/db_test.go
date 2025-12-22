package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"grigorii.smolianinov/task-6/internal/db"
)

func TestDBService_AllScenarios(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening mock: %s", err)
	}
	defer dbConn.Close()

	service := db.New(dbConn)

	tests := []struct {
		name    string
		method  string
		query   string
		mock    func()
		wantErr bool
	}{
		{
			name:   "GetNames Success",
			method: "GetNames",
			query:  "SELECT name FROM users",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:   "GetNames Query Error",
			method: "GetNames",
			query:  "SELECT name FROM users",
			mock: func() {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errors.New("query fail"))
			},
			wantErr: true,
		},
		{
			name:   "GetNames Scan Error",
			method: "GetNames",
			query:  "SELECT name FROM users",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:   "GetNames Rows Err",
			method: "GetNames",
			query:  "SELECT name FROM users",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errors.New("row error"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:   "GetUniqueNames Success",
			method: "GetUniqueNames",
			query:  "SELECT DISTINCT name FROM users",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:   "GetUniqueNames Query Error",
			method: "GetUniqueNames",
			query:  "SELECT DISTINCT name FROM users",
			mock: func() {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("query fail"))
			},
			wantErr: true,
		},
		{
			name:   "GetUniqueNames Scan Error",
			method: "GetUniqueNames",
			query:  "SELECT DISTINCT name FROM users",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:   "GetUniqueNames Rows Err",
			method: "GetUniqueNames",
			query:  "SELECT DISTINCT name FROM users",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errors.New("row error"))
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // [cite: 355]
			tt.mock()
			var err error
			if tt.method == "GetNames" {
				_, err = service.GetNames()
			} else {
				_, err = service.GetUniqueNames()
			}

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
