package service

import (
	"context"
	"errors"
	"fmt"
	"test/api/models"
	"test/pkg/check"
	"test/pkg/logger"
	"test/pkg/security"
	"test/storage"

	"github.com/jackc/pgx"
)

type customerService struct {
	storage storage.IStorage
	log     logger.ILogger
}

func NewCustomerService(storage storage.IStorage, log logger.ILogger) customerService {
	return customerService{
		storage: storage,
		log:     log,
	}
}

func (c customerService) Create(ctx context.Context, createUser models.CreateCustomer) (models.Customer, error) {
	c.log.Info("User create service layer", logger.Any("createUser", createUser))

	password, err := security.HashPassword(createUser.Password)
	if err != nil {
		c.log.Error("error while hashing password", logger.Error(err))
		return models.Customer{}, err
	}
	createUser.Password = password

	pKey, err := c.storage.Customer().Create(ctx, createUser)
	if err != nil {
		c.log.Error("error while creating user", logger.Error(err))
		return models.Customer{}, err
	}

	customer, err := c.storage.Customer().GetByID(ctx, models.PrimaryKey{
		ID: pKey,
	})

	return customer, nil
}

func (c customerService) GetUser(ctx context.Context, pKey models.PrimaryKey) (models.Customer, error) {
	user, err := c.storage.Customer().GetByID(ctx, pKey)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("ERROR in service layer while getting user by id", err.Error())
			return models.Customer{}, err
		}
	}

	return user, nil
}

func (c customerService) GetUsers(ctx context.Context, request models.GetListRequest) (models.CustomersResponse, error) {
	c.log.Info("Get user list service layer", logger.Any("request", request))
	usersResponse, err := c.storage.Customer().GetList(ctx, request)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			c.log.Error("error while getting users list", logger.Error(err))
			return models.CustomersResponse{}, err
		}
	}

	return usersResponse, err
}

func (c customerService) Update(ctx context.Context, updateUser models.UpdateCustomer) (models.Customer, error) {
	pKey, err := c.storage.Customer().Update(ctx, updateUser)
	if err != nil {
		c.log.Error("ERROR in service layer while updating updateUser", logger.Error(err))
		return models.Customer{}, err
	}

	user, err := c.storage.Customer().GetByID(ctx, models.PrimaryKey{
		ID: pKey,
	})
	if err != nil {
		c.log.Error("ERROR in service layer while getting user after update", logger.Error(err))
		return models.Customer{}, err
	}

	return user, nil
}

func (c customerService) Delete(ctx context.Context, key models.PrimaryKey) error {
	err := c.storage.Customer().Delete(ctx, key)
	return err
}

func (c customerService) UpdatePassword(ctx context.Context, request models.UpdateCustomerPassword) error {
	oldPassword, err := c.storage.Customer().GetPassword(ctx, request.ID)
	if err != nil {
		c.log.Error("ERROR in service layer while getting user password", logger.Error(err))
		return err
	}

	if oldPassword != request.OldPassword {
		c.log.Error("ERROR in service old password is not correct")
		return errors.New("old password did not match")
	}

	if err = check.ValidatePassword(request.NewPassword); err != nil {
		c.log.Error("ERROR in service layer new password is weak", logger.Error(err))
		return err
	}

	if err = c.storage.Customer().UpdatePassword(context.Background(), request); err != nil {
		c.log.Error("ERROR in service layer while updating password", logger.Error(err))
		return err
	}

	return nil
}
