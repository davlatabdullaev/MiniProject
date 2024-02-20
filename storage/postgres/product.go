package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"test/api/models"
	"test/pkg/logger"
	"test/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type productRepo struct {
	db    *pgxpool.Pool
	log   logger.ILogger
}

func NewProductRepo(db *pgxpool.Pool, log logger.ILogger) storage.IProductStorage {
	return &productRepo{
		db:    db,
		log:   log,
	}
}

func (p *productRepo) Create(ctx context.Context, product models.CreateProduct) (string, error) {
	id := uuid.New()
	query := `insert into products(id, name, price, original_price, quantity, category_id, branch_id) 
						values($1, $2, $3, $4, $5, $6, $7)`

	if rowsAffected, err := p.db.Exec(ctx, query,
		id,
		product.Name,
		product.Price,
		product.Quantity,
		); err != nil {
		if r := rowsAffected.RowsAffected(); r == 0 {
			p.log.Error("error is in rows affected", logger.Error(err))
			return "", err
		}
		p.log.Error("error is while inserting product", logger.Error(err))
		return "", err
	}

	return id.String(), nil
}

func (p *productRepo) GetByID(ctx context.Context, key models.PrimaryKey) (models.Product, error) {
	var createdAt, updatedAt = sql.NullString{}, sql.NullString{}
	product := models.Product{}
	query := `select id, name, price, original_price, quantity, category_id, branch_id, created_at, updated_at
							from products where id = $1 and deleted_at = 0`
	if err := p.db.QueryRow(ctx, query, key.ID).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Quantity,
		&createdAt,
		&updatedAt); err != nil {
		p.log.Error("error is while selecting product by id", logger.Error(err))
		return models.Product{}, err
	}

	if createdAt.Valid {
		product.CreatedAt = createdAt.String
	}

	if updatedAt.Valid {
		product.UpdatedAt = updatedAt.String
	}
	return product, nil
}

func (p *productRepo) GetList(ctx context.Context, request models.GetListRequest) (models.ProductsResponse, error) {
	var (
		products             = []models.Product{}
		page                 = request.Page
		offset               = (page - 1) * request.Limit
		search               = request.Search
		query, countQuery    string
		count                = 0
		createdAt, updatedAt = sql.NullString{}, sql.NullString{}
	)

	countQuery = `select count(1) from products where deleted_at = 0 `

	if search != "" {
		countQuery += fmt.Sprintf(` and (name ilike '%%%s%%' or 
			CAST(price AS TEXT) ilike '%s' or CAST(quantity AS TEXT) ilike '%s')`, search, search, search)
	}

	if err := p.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		p.log.Error("error is while scanning count", logger.Error(err))
		return models.ProductsResponse{}, err
	}

	query = `select id, name, price, original_price, quantity, category_id, branch_id, created_at, updated_at
								from products where deleted_at = 0`

	if search != "" {
		query += fmt.Sprintf(` and (name ilike '%%%s%%' or 
			CAST(price AS TEXT) ilike '%s' or CAST(quantity AS TEXT) ilike '%s')`, search, search, search)
	}

	query += ` order by created_at desc LIMIT $1 OFFSET $2`

	rows, err := p.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		p.log.Error("error is while selecting products", logger.Error(err))
		return models.ProductsResponse{}, err
	}

	for rows.Next() {
		product := models.Product{}
		if err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
			&createdAt,
			&updatedAt); err != nil {
			p.log.Error("error is while scanning products", logger.Error(err))
			return models.ProductsResponse{}, err
		}
		if createdAt.Valid {
			product.CreatedAt = createdAt.String
		}

		if updatedAt.Valid {
			product.UpdatedAt = updatedAt.String
		}
		products = append(products, product)
	}
	return models.ProductsResponse{
		Products: products,
		Count:    count,
	}, err
}

func (p *productRepo) Update(ctx context.Context, product models.UpdateProduct) (string, error) {
	query := `update products set name = $1, price = $2, original_price = $3, quantity = $4, 
                    category_id = $5, updated_at = now()  where id = $6`

	if _, err := p.db.Exec(ctx, query,
		&product.Name,
		&product.Price,
		&product.Quantity,
		&product.ID); err != nil {
		p.log.Error("error is while updating product", logger.Error(err))
		return "", err
	}

	return product.ID, nil
}

func (p *productRepo) Delete(ctx context.Context, key models.PrimaryKey) error {
	query := `update products set deleted_at = extract(epoch from current_timestamp) where id = $1`

	if rowsAffected, err := p.db.Exec(ctx, query, key.ID); err != nil {
		if r := rowsAffected.RowsAffected(); r == 0 {
			p.log.Error("error is in rows affected", logger.Error(err))
			return err
		}
		p.log.Error("error is while deleting product", logger.Error(err))
		return err
	}
	return nil
}
