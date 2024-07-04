package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"order-service/config"
	"order-service/pkg/logger"
	orderpb "order-service/protos/order-service"
	"order-service/service"

	grpcClient "order-service/service/grpc_client"

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

	logg := logger.New("debug", "order-service")
	defer logger.Cleanup(logg)

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	grpcClient, err := grpcClient.New(*cfg)
	if err != nil {
		log.Fatal("grpc client dail error", logger.Error(err))
	}

	db, err := sql.Open("postgres", "host="+cfg.Postgres.Host+" port="+cfg.Postgres.Port+" user="+cfg.Postgres.User+" password="+cfg.Postgres.Password+" dbname="+cfg.Postgres.Database+" sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	orderService := service.NewOrderService(db, logg, grpcClient)

	fmt.Println("Server is running on port :", cfg.OrderServicePort)
	lis, err := net.Listen("tcp", ":"+cfg.OrderServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(s, orderService)
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
