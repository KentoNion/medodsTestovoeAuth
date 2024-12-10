package auth

import (
	"medodsTestovoe/auth/pkg"
	"time"
)

type Token struct {
	UserID string
	Secret string
	IP     string
}

type AuthTokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (t Token) MapToAcces(cl pkg.Clock) jwt.Claims {
	return jwt.MapClaims{
		"user_id": t.UserID,
		"secret":  t.Secret,
		"ip":      t.IP,
		"exp":     cl.Now().Add(time.Hour * 24).Unix(), //интерфейс получения времени + 24 часа
	}
}

func (t Token) MapToRefresh(cl pkg.Clock) jwt.Claims {
	return jwt.MapClaims{
		"user_id": t.UserID,
		"secret":  t.Secret,
		"ip":      t.IP,
		"exp":     cl.Now().AddDate(1, 0, 0), //интерфейс получения времени + 1 год
	}
}
