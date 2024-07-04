package repo

import "order-service/models"

type OrderStorage interface {
	CreateOrder(order *models.Order) error
	GetOrder(id string) (*models.Order, error)
	DeleteOrder(id string) error
	ListOrders() ([]*models.Order, error)
}
