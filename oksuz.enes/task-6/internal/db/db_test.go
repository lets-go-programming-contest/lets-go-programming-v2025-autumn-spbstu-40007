package db

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetWifiStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	defer db.Close()

	repo := &Repository{DB: db}

	rows := sqlmock.NewRows([]string{"status"}).AddRow("active")
	mock.ExpectQuery("SELECT status FROM wifi_table WHERE id = ?").WithArgs(1).WillReturnRows(rows)

	status, err := repo.GetWifiStatus(1)
	assert.NoError(t, err)
	assert.Equal(t, "active", status)

	mock.ExpectQuery("SELECT status FROM wifi_table WHERE id = ?").WithArgs(2).WillReturnError(sql.ErrNoRows)
	_, err = repo.GetWifiStatus(2)
	assert.Error(t, err)
}
