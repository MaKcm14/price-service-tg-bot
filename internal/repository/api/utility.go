package api

import (
	"fmt"
	"io"
	"log/slog"
)

func ReadResponseBody(source io.Reader, logger *slog.Logger, serviceType string) ([]byte, error) {
	respBody := make([]byte, 0, 100000)

	for {
		buffer := make([]byte, 10000)
		n, err := source.Read(buffer)

		if n != 0 && (err == nil || err == io.EOF) {
			respBody = append(respBody, buffer[:n]...)
		} else if err != nil && err != io.EOF {
			logger.Warn(fmt.Sprintf("error of the %v: %v: %v", serviceType, ErrBufferReading, err))
			return nil, fmt.Errorf("%w: %v", ErrBufferReading, err)
		} else if err == io.EOF {
			break
		}
	}

	return respBody, nil
}
