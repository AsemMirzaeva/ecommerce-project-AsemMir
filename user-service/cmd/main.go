package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"user-service/internal/config"
	"user-service/internal/interceptor/logger"
	pb "user-service/proto"
	storage "user-service/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	db *storage.UserRepo
}

type User struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name  string
	Email string
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

	server := userServiceServer{db: storage.NewUserRepo(db)}

	address := fmt.Sprintf(":%s", cfg.UserServicePort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, server)

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
