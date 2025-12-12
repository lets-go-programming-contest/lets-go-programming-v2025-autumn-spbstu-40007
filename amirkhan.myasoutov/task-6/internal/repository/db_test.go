package repository_test

import (
	"context"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ami0-0/task-6/internal/repository"
)

func TestGetActiveUsers(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()

	svc := repository.NewDataService(db)
	
	rows := sqlmock.NewRows([]string{"username"}).AddRow("alice").AddRow("bob")
	mock.ExpectQuery("SELECT username FROM users").WillReturnRows(rows)

	res, err := svc.GetActiveUsers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []string{"alice", "bob"}, res)
}