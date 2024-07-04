package storage

import (
	"database/sql"
	"order-service/storage/postgres"
	"order-service/storage/repo"
)

type IStorage interface {
	Order() repo.OrderStorage
}

type storagePg struct {
	db        *sql.DB
	OrderRepo repo.OrderStorage
}

func (s storagePg) Order() repo.OrderStorage {
	return s.OrderRepo
}

func NewStoragePg(db *sql.DB) *storagePg {
	return &storagePg{db, postgres.NewPostgresRepository(db)}
}
