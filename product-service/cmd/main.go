package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"product-service/config"
	"product-service/internal"
	"product-service/pkg/logger"
	productpb "product-service/protos/product-service"
	"product-service/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	logg := logger.New("debug", "product-service")
	defer logger.Cleanup(logg)

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := sql.Open("postgres", "host="+cfg.Postgres.Host+" port="+cfg.Postgres.Port+" user="+cfg.Postgres.User+" password="+cfg.Postgres.Password+" dbname="+cfg.Postgres.Database+" sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	productService := service.NewProductService(db, logg)

	fmt.Println("Server is running on port :", cfg.ProductServicePort)
	lis, err := net.Listen("tcp", ":"+cfg.ProductServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.StreamInterceptor(internal.StreamInterceptor),
		grpc.UnaryInterceptor(internal.UnaryInterceptor))

	productpb.RegisterProductServiceServer(s, productService)
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
