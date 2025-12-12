package db_test

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ami0-0/task-6/internal/db"
)

func TestGetNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = mockDB.Close() }()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()
	assert.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob"}, names)
}