package handler

import (
	ord "api-gateway/protos/order-service"
	prod "api-gateway/protos/product-service"
	user "api-gateway/protos/user-service"
)

type Handler struct {
	user    user.UserServiceClient
	product prod.ProductServiceClient
	order   ord.OrderServiceClient
}

func NewHandler(user user.UserServiceClient, product prod.ProductServiceClient, order ord.OrderServiceClient) *Handler {
	return &Handler{user: user, product: product}
}
