package auth

import (
	"time"
)

type Token struct {
	UserID string
	Secret string
	IP     string
}

func (t Token) MapToAcces(cl pkg.Clock) jwt.Claims {
	return jwt.MapClaims{
		"user_id": t.UserID,
		"secret":  t.Secret,
		"ip":      t.IP,
		"exp":     cl.Now().Add(time.Hour * 24).Unix(), //интерфейс получения времени + 24 часа
	}
}
