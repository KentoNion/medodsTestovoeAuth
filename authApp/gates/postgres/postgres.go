package store

import (
	"context"
	"github.com/bool64/sqluct"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"medodsTestovoe/auth/pkg"
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

func (s Store) Save(ctx context.Context, token pkg.Refresh, userID string) error {
	query := "INSERT INTO tokens (user_id, token) VALUES ($1, $2) ON CONFLICT (token) DO UPDATE SET token = $2"
	_, err := s.db.ExecContext(ctx, query, userID, token)
	if err != nil {
		return errors.Wrap(err, "failed to save token")
	}
	return nil
}

func (s Store) Get(ctx context.Context, token pkg.Refresh) (bool, error) {
	query := "SELECT 1 FROM tokens where token = $1"
	rows, err := s.db.QueryContext(ctx, query, token)
	defer rows.Close()
	if err != nil {
		return false, errors.Wrap(err, "failed to query")
	}
	return rows.Next(), nil
}

func (s Store) Delete(ctx context.Context, token pkg.Refresh) error {
	query := "DELETE FROM tokens WHERE token = $1"
	_, err := s.db.ExecContext(ctx, query, token)
	if err != nil {
		return errors.Wrap(err, "failed to delete")
	}
	return nil
}
