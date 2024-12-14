package store

import (
	"context"
	"fmt"
	"github.com/bool64/sqluct"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"medodsTestovoe/auth/pkg"
)

func hashToken(token string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(hashedBytes), nil
}

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

func (s Store) Save(ctx context.Context, token pkg.Refresh, userID string, ip string) error {
	hash, err := hashToken(string(token))
	if err != nil {
		return err
	}
	query := "INSERT INTO tokens (user_id, token, ip) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token = $2, ip = $3"
	_, err = s.db.ExecContext(ctx, query, userID, hash, ip)
	if err != nil {
		return errors.Wrap(err, "failed to save token")
	}
	return nil
}

func (s Store) Get(ctx context.Context, userID string, token pkg.Refresh) (bool, string, error) {
	var storedHash string
	var storedIP string
	query := "SELECT token, ip FROM tokens WHERE user_id = $1"
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&storedHash, &storedIP)
	if err != nil {
		return false, "", errors.Wrap(err, "failed to query")
	}

	// Сравнение токена с сохранённым хэшем
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(token))
	if err != nil {
		return false, "", errors.New("invalid token, not found")
	}
	return true, storedIP, nil
}

func (s Store) Delete(ctx context.Context, userID string) error {
	query := "DELETE FROM tokens WHERE user_id = $1"
	_, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return errors.Wrap(err, "failed to delete")
	}
	return nil
}

func (s Store) CheckUserExist(ctx context.Context, userID string) (bool, error) {
	query := "SELECT 1 FROM tokens where user_id = $1"
	rows, err := s.db.QueryContext(ctx, query, userID)
	defer rows.Close()
	if err != nil {
		return false, errors.Wrap(err, "failed to query")
	}
	return rows.Next(), nil
}
