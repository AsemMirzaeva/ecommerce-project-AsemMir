package handler

import (
	"context"

	"api-gateway/internal/config"
	"api-gateway/proto"
	"log"

	"google.golang.org/grpc"
)

type Handler struct {
	proto.UnimplementedAPIGatewayServiceServer
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	conn, err := grpc.Dial(h.cfg.UserServiceHost+":"+h.cfg.UserServicePort, grpc.WithInsecure())
	if err != nil {
		log.Logger.Printf("error connecting to user service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := proto.NewUserServiceClient(conn)
	return client.CreateUser(ctx, req)
}

func (h *Handler) AddProduct(ctx context.Context, req *proto.AddProductRequest) (*proto.AddProductResponse, error) {
	conn, err := grpc.Dial(h.cfg.ProductServiceHost+":"+h.cfg.ProductServicePort, grpc.WithInsecure())
	if err != nil {
		log.Logger.Printf("error connecting to product service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := proto.NewProductServiceClient(conn)
	return client.AddProduct(ctx, req)
}

func (h *Handler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	conn, err := grpc.Dial(h.cfg.OrderServiceHost+":"+h.cfg.OrderServicePort, grpc.WithInsecure())
	if err != nil {
		log.Logger.Printf("error connecting to order service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := proto.NewOrderServiceClient(conn)
	return client.CreateOrder(ctx, req)
}

func (h *Handler) ListProducts(req *proto.ListProductsRequest, stream proto.APIGatewayService_ListProductsServer) error {
	conn, err := grpc.Dial(h.cfg.ProductServiceHost+":"+h.cfg.ProductServicePort, grpc.WithInsecure())
	if err != nil {
		log.Logger.Printf("error connecting to product service: %v", err)
		return err
	}
	defer conn.Close()

	client := proto.NewProductServiceClient(conn)
	productStream, err := client.ListProducts(context.Background(), req)
	if err != nil {
		log.Logger.Printf("error listing products: %v", err)
		return err
	}

	for {
		productResp, err := productStream.Recv()
		if err == grpc.ErrServerStopped {
			break
		}
		if err != nil {
			log.Logger.Printf("error receiving product: %v", err)
			return err
		}

		if err := stream.Send(productResp); err != nil {
			log.Logger.Printf("error sending product: %v", err)
			return err
		}
	}

	return nil
}

func (h *Handler) CreateOrders(stream proto.APIGatewayService_CreateOrdersServer) error {
	conn, err := grpc.Dial(h.cfg.OrderServiceHost+":"+h.cfg.OrderServicePort, grpc.WithInsecure())
	if err != nil {
		log.Logger.Printf("error connecting to order service: %v", err)
		return err
	}
	defer conn.Close()

	client := proto.NewOrderServiceClient(conn)
	orderStream, err := client.CreateOrders(context.Background())
	if err != nil {
		log.Logger.Printf("error creating orders: %v", err)
		return err
	}

	for {
		req, err := stream.Recv()
		if err == grpc.ErrServerStopped {
			break
		}
		if err != nil {
			log.Logger.Printf("error receiving order: %v", err)
			return err
		}

		if err := orderStream.Send(req); err != nil {
			log.Logger.Printf("error sending order: %v", err)
			return err
		}

		resp, err := orderStream.Recv()
		if err != nil {
			log.Logger.Printf("error receiving order response: %v", err)
			return err
		}

		if err := stream.Send(resp); err != nil {
			log.Logger.Printf("error sending order response: %v", err)
			return err
		}
	}

	return nil
}
