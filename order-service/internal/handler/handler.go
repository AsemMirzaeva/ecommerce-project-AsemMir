package handler

import (
	"context"

	"log"
	"order-service/internal/config"
	proto "order-service/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Handler struct {
	proto.UnimplementedOrderServiceServer
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	orderID := uuid.New().String()
	// Assuming you have an order repository to save the order
	// err := h.orderRepo.CreateOrder(orderID, req.UserId, req.ProductId, req.Quantity)
	if err != nil {
		log.Logger.Printf("error creating order: %v", err)
		return nil, err
	}

	return &proto.CreateOrderResponse{Id: orderID}, nil
}

func (h *Handler) CreateOrders(stream proto.OrderService_CreateOrdersServer) error {
	for {
		req, err := stream.Recv()
		if err == grpc.ErrServerStopped {
			break
		}
		if err != nil {
			log.Logger.Printf("error receiving orders: %v", err)
			return err
		}

		// Process the orders
		orderID := uuid.New().String()
		// err := h.orderRepo.CreateOrder(orderID, req.UserId, req.ProductIds)
		if err != nil {
			log.Logger.Printf("error creating orders: %v", err)
			return err
		}

		if err := stream.Send(&proto.CreateOrdersResponse{
			Id: orderID,
		}); err != nil {
			log.Logger.Printf("error sending orders: %v", err)
			return err
		}
	}

	return nil
}
