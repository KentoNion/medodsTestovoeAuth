package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"medodsTestovoe/auth/pkg"
)

type authStore interface {
	Save(ctx context.Context, token pkg.Refresh, userID string) error
	Get(ctx context.Context, token pkg.Refresh) (bool, error)
	Delete(ctx context.Context, token pkg.Refresh) error
}

type notifier interface {
	NotifyNewLogin(ctx context.Context, userID string, oldIP string, newIP string) error
}

type Service struct {
	secretKey string
	store     authStore
	notifier  notifier
	cl        pkg.Clock
}

func NewService(secretKey string, store authStore, notifier notifier, cl pkg.Clock) *Service {
	return &Service{
		secretKey: secretKey,
		store:     store,
		notifier:  notifier,
		cl:        cl,
	}
}

func (s *Service) Authorize(ctx context.Context, secret string, userID string, ip string) (Tokens AuthTokens, err error) {
	//JWT, SHA512, не храним
	token := Token{
		UserID: userID,
		Secret: secret,
		IP:     ip,
	}
	result := AuthTokens{
		Access:  "",
		Refresh: "",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.MapToRefresh(s.cl))
	if err != nil {
		return result, errors.Wrap(err, "failed to make access token")
	}
	var refresh pkg.Refresh
	refreshStr, err := refreshToken.SignedString([]byte(s.secretKey))
	refresh = pkg.Refresh(refreshStr)
	if err != nil {
		return result, errors.Wrap(err, "failed to make refresh token")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.MapToAccess(s.cl, refreshStr))
	var access pkg.Access
	accessStr, err := accessToken.SignedString([]byte(s.secretKey))
	access = pkg.Access(accessStr)
	err = s.store.Save(ctx, refresh, userID)
	if err != nil {
		return result, err
	}

	result = AuthTokens{
		Access:  access,
		Refresh: refresh,
	}
	return result, nil
}

func (s *Service) Refresh(ctx context.Context, refresh pkg.Refresh, ip string) (newTokens AuthTokens, err error) {
	//храним в базе в виде хеша
	result := AuthTokens{
		Access:  "",
		Refresh: "",
	}
	refreshStr := string(refresh)
	exists, err := s.store.Get(ctx, refresh)
	if err != nil {
		return result, errors.Wrap(err, "failed to check access token")
	}
	if !exists {
		return result, ErrRefreshTokenNotFound
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return result, errors.Wrap(err, "failed to parse refresh token")
	}
	if !token.Valid {
		return result, ErrWrongToken
	}
	refreshToken := Token{}
	if err := refreshToken.Fill(claims); err != nil {
		return result, errors.Wrapf(err, "failed to parse refresh token for user %s from %s", refreshToken.UserID, ip)
	}

	if ip != refreshToken.IP { // проверяем соответствует ли текущий ip старому ip при authorize
		err := s.notifier.NotifyNewLogin(ctx, refreshToken.UserID, ip, refreshToken.IP)
		if err != nil {
			return result, errors.Wrapf(err, "failed to notify new login for user %s from %s, old IP - %s", refreshToken.UserID, ip, refreshToken.IP)
		}
	}

	//генерируем новый access
	accessToken := Token{}
	if err := accessToken.Fill(claims); err != nil {
		return result, errors.Wrapf(err, "failed to parse access token for user %s from %s", refreshToken.UserID, ip)
	}
	if refreshToken.Secret != accessToken.Secret {
		return result, ErrWrongToken
	}

	newAccessJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, accessToken.MapToAccess(s.cl, refreshStr))
	newAccessStr, err := newAccessJWT.SignedString([]byte(s.secretKey))
	if err != nil {
		return result, errors.Wrapf(err, "failed to make access token for user %s from %s", refreshToken.UserID, ip)
	}
	newAccess := pkg.Access(newAccessStr)

	result = AuthTokens{
		Access:  newAccess,
		Refresh: refresh,
	}
	return result, nil
}

/*
Логика рефреша такая: По истечению аксес токена нужно произвести рефреш, мы считываем тот рефреш что даёт пользователь, сверяем его с тем что с бд, если там есть такой то:
1. Достаём старый ip из рефреш токена
2. Сравниваем его с переданным ip
3. Если ip не совпадают отправляем уведомление в моковый нотификатор
4. Генерируем новый access токен
5. Передаём новый аксес и старый рефреш если всё ок.
*/
