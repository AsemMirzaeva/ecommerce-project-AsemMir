package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	proto "product-service/proto"
)

type ProductRepo struct {
	DB *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (s *ProductRepo) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	id := uuid.New().String()
	query := "INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)"
	_, err := s.DB.ExecContext(ctx, query, id, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.RegisterResponse{
		Id:       id,
		Username: req.Username,
		Email:    req.Email,
	}, nil
}

func (s *ProductRepo) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	query := "SELECT id, username, email FROM users WHERE id=$1"
	row := s.DB.QueryRowContext(ctx, query, req.Id)

	var user proto.GetUserResponse
	err := row.Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *ProductRepo) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	query := "UPDATE users SET username=$1, email=$2 WHERE id=$3"
	_, err := s.DB.ExecContext(ctx, query, req.Username, req.Email, req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateUserResponse{
		Id:       req.Id,
		Username: req.Username,
		Email:    req.Email,
	}, nil
}

func (s *ProductRepo) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	query := "DELETE FROM users WHERE id=$1"
	_, err := s.DB.ExecContext(ctx, query, req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteUserResponse{Status: "Deleted"}, nil
}

func (s *ProductRepo) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	query := "SELECT id FROM users WHERE email=$1 AND password=$2"
	row := s.DB.QueryRowContext(ctx, query, req.Email, req.Password)

	var id string
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Generate a token (this is a placeholder, implement your token generation logic)
	token := "generated_token_based_on_user_id"

	return &proto.LoginResponse{Token: token}, nil
}
