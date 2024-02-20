package service

import (
	"context"
	"test/pkg/logger"
	"test/storage"
	"test/api/models"
)

type productService struct {
	storage storage.IStorage
	log     logger.ILogger
}

func NewProductService(storage storage.IStorage, log logger.ILogger) productService {
	return productService{
		storage: storage,
		log:     log,
	}
}

func (p productService) Create(ctx context.Context, product models.CreateProduct) (models.Product, error) {
	p.log.Info("product create service layer", logger.Any("product", product))

	id, err := p.storage.Product().Create(ctx, product)
	if err != nil {
		p.log.Error("error in service layer while creating product", logger.Error(err))
		return models.Product{}, err
	}

	createdProduct, err := p.storage.Product().GetByID(ctx, models.PrimaryKey{ID: id})
	if err != nil {
		p.log.Error("error in service layer while getting by id", logger.Error(err))
		return models.Product{}, err
	}

	return createdProduct, nil
}

func (p productService) Get(ctx context.Context, key models.PrimaryKey) (models.Product, error) {
	product, err := p.storage.Product().GetByID(ctx, key)
	if err != nil {
		p.log.Error("error in service layer while getting by id", logger.Error(err))
		return models.Product{}, err
	}

	return product, nil
}

func (p productService) GetList(ctx context.Context, request models.GetListRequest) (models.ProductsResponse, error) {
	p.log.Info("product get list service layer", logger.Any("product", request))

	products, err := p.storage.Product().GetList(ctx, request)
	if err != nil {
		p.log.Error("error in service layer while getting list", logger.Error(err))
		return models.ProductsResponse{}, err
	}

	return products, nil
}

func (p productService) Update(ctx context.Context, product models.UpdateProduct) (models.Product, error) {
	id, err := p.storage.Product().Update(ctx, product)
	if err != nil {
		p.log.Error("error in service layer while update", logger.Error(err))
		return models.Product{}, err
	}

	updatedProduct, err := p.storage.Product().GetByID(ctx, models.PrimaryKey{ID: id})
	if err != nil {
		p.log.Error("error in service layer while getting by id", logger.Error(err))
		return models.Product{}, err
	}

	return updatedProduct, nil
}

func (p productService) Delete(ctx context.Context, key models.PrimaryKey) error {
	err := p.storage.Product().Delete(ctx, key)

	return err
}
