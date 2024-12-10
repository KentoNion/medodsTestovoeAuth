package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"medodsTestovoe/auth/pkg"
)

type authStore interface {
	Save(ctx context.Context, token string, userID string) error
	Get(ctx context.Context, token string) (bool, error)
	delete(ctx context.Context, token string) error
}

type notifier interface {
	NotifyNewLogin(ctx context.Context, userID string) error
}

type service struct {
	secretKey string
	store     authStore
	notifier  notifier
	cl        pkg.Clock
}

func NewService(secretKey string, store authStore, notifier notifier, cl pkg.Clock) *service {
	return &service{
		secretKey: secretKey,
		store:     store,
		notifier:  notifier,
		cl:        cl,
	}
}

func (s *service) Authorize(ctx context.Context, secret string, userID string, ip string) (Tokens AuthTokens, err error) {
	//JWT, SHA512, не храним
	token := Token{
		UserID: userID,
		Secret: secret, //uuid.New().String() <-----------------------------------------------------------
		IP:     ip,
	}
	result := AuthTokens{
		Access:  "",
		Refresh: "",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.MapToAcces(s.cl))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.MapToRefresh(s.cl))
	access, err := accessToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return result, errors.Wrap(err, "failed to make access token")
	}

	refresh, err := refreshToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return result, errors.Wrap(err, "failed to make refresh token")
	}

	err = s.store.Save(ctx, userID, refresh)
	if err != nil {
		return result, err
	}

	result = AuthTokens{
		Access:  access,
		Refresh: refresh,
	}
	return result, nil
}

func (s *service) Refresh(ctx context.Context, refresh string) (newTokens AuthTokens, err error) {
	//храним в базе в виде хеша

}
