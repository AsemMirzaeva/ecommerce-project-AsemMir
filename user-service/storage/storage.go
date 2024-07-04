package storage

import (
	"database/sql"
	"user-service/storage/postgres"
	"user-service/storage/repo"
)

type IStorage interface {
	User() repo.UserStorage
}

type storagePg struct {
	db       *sql.DB
	userRepo repo.UserStorage
}

func (s storagePg) User() repo.UserStorage {
	return s.userRepo
}

func NewStoragePg(db *sql.DB) *storagePg {
	return &storagePg{db, postgres.NewPostgresRepository(db)}
}
