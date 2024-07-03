package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"product-service/internal/config"
	"product-service/internal/interceptor/logger"
	pb "product-service/proto"
	"product-service/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type productServiceService struct {
	pb.UnimplementedProductServiceServer
	db *storage.ProductRepo
}

type Product struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name  string
	Price float64
}

func (s *productServiceService) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	product := Product{Name: req.GetName(), Price: req.GetPrice()}
	result := s.db.Create(&product)
	if result.Error != nil {
		logger.ErrorLogger.Println("Failed to add product:", result.Error)
		return nil, result.Error
	}
	logger.InfoLogger.Println("Product added successfully:", product.ID)
	return &pb.AddProductResponse{Id: product.ID.String()}, nil
}

func (s *productServiceService) ListProducts(req *pb.ListProductsRequest, stream pb.ProductService_ListProductsServer) error {
	var products []Product
	s.db.Find(&products)
	for _, product := range products {
		if err := stream.Send(&pb.ListProductsResponse{Name: product.Name, Price: product.Price}); err != nil {
			logger.ErrorLogger.Println("Failed to list products:", err)
			return err
		}
	}
	return nil
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := storage.ConnectDB(*cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	server := productServiceService{db: storage.NewProductRepo(db)}

	address := fmt.Sprintf(":%s", cfg.ProductServicePort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterProductServiceServer(s, server)

	log.Printf("Server listening on port %v", lis.Addr())

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

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
