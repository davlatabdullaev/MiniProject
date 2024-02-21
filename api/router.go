package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"test/api/handler"
	"test/pkg/logger"
	"test/service"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// New ...
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
func New(services service.IServiceManager, log logger.ILogger) *gin.Engine {
	h := handler.New(services, log)

	r := gin.New()

	//r.Use(authenticateMiddleware)
	r.Use(gin.Logger())

	{
		// auth endpoints
		r.POST("/auth/customer/login", h.CustomerLogin)

		// user endpoints
		r.POST("/customer", h.CreateCustomer)
		r.GET("/customer/:id", h.GetCustomer)
		r.GET("/customers", h.GetCustomerList)
		r.PUT("/customer/:id", h.UpdateCustomer)
		r.DELETE("/customer/:id", h.DeleteCustomer)
		r.PATCH("/customer/:id", h.UpdateCustomerPassword)
		// product endpoints
		r.POST("/product", h.CreateProduct)
		r.GET("/product/:id", h.GetProduct)
		r.GET("/products", h.GetProductList)
		r.PUT("/product/:id", h.UpdateProduct)
		r.DELETE("/product/:id", h.DeleteProduct)

		// basket endpoints
		r.POST("/basket", h.CreateBasket)
		r.GET("/basket/:id", h.GetBasket)
		r.GET("/baskets", h.GetBasketList)
		r.PUT("basket/:id", h.UpdateBasket)
		r.DELETE("basket/:id", h.DeleteBasket)

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		r.Use(traceRequest)
		r.Use(authenticateMiddleware)

	}

	return r
}

func authenticateMiddleware(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
	} else {
		c.Next()
	}
}

func traceRequest(c *gin.Context) {
	beforeRequest(c)

	c.Next()

	afterRequest(c)
}

func beforeRequest(c *gin.Context) {
	startTime := time.Now()

	c.Set("start_time", startTime)

	log.Println("start time:", startTime.Format("2006-01-02 15:04:05.0000"), "path:", c.Request.URL.Path)
}

func afterRequest(c *gin.Context) {
	startTime, exists := c.Get("start_time")
	if !exists {
		startTime = time.Now()
	}

	duration := time.Since(startTime.(time.Time)).Seconds()

	log.Println("end time:", time.Now().Format("2006-01-02 15:04:05.0000"), "duration:", duration, "method:", c.Request.Method)
	fmt.Println()
}
