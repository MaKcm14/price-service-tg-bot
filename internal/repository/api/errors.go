package api

import "errors"

var (
	ErrApiInteraction  = errors.New("error of sending the request to the price-service api")
	ErrBufferReading   = errors.New("error of reading the response's body")
	ErrJSONParser      = errors.New("error of parsing the JSON")
	ErrOfMarketsParams = errors.New("error of converting the markets to the URL-view: it's empty")
)
