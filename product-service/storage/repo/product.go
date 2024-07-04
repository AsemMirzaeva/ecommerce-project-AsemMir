package repo

import "product-service/models"

type ProductStorage interface {
	AddProduct(product *models.Product) (*models.Product, error)
    GetProduct(id string) (*models.Product, error)
    UpdateProduct(product *models.Product) (*models.Product, error)
    DeleteProduct(id string) error
    ListProducts(limit, page int) ([]*models.Product, error)
}