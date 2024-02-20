package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"test/api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateCustomer godoc
// @Router       /customer [POST]
// @Summary      Creates a new customer
// @Description  create a new customer
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        customer body models.CreateCustomer false "customer"
// @Success      201  {object}  models.Customer
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateCustomer(c *gin.Context) {
	createCustomer := models.CreateCustomer{}

	if err := c.ShouldBindJSON(&createCustomer); err != nil {
		handleResponse(c, h.log, "error while reading body from client", http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := h.services.Customer().Create(ctx, createCustomer)
	if err != nil {
		handleResponse(c, h.log, "error while creating Customer", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusCreated, resp)
}

// GetCustomer godoc
// @Router       /Customer/{id} [GET]
// @Summary      Gets Customer
// @Description  get Customer by ID
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer"
// @Success      200  {object}  models.Customer
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetCustomer(c *gin.Context) {
	var err error

	uid := c.Param("id")

	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "invalid uuid type", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	Customer, err := h.services.Customer().GetCustomer(ctx, models.PrimaryKey{
		ID: id.String(),
	})
	if err != nil {
		handleResponse(c, h.log, "error while getting Customer by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, Customer)
}

// GetCustomerList godoc
// @Router       /Customers [GET]
// @Summary      Get Customer list
// @Description  get Customer list
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param        page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.CustomersResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetCustomerList(c *gin.Context) {
	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, h.log, "error while parsing page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, h.log, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := h.services.Customer().GetCustomersList(ctx, models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(c, h.log, "error while getting Customers", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "success!", http.StatusOK, resp)
}

// UpdateCustomer godoc
// @Router       /Customer/{id} [PUT]
// @Summary      Update Customer
// @Description  update Customer
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param 		 id path string true "Customer_id"
// @Param        Customer body models.UpdateCustomer true "Customer"
// @Success      200  {object}  models.Customer
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateCustomer(c *gin.Context) {
	updateCustomer := models.UpdateCustomer{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, h.log, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateCustomer.ID = uid

	if err := c.ShouldBindJSON(&updateCustomer); err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := h.services.Customer().Update(ctx, updateCustomer)
	if err != nil {
		handleResponse(c, h.log, "error while updating Customer", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, resp)
}

// DeleteCustomer godoc
// @Router       /Customer/{id} [DELETE]
// @Summary      Delete Customer
// @Description  delete Customer
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param 		 id path string true "Customer_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteCustomer(c *gin.Context) {
	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, h.log, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = h.services.Customer().Delete(ctx, models.PrimaryKey{
		ID: id.String(),
	}); err != nil {
		handleResponse(c, h.log, "error while deleting Customer by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, "data successfully deleted")
}

// UpdateCustomerPassword godoc
// @Router       /Customer/{id} [PATCH]
// @Summary      Update Customer password
// @Description  update Customer password
// @Tags         customer
// @Accept       json
// @Produce      json
// @Param 		 id path string true "Customer_id"
// @Param        Customer body models.UpdateCustomerPassword true "Customer"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateCustomerPassword(c *gin.Context) {
	updateCustomerPassword := models.UpdateCustomerPassword{}

	if err := c.ShouldBindJSON(&updateCustomerPassword); err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handleResponse(c, h.log, "error while parsing uuid", http.StatusBadRequest, err.Error())
		return
	}

	updateCustomerPassword.ID = uid.String()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = h.services.Customer().UpdatePassword(ctx, updateCustomerPassword); err != nil {
		handleResponse(c, h.log, "error while updating Customer password", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, "password successfully updated")
}
