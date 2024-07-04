package main

import (
    "database/sql"
    "fmt"
    "log"
    "net"
    "order-service/config"
    "order-service/handlers"
    "order-service/repository"
    "order-service/service"
    orderpb "order-service/proto/orderproto"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    "google.golang.org/grpc"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    cfg := config.LoadConfig()
    db, err := sql.Open("postgres", fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
    ))
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
    }
    defer db.Close()

    repo := repository.NewPostgresRepository(db)
    svc := service.NewOrderService(repo)

    productClient, err := service.NewProductClient(":" + cfg.USER_SERVICE_PORT)
    if err != nil {
        log.Fatalf("failed to create product client: %v", err)
    }

    svc.SetProductClient(productClient)

    server := handlers.NewServer(svc)

    fmt.Println("Server is running on port:", cfg.PORT)
    lis, err := net.Listen("tcp", ":" + cfg.PORT)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    orderpb.RegisterOrderServiceServer(s, server)

    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
