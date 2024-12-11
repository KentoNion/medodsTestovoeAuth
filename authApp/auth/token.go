package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"medodsTestovoe/auth/pkg"
	"time"
)

type Token struct {
	UserID string
	Secret string
	IP     string
}

func (t Token) MapToAccess(cl pkg.Clock, refresh string) jwt.Claims {
	return jwt.MapClaims{
		"user_id": t.UserID,
		"refresh": refresh,
		"ip":      t.IP,
		"exp":     cl.Now().Add(time.Hour * 24).Unix(), //интерфейс получения времени + 24 часа
	}
}

func (t Token) MapToRefresh(cl pkg.Clock) jwt.Claims {
	return jwt.MapClaims{
		"user_id": t.UserID,
		"secret":  t.Secret,
		"ip":      t.IP,
		"exp":     cl.Now().AddDate(0, 1, 0), //интерфейс получения времени + 1 месяц
	}
}

func (t Token) Fill(claims jwt.MapClaims) error {
	var ok bool
	t.IP, ok = claims["ip"].(string)
	if !ok {
		return errors.New("failed to parse ip")
	}
	t.Secret, ok = claims["secret"].(string)
	if !ok {
		return errors.New("failed to parse secret")
	}
	t.UserID, ok = claims["user_id"].(string)
	if !ok {
		return errors.New("failed to parse user_id")
	}
	return nil
}

type AuthTokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

var ErrWrongToken = errors.New("Wrong token")

var ErrRefreshTokenNotFound = errors.New("Refresh token not found")
