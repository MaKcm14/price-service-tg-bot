package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

// PriceServiceApi defines the logic of price-service API-interaction
type PriceServiceApi struct {
	socket    string
	converter urlConverter
	logger    *slog.Logger
}

func NewPriceServiceApi(socket string, log *slog.Logger) PriceServiceApi {
	return PriceServiceApi{
		socket: socket,
		logger: log,
	}
}

// readResponse defines the logic of reading the response from the price-service.
func (p PriceServiceApi) readResponse(resp *http.Response) ([]byte, error) {
	const op = "api.read-response"

	body, err := ReadResponseBody(resp.Body, p.logger, op)

	if err != nil {
		return nil, err
	} else if resp.StatusCode > 299 {
		var errDesc string
		var res = make(map[string]string)

		json.Unmarshal(body, &errDesc)

		if desc, flagExist := res["error"]; flagExist {
			errDesc = desc
		}

		p.logger.Warn(fmt.Sprintf("error of the %s: %v: %v", op, ErrApiInteraction, errDesc))
		return nil, fmt.Errorf("error of the %s: %w: %s", op, ErrApiInteraction, errDesc)
	}

	return body, nil
}

// getProducts defines the main logic of getting the products for the different modes.
func (p PriceServiceApi) getProducts(url string, op string) (ProductResponse, error) {
	products := NewProductResponse()

	resp, err := http.Get(url)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %v: %v", op, err))
		return ProductResponse{}, fmt.Errorf("error of the %w: %v", ErrApiInteraction, err)
	}
	defer resp.Body.Close()

	body, err := p.readResponse(resp)

	if err != nil {
		return ProductResponse{}, err
	}

	err = json.Unmarshal(body, &products)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %v: %v", op, ErrJSONParser, err))
		return ProductResponse{}, fmt.Errorf("error of the %s: %w: %v", op, ErrJSONParser, err)
	}

	return products, nil
}

// sendRequest defines the logic of sending the POST-request for async products getting.
func (p PriceServiceApi) sendPostAsyncProdsRequest(url string, op string, body io.Reader) error {
	resp, err := http.Post(url, "application/json", body)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %v: %v", op, ErrApiInteraction, err))
		return fmt.Errorf("error of the %s: %w: %v", op, ErrApiInteraction, err)
	} else if resp.StatusCode > 299 {
		p.logger.Error(fmt.Sprintf("error of the %s: %v", op, ErrApiInteraction))
		return fmt.Errorf("error of the %s: %w", op, ErrApiInteraction)
	}

	return nil
}

// GetProductsByBestPrice defines the logic of getting the products by best price
// according to the user's request data.
func (p PriceServiceApi) GetProductsByBestPrice(request dto.ProductRequest) (map[string]entities.ProductSample, error) {
	const op = "api.best-price-getter"

	basePath := fmt.Sprintf("http://%s/products/filter/price/best-price", p.socket)

	markets, err := p.converter.convertMarkets(request.Markets)

	if err != nil {
		err = fmt.Errorf("error of the %s: error of markets params: %w", op, err)
		p.logger.Warn(err.Error())
		return nil, err
	}

	url := p.converter.buildURL(basePath,
		"query", request.Query, "markets", markets)

	products, err := p.getProducts(url, op)

	if err != nil {
		return nil, err
	}

	return products.Sample, nil
}

// SendAsyncBestPriceRequest sends async request for getting the products through the kafka.
func (p PriceServiceApi) SendAsyncBestPriceRequest(request dto.ProductRequest, headers map[string]string) error {
	const op = "api.best-price-async-sender"

	basePath := fmt.Sprintf("http://%s/products/filter/price/best-price/async", p.socket)

	markets, err := p.converter.convertMarkets(request.Markets)

	if err != nil {
		err = fmt.Errorf("error of the %s: error of markets params: %w", op, err)
		p.logger.Warn(err.Error())
		return err
	}

	url := p.converter.buildURL(basePath,
		"query", request.Query, "markets", markets)

	extraHeaders := newExtraHeaders(headers)

	res, _ := json.Marshal(extraHeaders)

	err = p.sendPostAsyncProdsRequest(url, op, bytes.NewReader(res))

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %v: %s", op, ErrApiInteraction, err))
		return fmt.Errorf("error of the %s: %w: %s", op, ErrApiInteraction, err)
	}

	return nil
}

// GetSupportedMarkets gets the markets that are supported by the price-service.
func (p PriceServiceApi) GetSupportedMarkets() (entities.SupportedMarkets, error) {
	const op = "api.get-supported-markets"

	resp, err := http.Get(fmt.Sprintf("http://%s/api/markets", p.socket))

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return entities.SupportedMarkets{}, fmt.Errorf("error of the %s: %w: %s", op, ErrApiInteraction, err)
	} else if resp.StatusCode > 299 {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, ErrApiInteraction))
		return entities.SupportedMarkets{}, fmt.Errorf("error of the %s: %w", op, ErrApiInteraction)
	}
	defer resp.Body.Close()

	var markets entities.SupportedMarkets

	body, err := ReadResponseBody(resp.Body, p.logger, op)

	if err != nil {
		return entities.SupportedMarkets{}, err
	}

	json.Unmarshal(body, &markets)

	return markets, nil
}
