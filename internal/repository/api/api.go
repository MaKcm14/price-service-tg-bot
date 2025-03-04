package api

import (
	"encoding/json"
	"fmt"
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

		p.logger.Warn(fmt.Sprintf("error of %s: %v: %v", op, ErrApiInteraction, errDesc))
		return nil, fmt.Errorf("error of the %s: %w: %s", op, ErrApiInteraction, errDesc)
	}

	return body, nil
}

// getProducts defines the main logic of getting the products for the different modes.
func (p PriceServiceApi) getProducts(url string, op string) (productResponse, error) {
	products := newProductResponse()

	resp, err := http.Get(url)

	if err != nil {
		p.logger.Warn(fmt.Sprintf("error of the %v: %v", op, err))
		return productResponse{}, fmt.Errorf("error of the %w: %v", ErrApiInteraction, err)
	}
	defer resp.Body.Close()

	body, err := p.readResponse(resp)

	if err != nil {
		return productResponse{}, err
	}

	err = json.Unmarshal(body, &products)

	if err != nil {
		p.logger.Warn(fmt.Sprintf("error of %s: %v: %v", op, ErrJSONParser, err))
		return productResponse{}, fmt.Errorf("error of %s: %w: %v", op, ErrJSONParser, err)
	}

	return products, nil
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

	//DEBUG:
	fmt.Println(url)
	//TODO: delete

	products, err := p.getProducts(url, op)

	if err != nil {
		return nil, err
	}

	return products.Sample, nil
}
