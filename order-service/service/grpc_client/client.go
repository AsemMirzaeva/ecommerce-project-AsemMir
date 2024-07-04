package grpcClient

import (
	"fmt"
	"order-service/config"
	pbp "order-service/protos/product-service"
	pbu "order-service/protos/user-service"

	"google.golang.org/grpc"
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
	// dail to Product-service
	connProduct, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.ProductServiceHost, cfg.ProductServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("user service dail host: %s port : %s", cfg.ProductServiceHost, cfg.ProductServicePort)
	}
	// dail to User-service
	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("user service dail host: %s port : %s", cfg.UserServiceHost, cfg.UserServicePort)
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
