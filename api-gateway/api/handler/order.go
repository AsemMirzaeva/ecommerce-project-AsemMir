package handler

import (
	"context"
	"net/http"

	order "api-gateway/protos/order-service"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type OrderHandler struct {
	client order.OrderServiceClient
}

func NewOrderHandler(conn *grpc.ClientConn) *OrderHandler {
	client := order.NewOrderServiceClient(conn)
	return &OrderHandler{client: client}
}

// CreateOrder godoc
// @Summary Create a new order
// @Security ApiKeyAuth
// @Description API for creating a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param Order body models.CreateOrder true "createOrderModel"
// @Success 200 {object} models.Order
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.CreateOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetOrder godoc
// @Summary Get order by ID
// @Security ApiKeyAuth
// @Description API for getting order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	req := order.GetOrderRequest{Id: id}
	resp, err := h.client.GetOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteOrder godoc
// @Summary Delete order by ID
// @Security ApiKeyAuth
// @Description API for deleting order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	req := order.DeleteOrderRequest{Id: id}
	resp, err := h.client.DeleteOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ListOrders godoc
// @Summary List all orders
// @Description Retrieve a list of all orders
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Success 200 {array} models.Order
// @Router /orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")
	req := order.ListOrdersRequest{
		Limit: limit,
		Page:  page,
	}
	resp, err := h.client.ListOrders(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// // CreateOrders godoc
// // @Summary Create multiple orders
// // @Security ApiKeyAuth
// // @Description API for creating multiple orders
// // @Tags Orders
// // @Accept json
// // @Produce json
// // @Success 200 {object} models.Order
// // @Router /orders [post]
// func (h *OrderHandler) CreateOrders(c *gin.Context) {
// 	stream, err := h.client.CreateOrders(context.Background())
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var req order.CreateOrderRequest
// 	for {
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			break
// 		}
// 		if err := stream.Send(&req); err != nil {
// 			break
// 		}
// 	}

// 	resp, err := stream.CloseAndRecv()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, resp)
// }
