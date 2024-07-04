package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"order-service/models"
	l "order-service/pkg/logger"
	pbo "order-service/protos/order-service"
	pbp "order-service/protos/product-service"
	pbu "order-service/protos/user-service"
	grpcClient "order-service/service/grpc_client"
	"order-service/storage"

	"github.com/google/uuid"
)

type OrderService struct {
	pbo.UnimplementedOrderServiceServer
	storage storage.IStorage
	logger  l.Logger
	client  grpcClient.IServiceManager
}

func NewOrderService(db *sql.DB, log l.Logger, client grpcClient.IServiceManager) *OrderService {
	return &OrderService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pbo.CreateOrderRequest) (*pbo.CreateOrderResponse, error) {
	productServiceClient := s.client.ProductService()
	product, err := productServiceClient.GetProduct(ctx, &pbp.GetProductRequest{Id: req.ProductId})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	userServiceClient := s.client.UserService()
	user, err := userServiceClient.GetUser(ctx, &pbu.GetUserRequest{Id: req.UserId})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	order := models.Order{
		ID:         uuid.New().String(),
		UserID:     user.Id,
		ProductID:  product.Id,
		Quantity:   req.Quantity,
		Status:     "created",
		TotalPrice: float32(math.Round(float64(req.Quantity)*float64(product.Price)*100) / 100),
	}
	fmt.Println(order.TotalPrice)
	err = s.storage.Order().CreateOrder(&order)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &pbo.CreateOrderResponse{
		Id:         order.ID,
		TotalPrice: float32(order.TotalPrice),
	}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *pbo.GetOrderRequest) (*pbo.GetOrderResponse, error) {
	order, err := s.storage.Order().GetOrder(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pbo.GetOrderResponse{
		Id:         order.ID,
		UserId:     order.UserID,
		ProductId:  order.ProductID,
		Quantity:   order.Quantity,
		TotalPrice: float32(order.TotalPrice),
	}, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, req *pbo.DeleteOrderRequest) (*pbo.DeleteOrderResponse, error) {
	err := s.storage.Order().DeleteOrder(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pbo.DeleteOrderResponse{Message: "Order deleted"}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, req *pbo.ListOrdersRequest) (*pbo.ListOrdersResponse, error) {
	orders, err := s.storage.Order().ListOrders()
	if err != nil {
		s.logger.Error(err.Error())
		return &pbo.ListOrdersResponse{}, nil
	}
	var result []*pbo.GetOrderResponse
	for _, order := range orders {
		result = append(result, &pbo.GetOrderResponse{
			Id:         order.ID,
			UserId:     order.UserID,
			ProductId:  order.ProductID,
			Quantity:   order.Quantity,
			TotalPrice: float32(order.TotalPrice),
		})
	}
	return &pbo.ListOrdersResponse{Orders: result}, nil

}
