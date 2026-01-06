package db

import (
	"database/sql"
	"fmt"
)

type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type Service struct {
	db Database
}

func New(db Database) Service {
	return Service{db: db}
}

func (s Service) GetNames() ([]string, error) {
	return s.fetchNames("SELECT name FROM users")
}

func (s Service) GetUniqueNames() ([]string, error) {
	return s.fetchNames("SELECT DISTINCT name FROM users")
}

func (s Service) fetchNames(query string) ([]string, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer rows.Close()

	out := make([]string, 0)

	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("rows scanning: %w", err)
		}

		out = append(out, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return out, nil
}
