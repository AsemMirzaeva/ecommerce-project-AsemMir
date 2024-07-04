package grpcClient

import (
	"fmt"
	"order-service/config"
	pbp "order-service/protos/product-service"
	pbu "order-service/protos/user-service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IServiceManager interface {
	ProductService() pbp.ProductServiceClient
	UserService() pbu.UserServiceClient
}

type serviceManager struct {
	cfg            config.Config
	productService pbp.ProductServiceClient
	userService    pbu.UserServiceClient
}

func New(cfg config.Config) (IServiceManager, error) {
	// Dial to Product-service
	connProduct, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.ProductServiceHost, cfg.ProductServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial Product service at %s:%s: %w", cfg.ProductServiceHost, cfg.ProductServicePort, err)
	}

	// Dial to User-service
	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial User service at %s:%s: %w", cfg.UserServiceHost, cfg.UserServicePort, err)
	}

	return &serviceManager{
		cfg:            cfg,
		productService: pbp.NewProductServiceClient(connProduct),
		userService:    pbu.NewUserServiceClient(connUser),
	}, nil
}

func (s *serviceManager) ProductService() pbp.ProductServiceClient {
	return s.productService
}

func (s *serviceManager) UserService() pbu.UserServiceClient {
	return s.userService
}
