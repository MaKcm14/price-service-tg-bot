package api

import (
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"

	"github.com/MaKcm14/price-service/pkg/entities"
)

// urlConverter defines the logic of building the URL.
type urlConverter struct{}

// convertMarkets converts the users' markets to the URL-string.
func (u urlConverter) convertMarkets(markets map[string]string) (string, error) {
	if len(markets) == 0 {
		return "", ErrOfMarketsParams
	}
	var marketsView string

	count := 0
	for market := range markets {
		marketsView += strings.ToLower(market)
		count++
		if count != len(markets) {
			marketsView += "+"
		}
	}

	return marketsView, nil
}

// buildURL defines the logic of building the URL using the current filters.
func (u urlConverter) buildURL(basePath string, filters ...string) string {
	if len(filters) == 0 {
		return basePath
	}

	basePath += "?"
	for i := 0; i != len(filters); i += 2 {
		basePath += fmt.Sprintf("%s=%s", url.QueryEscape(filters[i]), url.QueryEscape(filters[i+1]))

		if i+2 < len(filters) {
			basePath += "&"
		}
	}

	return basePath
}

// productResponse defines the user's product response from API.
type ProductResponse struct {
	Sample map[string]entities.ProductSample `json:"samples"`
}

func NewProductResponse() ProductResponse {
	return ProductResponse{
		Sample: make(map[string]entities.ProductSample),
	}
}

// ReadResponseBody defines the logic of reading the response from the source.
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

// Header defines the data structure of the extra-header.
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ExtraHeader defines the data structure of extra-headers that must be included in the async response.
type ExtraHeaders struct {
	Headers []Header `json:"headers"`
}

func NewExtraHeaders(headers map[string]string) ExtraHeaders {
	extraHeaders := ExtraHeaders{
		make([]Header, 0, len(headers)),
	}

	for key, val := range headers {
		extraHeaders.Headers = append(extraHeaders.Headers, Header{
			Key:   key,
			Value: val,
		})
	}

	return extraHeaders
}
