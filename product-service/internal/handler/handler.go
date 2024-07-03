package handler

import (
    "context"

    "product-service/internal/config"
    "product-service/internal/proto"

    "github.com/google/uuid"
)

type Handler struct {
    proto.UnimplementedProductServiceServer
    cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
    return &Handler{cfg: cfg}
}

func (h *Handler) AddProduct(ctx context.Context, req *proto.AddProductRequest) (*proto.AddProductResponse, error) {
    productID := uuid.New().String()
    // Assuming you have a product repository to save the product
    // err := h.productRepo.AddProduct(productID, req.Name, req.Price)
    if err != nil {
        log.Logger.Printf("error adding product: %v", err)
        return nil, err
    }

    return &proto.AddProductResponse{Id: productID}, nil
}

func (h *Handler) ListProducts(req *proto.ListProductsRequest, stream proto.ProductService_ListProductsServer) error {
    // Assuming you have a product repository to list all products
    // products, err := h.productRepo.ListProducts()
    if err != nil {
        log.Logger.Printf("error listing products: %v", err)
        return err
    }

    for _, product := range products {
        if err := stream.Send(&proto.ListProductsResponse{
            Products: []*proto.Product{
                {
                    Id:    product.ID,
                    Name:  product.Name,
                    Price: product.Price,
                },
            },
        }); err != nil {
            log.Logger.Printf("error sending product: %v", err)
            return err
        }
    }

    return nil
}
