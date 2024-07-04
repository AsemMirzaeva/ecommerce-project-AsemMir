package internal

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func logger(format string, a ...any) {
	fmt.Printf("LOG:\t"+format+"\n", a...)
}

type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) SendMsg(m any) error {
	logger("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func NewWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func StreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, NewWrappedStream(ss))
	if err != nil {
		logger("RPC failed with error: %v", err)
	}
	return err
}

func UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	logger("Inside the unaryInterceptor")
	m, err := handler(ctx, req)
	if err != nil {
		logger("RPC failed with error: %v", err)
	}
	logger("Exiting the unaryInterceptor")
	return m, err
}
