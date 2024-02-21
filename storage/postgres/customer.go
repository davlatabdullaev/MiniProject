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
)

type customerRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewCustomerRepo(db *pgxpool.Pool, log logger.ILogger) storage.ICustomerStorage {
	return &customerRepo{
		db:  db,
		log: log,
	}
}

func (c *customerRepo) Create(ctx context.Context, createUser models.CreateCustomer) (string, error) {

	uid := uuid.New()

	if _, err := c.db.Exec(ctx, `insert into 
			customers values ($1, $2, $3, $4)
			`,
		uid,
		createUser.FullName,
		createUser.Phone,
		createUser.Password,
	); err != nil {
		c.log.Error("error while inserting data", logger.Error(err))
		return "", err
	}

	return uid.String(), nil
}

func (c *customerRepo) GetByID(ctx context.Context, pKey models.PrimaryKey) (models.Customer, error) {
	var createdAt, updatedAt = sql.NullTime{}, sql.NullString{}
	user := models.Customer{}

	query := `
		select id, full_name, phone, created_at, updated_at 
						from customers where id = $1 and deleted_at = 0 
`
	if err := c.db.QueryRow(ctx, query, pKey.ID).Scan(
		&user.ID,
		&user.FullName,
		&user.Phone,
		&createdAt,
		&updatedAt,
	); err != nil {
		c.log.Error("error while scanning user", logger.Error(err))
		return models.Customer{}, err
	}

	if createdAt.Valid {
		user.CreatedAt = createdAt.Time.String()
	}

	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.String
	}

	return user, nil
}

func (c *customerRepo) GetList(ctx context.Context, request models.GetListRequest) (models.CustomersResponse, error) {
	var (
		customers            = []models.Customer{}
		count                = 0
		countQuery, query    string
		page                 = request.Page
		offset               = (page - 1) * request.Limit
		search               = request.Search
		createdAt, updatedAt = sql.NullTime{}, sql.NullString{}
	)

	countQuery = `
		SELECT count(1) from customers where deleted_at = 0 `

	if search != "" {
		countQuery += fmt.Sprintf(` and (phone ilike '%s' or full_name ilike '%s')`, search, search)
	}

	if err := c.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of users", err.Error())
		return models.CustomersResponse{}, err
	}

	query = `
		SELECT id, full_name, phone, created_at, updated_at
			FROM customers
			    WHERE deleted_at = 0
			    `

	if search != "" {
		query += fmt.Sprintf(` and (phone ilike '%s' or full_name ilike '%s') `, search, search)
	}

	query += ` order by created_at desc LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.CustomersResponse{}, err
	}

	for rows.Next() {
		user := models.Customer{}

		if err = rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Phone,
			&createdAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.CustomersResponse{}, err
		}
		if createdAt.Valid {
			user.CreatedAt = createdAt.Time.String()
		}

		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.String
		}
		customers = append(customers, user)
	}

	return models.CustomersResponse{
		Customers: customers,
		Count:     count,
	}, nil
}

func (c *customerRepo) Update(ctx context.Context, request models.UpdateCustomer) (string, error) {
	query := `
		update customers 
			set full_name = $1, phone = $2, updated_at = now()
				where id = $4`

	if _, err := c.db.Exec(ctx, query, request.FullName, request.Phone, request.ID); err != nil {
		fmt.Println("error while updating user data", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (c *customerRepo) Delete(ctx context.Context, request models.PrimaryKey) error {
	query := `update customers set deleted_at = extract(epoch from current_timestamp) where id = $1`

	if _, err := c.db.Exec(ctx, query, request.ID); err != nil {
		fmt.Println("error while deleting customer by id", err.Error())
		return err
	}

	return nil
}

func (c *customerRepo) GetPassword(ctx context.Context, id string) (string, error) {
	password := ""

	query := `
		select password from customers 
		                where id = $1`

	if err := c.db.QueryRow(ctx, query, id).Scan(&password); err != nil {
		fmt.Println("Error while scanning password from customers", err.Error())
		return "", err
	}

	return password, nil
}

func (c *customerRepo) UpdatePassword(ctx context.Context, request models.UpdateCustomerPassword) error {
	query := `
		update customers 
				set password = $1, updated_at = now()
					where id = $2 `

	if _, err := c.db.Exec(ctx, query, request.NewPassword, request.ID); err != nil {
		fmt.Println("error while updating password for user", err.Error())
		return err
	}

	return nil
}

func (c *customerRepo) GetCustomerCredentialsByLogin(ctx context.Context, login string) (models.Customer, error) {
	user := models.Customer{}

	query := `
		select id, password from customers 
		                where login = $1`

	if err := c.db.QueryRow(ctx, query, login).Scan(&user.ID, &user.Password); err != nil {
		fmt.Println("Error while scanning password from customers", err.Error())
		return models.Customer{}, err
	}

	return user, nil
}
