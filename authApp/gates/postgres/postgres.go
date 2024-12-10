package postgres

import (
	"context"
	"github.com/bool64/sqluct"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
	sm sqluct.Mapper
}

func NewDB(db *sqlx.DB) *Store {
	return &Store{
		db: db,
		sm: sqluct.Mapper{Dialect: sqluct.DialectPostgres},
	}
}

func (s Store) Save(ctx context.Context) error {

}

func (s Store) Get(ctx context.Context) (string, error) {

}

func (s Store) Delete(ctx context.Context) error {

}
