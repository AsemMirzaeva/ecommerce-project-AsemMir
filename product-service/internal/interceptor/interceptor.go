package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("gRPC method: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return resp, err
}

func StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Printf("gRPC stream method: %s", info.FullMethod)
	err := handler(srv, ss)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return err
}
