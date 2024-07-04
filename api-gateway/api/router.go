package api

import (
	_ "api-gateway/api/docs" // swag

	"api-gateway/api/handler"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

// @Title Welcome to swagger service
// @Version 1.0
// @Description  This is Asem's E-Commerce
func NewRouter(conn, connP, orderConn *grpc.ClientConn) *gin.Engine {

	router := gin.Default()

	userHandler := handler.NewUserHandler(conn)
	productHandler := handler.NewProductHandler(connP)
	orderHandler := handler.NewOrderHandler(orderConn)

	// User routes
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/:id", userHandler.GetUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
	router.GET("/users", userHandler.ListUsers)

	// // Product routes
	router.POST("/products", productHandler.AddProduct)
	router.GET("/products/:id", productHandler.GetProduct)
	router.PUT("/products/:id", productHandler.UpdateProduct)
	router.DELETE("/products/:id", productHandler.DeleteProduct)
	router.GET("/products", productHandler.ListProducts)

	// // Order routes
	router.POST("/orders", orderHandler.CreateOrder)
	router.GET("/orders/:id", orderHandler.GetOrder)
	router.DELETE("/orders/:id", orderHandler.DeleteOrder)
	router.GET("/orders", orderHandler.ListOrders)
	router.POST("/orders/bulk", orderHandler.CreateOrders)

	url := ginSwagger.URL("swaggerdoc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
