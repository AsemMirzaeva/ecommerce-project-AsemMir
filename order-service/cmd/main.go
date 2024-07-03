package main

import (
	"context"
	"net"
	"order-service/storage"
	"os"
	"os/signal"
	pb "order-service/proto"
	"path/to/pkg/interceptor"
	"path/to/pkg/logger"
	productpb "path/to/product/proto"
	userpb "path/to/user/proto"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	db *storage.OrderRepo
}

type Order struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	conn, err := grpc.Dial(os.Getenv("USERSERVICE_HOST")+":"+os.Getenv("USERSERVICE_PORT"), grpc.WithInsecure())
	if err != nil {
		logger.ErrorLogger.Println("Failed to connect to UserService:", err)
		return nil, err
	}
	defer conn.Close()

	userClient := userpb.NewUserServiceClient(conn)
	_, err = userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.GetUserId()})
	if err != nil {
		logger.ErrorLogger.Println("Failed to get user from UserService:", err)
		return nil, err
	}

	conn, err = grpc.Dial(os.Getenv("PRODUCTSERVICE_HOST")+":"+os.Getenv("PRODUCTSERVICE_PORT"), grpc.WithInsecure())
	if err != nil {
		logger.ErrorLogger.Println("Failed to connect to ProductService:", err)
		return nil, err
	}
	defer conn.Close()

	productClient := productpb.NewProductServiceClient(conn)
	_, err = productClient.GetProduct(ctx, &productpb.GetProductRequest{Id: req.GetProductId()})
	if err != nil {
		logger.ErrorLogger.Println("Failed to get product from ProductService:", err)
		return nil, err
	}

	order := Order{UserID: uuid.MustParse(req.GetUserId()), ProductID: uuid.MustParse(req.GetProductId()), Quantity: req.GetQuantity()}
	result := s.db.Create(&order)
	if result.Error != nil {
		logger.ErrorLogger.Println("Failed to create order:", result.Error)
		return nil, result.Error
	}
	logger.InfoLogger.Println("Order created successfully:", order.ID)
	return &pb.CreateOrderResponse{Id: order.ID.String()}, nil
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		logger.ErrorLogger.Fatalf("Error loading .env file")
	}

	logger.InitLogger()
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Order{})

	lis, err := net.Listen("tcp", os.Getenv("ORDERSERVICE_HOST")+":"+os.Getenv("ORDERSERVICE_PORT"))
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryInterceptor),
		grpc.StreamInterceptor(interceptor.StreamInterceptor),
	)
	pb.RegisterOrderServiceServer(s, &server{db: db})
	logger.InfoLogger.Printf("Order service listening on %s:%s", os.Getenv("ORDERSERVICE_HOST"), os.Getenv("ORDERSERVICE_PORT"))

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.InfoLogger.Println("Shutting down server...")
		s.GracefulStop()
		logger.InfoLogger.Println("Server stopped")
		os.Exit(0)
	}()

	if err := s.Serve(lis); err != nil {
		logger.ErrorLogger.Fatalf("Failed to serve: %v", err)
	}
}
