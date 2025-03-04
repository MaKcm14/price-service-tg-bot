package redis

import "errors"

var (
	ErrConnToRedis    = errors.New("error of connection to the redis")
	ErrOfRedisRequest = errors.New("error of the sending and executing the request by redis-server")
)
