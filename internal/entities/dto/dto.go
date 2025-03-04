package dto

import "github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities"

// ProductRequest defines the user's product request.
type ProductRequest struct {
	Markets map[string]string
	Query   string

	Mode entities.Mode
}

func NewProductRequest(mode entities.Mode) ProductRequest {
	return ProductRequest{
		Markets: make(map[string]string),
		Mode:    mode,
	}
}
