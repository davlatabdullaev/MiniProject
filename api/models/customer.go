package models

type Customer struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt int    `json:"deleted_at"`
}

type CreateCustomer struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UpdateCustomer struct {
	ID       string `json:"-"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type CustomersResponse struct {
	Customers []Customer `json:"customers"`
	Count     int        `json:"count"`
}
