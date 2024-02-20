package models

type Basket struct {
	ID         string `json:"id"`
	ProductID  string `json:"product_id"`
	CustomerID string `json:"customer_id"`
	TotalSum   int    `json:"total_sum"`
	Quantity   int    `json:"quantity"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  int    `json:"deleted_at"`
}

type CreateBasket struct {
	ProductID  string `json:"product_id"`
	CustomerID string `json:"customer_id"`
	TotalSum   int    `json:"total_sum"`
	Quantity   int    `json:"quantity"`
}

type UpdateBasket struct {
	ID         string `json:"-"`
	ProductID  string `json:"product_id"`
	CustomerID string `json:"customer_id"`
	TotalSum   int    `json:"total_sum"`
	Quantity   int    `json:"quantity"`
}

type BasketResponse struct {
	Baskets []Basket `json:"baskets"`
	Count   int      `json:"count"`
}
