package db

import (
	"database/sql"

	_ "github.com/golang/mock/mockgen/model"
)

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
