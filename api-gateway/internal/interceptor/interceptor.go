package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerIntenceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerIntenceptor, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("UnaryServerIntenceptor PRE", info.FullMethod)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("no metadata")
	} else {
		log.Printf("Metadata %v+\n", md)
	}

	m, err := handler(ctx, req)

	log.Println("UnaryServerInterceptor POST", info.FullMethod)

	return m, err

}

type wrappedStream struct {
	grpc.ServerStream
}
