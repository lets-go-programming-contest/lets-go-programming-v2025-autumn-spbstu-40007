package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type SQLHandler interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type DataService struct {
	executor SQLHandler
}

func NewDataService(db SQLHandler) DataService {
	return DataService{executor: db}
}

func (s DataService) GetActiveUsers(ctx context.Context) ([]string, error) {
	const query = "SELECT username FROM users WHERE status = 'active'"
	rows, err := s.executor.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}