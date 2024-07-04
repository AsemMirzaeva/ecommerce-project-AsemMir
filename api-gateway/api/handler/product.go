package handler

import (
	prod "api-gateway/protos/product-service"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type ProductHandler struct {
	client prod.ProductServiceClient
}

func NewProductHandler(conn *grpc.ClientConn) *ProductHandler {
	client := prod.NewProductServiceClient(conn)
	return &ProductHandler{client: client}
}

// CreateProduct ...
// @Summary CreateProduct
// @Security ApiKeyAuth
// @Description Api for creating a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param Product body models.CreateProduct true "createProductModel"
// @Success 200 {object} models.Product
// @Router /products [post]
func (h *ProductHandler) AddProduct(c *gin.Context) {
	var req prod.AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.AddProduct(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProduct ...
// @Summary GetProduct
// @Security ApiKeyAuth
// @Description Api for creatgetting a product by id
// @Tags Products
// @Accept json
// @Produce json
// @Param  id path string true "Product_id"
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	req := prod.GetProductRequest{Id: id}
	resp, err := h.client.GetProduct(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateProduct updates Product by id
// @Summary UpdateProduct
// @Security ApiKeyAuth
// @Description Api for updating products by id
// @Tags Products
// @Accept json
// @Produce json
// @Param  id path string true "Product_id"
// @Param Product body models.UpdateProduct true "updateProductModel"
// @Success 200 {object} models.Product
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req = prod.UpdateProductRequest{
		Id: id,
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.client.UpdateProduct(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteProduct deletes Product by id
// @Summary DeleteProduct
// @Security ApiKeyAuth
// @Description Api for deleting products by id
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.Product
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)
	req := prod.DeleteProductRequest{Id: id}
	resp, err := h.client.DeleteProduct(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ListProducts godoc
// @Summary List all Products
// @Description Retrieve a list of all products
// @Tags Products
// @Accept  json
// @Produce  json
// @Param limit query string false "Limit"
// @Param page query string false "Page"
// @Success 200 {array} models.Product
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	req := prod.ListProductsRequest{
		Limit: limit,
		Page:  page,
	}
	stream, err := h.client.ListProducts(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var products []prod.ListProductsResponse
	for {
		product, err := stream.Recv()
		if err != nil {
			break
		}
		products = append(products, *product)
	}
	c.JSON(http.StatusOK, products)
}
