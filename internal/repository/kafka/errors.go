package kafka

import "errors"

var (
	ErrKafkaConnection = errors.New("error of connecting to the Kakfa's cluster")
)
