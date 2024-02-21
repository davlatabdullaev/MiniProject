package models

type Basket struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	TotalSum   int    `json:"total_sum"`
	Quantity   int    `json:"quantity"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  int    `json:"deleted_at"`
}

type CreateBasket struct {
	CustomerID string `json:"customer_id"`
	TotalSum   int    `json:"total_sum"`
	Quantity   int    `json:"quantity"`
}

type UpdateBasket struct {
	ID         string `json:"-"`
	CustomerID string `json:"customer_id"`
	TotalSum   int    `json:"total_sum"`
	Quantity   int    `json:"quantity"`
}

type BasketResponse struct {
	Baskets []Basket `json:"baskets"`
	Count   int      `json:"count"`
}
