package postgres

import "errors"

var (
	ErrDBConnection    = errors.New("error of setting the connection to the DB")
	ErrQueryExec       = errors.New("error while query execution")
	ErrCacheConnection = errors.New("error of setting the connection with the cache")
)
