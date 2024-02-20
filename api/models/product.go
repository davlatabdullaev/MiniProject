package models

type Product struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"product"`
	Quantity  int    `json:"quantity"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt int    `json:"deleted_at"`
}

type CreateProduct struct {
	Name     string `json:"name"`
	Price    int    `json:"product"`
	Quantity int    `json:"quantity"`
}

type UpdateProduct struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	Price    int    `json:"product"`
	Quantity int    `json:"quantity"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}
