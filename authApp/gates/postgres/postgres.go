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

func (s Store) Save(ctx context.Context, token pkg.Hash, userID string, ip string) error {
	query := "INSERT INTO tokens (user_id, token, ip) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token = $2, ip = $3"
	_, err := s.db.ExecContext(ctx, query, userID, token, ip)

	if err != nil {
		return errors.Wrap(err, "failed to save token")
	}
	return nil
}

func (s Store) Get(ctx context.Context, userID string) (pkg.Hash, string, error) {
	var storedHash pkg.Hash
	var storedIP string
	query := "SELECT token, ip FROM tokens WHERE user_id = $1"
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&storedHash, &storedIP)
	if err != nil { //если ошибка другая, значит какая-то лажа
		return "", "", err
	}
	return storedHash, storedIP, nil
}

func (s Store) Delete(ctx context.Context, userID string) error {
	query := "DELETE FROM tokens WHERE user_id = $1"
	_, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return errors.Wrap(err, "failed to delete")
	}
	return nil
}
