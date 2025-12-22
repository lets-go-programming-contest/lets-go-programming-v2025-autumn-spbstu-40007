package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"grigorii.smolianinov/task-6/internal/db"
)

func TestDBService_GetNames_Coverage(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening mock db: %s", err)
	}
	defer dbConn.Close()

	service := db.New(dbConn)

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Scan Error",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(nil).
					RowError(0, errors.New("scan fail"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "Rows Iteration Error",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					CloseError(errors.New("close error"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			_, err := service.GetNames()
			assert.Error(t, err)
		})
	}
}

func TestDBService_GetUniqueNames_Coverage(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening mock db: %s", err)
	}
	defer dbConn.Close()

	service := db.New(dbConn)

	t.Run("Scan Error Unique", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		_, err := service.GetUniqueNames()
		assert.Error(t, err)
	})

	t.Run("Rows Err Unique", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Bob").CloseError(errors.New("err"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		_, err := service.GetUniqueNames()
		assert.Error(t, err)
	})

	t.Run("Query Error Unique", func(t *testing.T) {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errors.New("query fail"))
		_, err := service.GetUniqueNames()
		assert.Error(t, err)
	})
}
