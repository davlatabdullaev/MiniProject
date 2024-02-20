package postgres

import (
	"context"
	"test/api/models"
	"test/config"
	"test/pkg/helper"
	"test/pkg/logger"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestProductRepo_Create(t *testing.T) {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)
	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Errorf("error while connection to db error: %v", err)
	}

	createproducts := models.CreateProduct{
		Name:     "product 1",
		Price:    12,
		Quantity: 2,
	}
	id, err := pgStore.Product().Create(context.Background(), createproducts)
	if err != nil {
		t.Error("error while inserting product", err)

	}
	idproduct, err := pgStore.Product().GetByID(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		t.Error("error", err)
	}
	if id == "" {
		t.Error("error while creating product")
	}
	assert.Equal(t, idproduct.Name, createproducts.Name)
	assert.Equal(t, idproduct.Quantity, createproducts.Quantity)
	assert.Equal(t, idproduct.Price, createproducts.Price)

}

func TestProductRepo_GetByID(t *testing.T) {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)
	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Errorf("error while connection to db error: %v", err)
	}
	products, err := pgStore.Product().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  1,
		Search: "",
	})

	if len(products.Products) == 0 {
		t.Error("error", err)

	}

	expectedproducts := products.Products[0].ID
	t.Run("succes", func(t *testing.T) {
		product, err := pgStore.Product().GetByID(context.Background(), models.PrimaryKey{ID: expectedproducts})

		if err != nil {
			t.Error("error while geting by id product", err)
		}
		if product.ID != expectedproducts {
			t.Errorf("expected: %q but got: %q", expectedproducts, product.ID)
		}
		if product.Name == "" {
			t.Error("expected: productname but got : nothing")
		}
		if product.Quantity <= 0 {
			t.Errorf("exepcted: more than 0 ,but got: %q", product.Quantity)
		}
		if product.Price <= 0 {

			t.Errorf("expeceted: more than 0 price but got %q", product.Price)
		}

	})

	t.Run("fail", func(t *testing.T) {
		productid := ""
		product, err := pgStore.Product().GetByID(context.Background(), models.PrimaryKey{
			ID: productid,
		})
		if err != nil {
			t.Error("error while getting product id", err)
		}
		if product.ID != productid {
			t.Errorf("expected: %q, but got %q", productid, product.ID)
		}
		if product.Name == "" {
			t.Error("expected: productname but got : nothing")
		}
		if product.Quantity <= 0 {
			t.Errorf("exepcted: more than 0 ,but got: %q", product.Quantity)
		}
		if product.Price <= 0 {

			t.Errorf("expeceted: more than 0 price but got %q", product.Price)
		}

	})

}

func TestProductRepo_GetList(t *testing.T) {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)
	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Errorf("error while connecting db %q", err)
	}

	products, err := pgStore.Product().GetList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 1000,
	})
	if err != nil {
		t.Error("error while getting list of products", err.Error())
	}
	if len(products.Products) != 5 {
		t.Errorf("expected 5 rows , but got %q", len(products.Products))
	}

	assert.Equal(t, len(products.Products), 5)

}

func TestProductRepo_Update(t *testing.T) {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)
	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Error("error while connecting to db ", err)
	}

	createProduct := models.CreateProduct{
		Name:     helper.GenerateProductName(),
		Price:    int(helper.GenerateRandomPrice(10.0, 100.0)),
		Quantity: 10,
	}

	productID, err := pgStore.Product().Create(context.Background(), createProduct)
	if err != nil {
		t.Error("erro while creating product in tetsing", err)
	}

	if err != nil {
		t.Errorf("error while creating product %v", err)
	}

	updateProduct := models.UpdateProduct{
		ID:       productID,
		Name:     helper.GenerateProductName(),
		Price:    int(helper.GenerateRandomPrice(10.0, 100.0)),
		Quantity: 10,
	}

	productupdateid, err := pgStore.Product().Update(context.Background(), updateProduct)
	if err != nil {
		t.Error("error updatinf product in testing", err)
	}

	product, err := pgStore.Product().GetByID(context.Background(), models.PrimaryKey{
		ID: productupdateid,
	})
	if err != nil {
		t.Error("error while geting by id in testing product", err)
	}
	if productupdateid == "" {
		t.Error("expected updated product id: but got empty string")
	}

	assert.Equal(t, product.ID, productupdateid)
	assert.Equal(t, product.Name, updateProduct.Name)
	assert.Equal(t, product.Price, updateProduct.Price)
	assert.Equal(t, product.Quantity, updateProduct.Quantity)

}

func TestProductRepo_Delete(t *testing.T) {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)
	pgStore, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Error("error while connecting to db ", err)
	}

	createproduct := models.CreateProduct{
		Name:     helper.GenerateProductName(),
		Price:    int(helper.GenerateRandomPrice(10.0, 100.0)),
		Quantity: 10,
	}

	productid, err := pgStore.Product().Create(context.Background(), createproduct)
	if err != nil {
		t.Error("error while creating product in delete tetsing", err)
	}

	if err = pgStore.Product().Delete(context.Background(), models.PrimaryKey{ID: productid}); err != nil {
		t.Error("error while deleteing product in testing", err)
	}
	t.Run("falied", func(t *testing.T) {
		if productid == "" {
			t.Error("expected product id but go nothing", err)
		}
	})

}
