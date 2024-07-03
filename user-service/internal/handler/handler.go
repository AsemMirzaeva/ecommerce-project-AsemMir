package handler

import (
	"context"
	"log"

	pb "user-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	users map[string]*pb.RegisterResponse // Simulating a user database
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userId := "user_" + req.Username
	h.users[userId] = &pb.RegisterResponse{
		Id:       userId,
		Username: req.Username,
		Email:    req.Email,
	}

	log.Println("Registered user:", h.users[userId])
	return &pb.RegisterResponse{
		Id:       userId,
		Username: req.Username,
		Email:    req.Email,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	for _, user := range h.users {
		if user.Email == req.Email {
			// This is a placeholder. Implement proper authentication.
			return &pb.LoginResponse{Token: "dummy-token-for-" + user.Id}, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, exists := h.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	return &pb.GetUserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	delete(h.users, req.Id)
	return &pb.DeleteUserResponse{Status: "User deleted successfully"}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, exists := h.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	user.Username = req.Username
	user.Email = req.Email
	h.users[req.Id] = user

	return &pb.UpdateUserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func RegisterUserHandler(server *grpc.Server) {
	userHandler := &UserHandler{
		users: make(map[string]*pb.RegisterResponse),
	}
	pb.RegisterUserServiceServer(server, userHandler)
}
