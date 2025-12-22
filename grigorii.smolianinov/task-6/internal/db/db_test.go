package db_test

import (
	"errors"
	"testing"

	"grigorii.smolianinov/task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDBService_GetNames(t *testing.T) {
	// Инициализация mock [cite: 485]
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening mock db: %s", err)
	}
	defer dbConn.Close()

	service := db.New(dbConn)

	// Table-driven test
	tests := []struct {
		name    string
		mock    func()
		want    []string
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			want:    []string{"Alice", "Bob"},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(errors.New("db fail"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := service.GetNames()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening mock db: %s", err)
	}
	defer dbConn.Close()

	service := db.New(dbConn)

	t.Run("Success Unique", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).AddRow("Admin")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		got, err := service.GetUniqueNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"Admin"}, got)
	})
}
