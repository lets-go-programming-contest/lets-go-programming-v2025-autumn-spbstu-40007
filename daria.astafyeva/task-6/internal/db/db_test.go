package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/itsdasha/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errQuery = errors.New("query failed")
	errScan  = errors.New("scan failed")
	errRows  = errors.New("rows iteration failed")
)

func TestDBService_New(t *testing.T) {
	t.Parallel()

	dbMock, _, _ := sqlmock.New()
	defer dbMock.Close()

	service := db.New(dbMock)

	assert.NotNil(t, service)
	assert.Same(t, dbMock, service.DB)
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mock sqlmock.Sqlmock)
		wantErr    bool
		want       []string
		errContain string
	}{
		{
			name: "success - multiple names",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Dasha").
					AddRow("Alex")
				mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)
			},
			want: []string{"Dasha", "Alex"},
		},
		{
			name: "success - empty result",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "query error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT name FROM users$").WillReturnError(errQuery)
			},
			wantErr:    true,
			errContain: "db query",
		},
		{
			name: "scan error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)
			},
			wantErr:    true,
			errContain: "rows scanning",
		},
		{
			name: "rows error after Next",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Dasha").
					RowError(0, errRows)
				mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)
			},
			wantErr:    true,
			errContain: "rows error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			tt.setup(mock)

			service := db.New(dbMock)
			names, err := service.GetNames()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContain)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				if tt.want == nil {
					assert.Nil(t, names)
				} else {
					assert.Equal(t, tt.want, names)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mock sqlmock.Sqlmock)
		wantErr    bool
		want       []string
		errContain string
	}{
		{
			name: "success",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Dasha").
					AddRow("Alex")
				mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)
			},
			want: []string{"Dasha", "Alex"},
		},
		{
			name: "query error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnError(errQuery)
			},
			wantErr:    true,
			errContain: "db query",
		},
		{
			name: "scan error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)
			},
			wantErr:    true,
			errContain: "rows scanning",
		},
		{
			name: "rows error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Dasha").
					RowError(0, errRows)
				mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)
			},
			wantErr:    true,
			errContain: "rows error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dbMock, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer dbMock.Close()

			tt.setup(mock)

			service := db.New(dbMock)
			names, err := service.GetUniqueNames()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContain)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, names)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
