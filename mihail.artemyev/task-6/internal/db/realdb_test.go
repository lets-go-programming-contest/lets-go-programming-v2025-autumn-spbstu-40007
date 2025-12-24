package db_test

import (
	"database/sql"
	"fmt"
)

type realDB struct {
	db *sql.DB
}

func (r *realDB) Query(query string, args ...any) (*sql.Rows, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	return rows, nil
}
