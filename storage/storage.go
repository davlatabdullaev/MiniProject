package storage

import (
	"context"
	"test/api/models"
)

type IStorage interface {
	Close()
	Basket() IBasketStorage
	Product() IProductStorage
	Customer() ICustomerStorage
}

type ICustomerStorage interface {
	Create(context.Context, models.CreateCustomer) (string, error)
	GetByID(context.Context, models.PrimaryKey) (models.Customer, error)
	GetList(context.Context, models.GetListRequest) (models.CustomersResponse, error)
	Update(context.Context, models.UpdateCustomer) (string, error)
	Delete(context.Context, models.PrimaryKey) error
	GetPassword(context.Context, string) (string, error)
	UpdatePassword(context.Context, models.UpdateCustomerPassword) error
	GetCustomerCredentialsByLogin(context.Context, string) (models.Customer, error)
}

type IProductStorage interface {
	Create(context.Context, models.CreateProduct) (string, error)
	GetByID(context.Context, models.PrimaryKey) (models.Product, error)
	GetList(context.Context, models.GetListRequest) (models.ProductsResponse, error)
	Update(context.Context, models.UpdateProduct) (string, error)
	Delete(context.Context, models.PrimaryKey) error
}

type IBasketStorage interface {
	Create(context.Context, models.CreateBasket) (string, error)
	GetByID(context.Context, models.PrimaryKey) (models.Basket, error)
	GetList(context.Context, models.GetListRequest) (models.BasketResponse, error)
	Update(context.Context, models.UpdateBasket) (string, error)
	Delete(context.Context, models.PrimaryKey) error
}
