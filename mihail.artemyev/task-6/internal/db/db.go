package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNilRows = errors.New("db returned nil rows without error")

type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type DBService struct {
	DB Database
}

func New(db Database) DBService {
	return DBService{DB: db}
}

func (s DBService) queryStrings(query string) ([]string, error) {
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	if rows == nil {
		return nil, ErrNilRows
	}
	defer rows.Close()

	var res []string

	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("rows scanning: %w", err)
		}
		res = append(res, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return res, nil
}

func (s DBService) GetNames() ([]string, error) {
	return s.queryStrings("SELECT name FROM users")
}

func (s DBService) GetUniqueNames() ([]string, error) {
	return s.queryStrings("SELECT DISTINCT name FROM users")
}
