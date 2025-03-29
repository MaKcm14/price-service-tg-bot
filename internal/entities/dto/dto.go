package dto

// ProductRequest defines the user's product request.
type ProductRequest struct {
	Markets map[string]string
	Query   string
}

func NewProductRequest() ProductRequest {
	return ProductRequest{
		Markets: make(map[string]string),
	}
}
