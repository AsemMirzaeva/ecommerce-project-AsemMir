package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"
	"user-service/models"
	l "user-service/pkg/logger"
	"user-service/storage"

	pbu "user-service/protos/user-service"
)

// UserService ...
type UserService struct {
	pbu.UnimplementedUserServiceServer
	storage storage.IStorage
	logger  l.Logger
}

// NewUserService ...
func NewUserService(db *sql.DB, log l.Logger) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *UserService) GetUser(ctx context.Context, req *pbu.GetUserRequest) (*pbu.GetUserResponse, error) {
	user, err := s.storage.User().GetUser(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pbu.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pbu.CreateUserRequest) (*pbu.CreateUserResponse, error) {
	user, err := s.storage.User().CreateUser(&models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pbu.CreateUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pbu.UpdateUserRequest) (*pbu.UpdateUserResponse, error) {
	user, err := s.storage.User().UpdateUser(&models.User{
		ID:       req.Id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pbu.UpdateUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pbu.DeleteUserRequest) (*pbu.DeleteUserResponse, error) {
	err := s.storage.User().DeleteUser(req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pbu.DeleteUserResponse{
		Message: "delete user successfully",
	}, nil
}

func (s *UserService) ListUsers(req *pbu.ListUsersRequest, stream pbu.UserService_ListUsersServer) error {
	limit, err := strconv.Atoi(req.Limit)
	if err != nil {
		return err
	}
	page, err := strconv.Atoi(req.Page)
	if err != nil {
		return err
	}
	if limit <= 0 {
		return errors.New("limit must be greater than zero")
	}

	users, err := s.storage.User().ListUsers(limit, page)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	for _, user := range users {
		if err := stream.Send(&pbu.ListUsersResponse{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}
