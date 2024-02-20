package service

import (
	"test/pkg/logger"
	"test/storage"
)

type IServiceManager interface {
	Customer() customerService
	Basket() basketService
	Product() productService
	AuthService() authService
}

type Service struct {
	customerService customerService
	basketService   basketService
	productService  productService
	authService     authService
}

func New(storage storage.IStorage, log logger.ILogger) Service {
	services := Service{}

	services.customerService = NewCustomerService(storage, log)
	services.basketService = NewBasketService(storage, log)
	services.productService = NewProductService(storage, log)
	services.authService = NewAuthService(storage, log)

	return services
}

func (s Service) Customer() customerService {
	return s.customerService
}

func (s Service) Basket() basketService {
	return s.basketService
}

func (s Service) Product() productService {
	return s.productService
}

func (s Service) AuthService() authService {
	return s.authService
}
