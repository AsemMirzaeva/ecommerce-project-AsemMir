package storage

import (
	"database/sql"
	"product-service/storage/postgres"
	"product-service/storage/repo"
)

type IStorage interface {
	Product() repo.ProductStorage
}

type storagePg struct {
	db          *sql.DB
	productRepo repo.ProductStorage
}

func (s storagePg) Product() repo.ProductStorage {
	return s.productRepo
}

func NewStoragePg(db *sql.DB) *storagePg {
	return &storagePg{db, postgres.NewPostgresRepository(db)}
}
