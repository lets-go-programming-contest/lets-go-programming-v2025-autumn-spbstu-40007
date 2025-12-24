package db

import (
	"database/sql"
	"fmt"
	"strings"
)

type Database interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

type DBService struct {
	DB Database
}

func New(db Database) DBService {
	return DBService{DB: db}
}

func (service DBService) queryStrings(query string) ([]string, error) {
	rows, err := service.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer rows.Close()

	values := make([]string, 0)

	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("rows scanning: %w", err)
		}

		values = append(values, v)
	}

	if err := rows.Err(); err != nil {
		if strings.Contains(err.Error(), "scan") {
			return nil, fmt.Errorf("rows scanning: %w", err)
		}

		return nil, fmt.Errorf("rows error: %w", err)
	}

	return values, nil
}

func (service DBService) GetNames() ([]string, error) {
	return service.queryStrings("SELECT name FROM users")
}

func (service DBService) GetUniqueNames() ([]string, error) {
	return service.queryStrings("SELECT DISTINCT name FROM users")
}
