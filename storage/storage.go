package storage

import (
	"github.com/Asliddin3/customer-servis/storage/postgres"
	"github.com/Asliddin3/customer-servis/storage/repo"
	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Customer() repo.CustomerStorageI
	Post() repo.PostStorageI
}

type storagePg struct {
	db           *sqlx.DB
	customerRepo repo.CustomerStorageI
	postRepo     repo.PostStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:           db,
		customerRepo: postgres.NewCustomerRepo(db),
	}
}
func (s storagePg) Customer() repo.CustomerStorageI {
	return s.customerRepo
}

func (s storagePg) Post() repo.PostStorageI {
	return s.postRepo
}
