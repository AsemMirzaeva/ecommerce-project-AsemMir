package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	userpb "path/to/user/proto"
	productpb "path/to/product/proto"
	orderpb "path/to/order/proto"
	"path/to/pkg/logger"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		logger.ErrorLogger.Fatalf("Error loading .env file")
	}

	logger.InitLogger()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, os.Getenv("USERSERVICE_HOST")+":"+os.Getenv("USERSERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway for User Service: %v", err)
	}

	err = productpb.RegisterProductServiceHandlerFromEndpoint(ctx, mux, os.Getenv("PRODUCTSERVICE_HOST")+":"+os.Getenv("PRODUCTSERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway for Product Service: %v", err)
	}

	err = orderpb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, os.Getenv("ORDERSERVICE_HOST")+":"+os.Getenv("ORDERSERVICE_PORT"), opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway for Order Service: %v", err)
	}

	logger.InfoLogger.Println("HTTP Gateway is running...")
	if err := http.ListenAndServe(os.Getenv("APIGATEWAY_HOST")+":"+os.Getenv("APIGATEWAY_PORT"), mux); err != nil {
		log.Fatalf("Failed to serve HTTP gateway: %v", err)
	}
}





// package main

// import (
//     "context"
//     "log"
//     "net"
//     "net/http"
//     "os"
//     "os/signal"
//     "syscall"

//     "api-gateway/internal/config"
//     "api-gateway/internal/handler"
//     "api-gateway/internal/proto"
   
//     "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
//     "google.golang.org/grpc"
//     "google.golang.org/grpc/credentials/insecure"
// )

// func main() {
//     log.InitLogger()

//     cfg := config.LoadConfig()

//     grpcServer := grpc.NewServer()

//     // Register gRPC server
//     proto.RegisterAPIGatewayServiceServer(grpcServer, handler.NewHandler(cfg))

//     // Start gRPC server
//     lis, err := net.Listen("tcp", ":"+cfg.APIGatewayPort)
//     if err != nil {
//         log.Logger.Fatalf("failed to listen: %v", err)
//     }

//     go func() {
//         log.Logger.Printf("starting gRPC server on port %s", cfg.APIGatewayPort)
//         if err := grpcServer.Serve(lis); err != nil {
//             log.Logger.Fatalf("failed to serve: %v", err)
//         }
//     }()

//     // Start HTTP server for gRPC-Gateway
//     mux := runtime.NewServeMux()
//     opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
//     err = proto.RegisterAPIGatewayServiceHandlerFromEndpoint(context.Background(), mux, ":"+cfg.APIGatewayPort, opts)
//     if err != nil {
//         log.Logger.Fatalf("failed to register HTTP server: %v", err)
//     }

//     go func() {
//         log.Logger.Printf("starting HTTP server on port %s", cfg.APIGatewayPort)
//         if err := http.ListenAndServe(":"+cfg.APIGatewayPort, mux); err != nil {
//             log.Logger.Fatalf("failed to serve: %v", err)
//         }
//     }()

//     // Graceful shutdown
//     sigChan := make(chan os.Signal, 1)
//     signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

//     <-sigChan

//     log.Logger.Println("gracefully shutting down...")
//     grpcServer.GracefulStop()
//     log.Logger.Println("shutdown complete")
// }
