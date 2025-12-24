package db_test

import "errors"

var (
	errScanError = errors.New("scan error")
	errRowsError = errors.New("rows error")
	errDBError   = errors.New("db error")
)
