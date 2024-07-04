package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"user-service/config"
	"user-service/pkg/logger"
	userpb "user-service/protos/user-service"
	"user-service/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	logg := logger.New("debug", "user-service")
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

	userService := service.NewUserService(db, logg)

	fmt.Println("Server is running on port :", cfg.UserServicePort)
	lis, err := net.Listen("tcp", ":"+cfg.UserServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, userService)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
