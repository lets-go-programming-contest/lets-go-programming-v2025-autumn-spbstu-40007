package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ami0-0/task-6/internal/db"
)

func TestGetActiveUsers(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = sqlDB.Close() }()

	svc := db.NewDataService(sqlDB)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"username"}).
			AddRow("admin").
			AddRow("user1")

		mock.ExpectQuery("SELECT username FROM users").WillReturnRows(rows)

		res, err := svc.GetActiveUsers(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT username FROM users").
			WillReturnError(errors.New("db fail"))

		_, err := svc.GetActiveUsers(ctx)
		assert.Error(t, err)
	})
}