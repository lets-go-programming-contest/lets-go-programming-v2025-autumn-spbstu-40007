package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/se1lzor/task-6/internal/db" // <-- замени на свой module path
)

var (
	errDB  = errors.New("db error")
	errRow = errors.New("row iteration error")
)

func TestService_GetNames_Table(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMock   func(m sqlmock.Sqlmock)
		want        []string
		wantErr     bool
		errIs       error
		errContains string
	}{
		{
			name: "success",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name: "query error",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT name FROM users").WillReturnError(errDB)
			},
			wantErr:     true,
			errIs:       errDB,
			errContains: "db query",
		},
		{
			name: "scan error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "rows scanning",
		},
		{
			name: "rows error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRow)
				m.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errIs:       errRow,
			errContains: "rows error",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			conn, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { _ = conn.Close() })

			svc := db.New(conn)
			tc.setupMock(mock)

			got, err := svc.GetNames()
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				if tc.errIs != nil {
					require.ErrorIs(t, err, tc.errIs)
				}
				if tc.errContains != "" {
					require.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_GetUniqueNames_Table(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(m sqlmock.Sqlmock)
		want      []string
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name: "query error",
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errDB)
			},
			wantErr: true,
		},
		{
			name: "scan error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "rows error",
			setupMock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errRow)
				m.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			conn, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { _ = conn.Close() })

			svc := db.New(conn)
			tc.setupMock(mock)

			got, err := svc.GetUniqueNames()
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
