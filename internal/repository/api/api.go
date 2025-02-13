package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

// PriceServiceApi defines the logic of price-service API-interaction
type PriceServiceApi struct {
	socket string
	logger *slog.Logger
}

func NewPriceServiceApi(socket string, log *slog.Logger) PriceServiceApi {
	return PriceServiceApi{
		socket: socket,
		logger: log,
	}
}

// GetProductsByBestPrice defines the logic of getting the products by best price
// according to the user's request data.
func (p PriceServiceApi) GetProductsByBestPrice(request dto.ProductRequest) (map[string]entities.ProductSample, error) {
	const op = "api.best-price-getter"

	products := struct {
		Sample map[string]entities.ProductSample `json:"samples"`
	}{
		make(map[string]entities.ProductSample, 1000),
	}

	url := fmt.Sprintf("http://%s/products/filter/price/best-price?query=%s&markets=", p.socket, url.QueryEscape(request.Query))

	count := 0
	for _, market := range request.Markets {
		url += market
		count++
		if count != len(request.Markets) {
			url += "+"
		}
	}

	resp, err := http.Get(url)

	if err != nil {
		p.logger.Warn(fmt.Sprintf("error of the %v: %v", op, err))
		return nil, fmt.Errorf("error of the %w: %v", ErrApiInteraction, err)
	}
	defer resp.Body.Close()

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

	err = json.Unmarshal(body, &products)

	if err != nil {
		p.logger.Warn(fmt.Sprintf("error of %s: %v: %v", op, ErrJSONParser, err))
		return nil, fmt.Errorf("error of %s: %w: %v", op, ErrJSONParser, err)
	}

	return products.Sample, nil
}
