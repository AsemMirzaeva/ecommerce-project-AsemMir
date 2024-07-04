package service

import (
	"context"
	"database/sql"
	"errors"
	"product-service/models"
	l "product-service/pkg/logger"
	pbp "product-service/protos/product-service"
	"product-service/storage"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ProductService struct {
	pbp.UnimplementedProductServiceServer
	storage storage.IStorage
	logger  l.Logger
}

func NewProductService(db *sql.DB, log l.Logger) *ProductService {
	return &ProductService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *ProductService) AddProduct(ctx context.Context, req *pbp.AddProductRequest) (*pbp.AddProductResponse, error) {
	id := uuid.New().String()
	product, err := s.storage.Product().AddProduct(&models.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		return nil, err
	}
	return &pbp.AddProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pbp.GetProductRequest) (*pbp.GetProductResponse, error) {
	product, err := s.storage.Product().GetProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return &pbp.GetProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *pbp.UpdateProductRequest) (*pbp.UpdateProductResponse, error) {
	product, err := s.storage.Product().UpdateProduct(&models.Product{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		return nil, err

	}
	return &pbp.UpdateProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil

}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pbp.DeleteProductRequest) (*pbp.DeleteProductResponse, error) {
	if err := s.storage.Product().DeleteProduct(req.Id); err != nil {
		return nil, err
	}
	return &pbp.DeleteProductResponse{Message: "Product deleted"}, nil
}

func (s *ProductService) ListProducts(req *pbp.ListProductsRequest, stream pbp.ProductService_ListProductsServer) error {
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

	products, err := s.storage.Product().ListProducts(limit, page)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	for _, product := range products {
		err := stream.Send(&pbp.ListProductsResponse{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}
