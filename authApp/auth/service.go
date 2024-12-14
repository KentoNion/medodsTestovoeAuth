package auth

import (
	"context"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"medodsTestovoe/auth/pkg"
)

type authStore interface {
	Save(ctx context.Context, token pkg.Refresh, userID string, ip string) error
	Get(ctx context.Context, userID string, token pkg.Refresh) (bool, string, error)
	Delete(ctx context.Context, userID string) error
	CheckUserExist(ctx context.Context, userID string) (bool, error)
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
	//заполняем токен
	token := Token{
		UserID: userID,
		Secret: secret,
		IP:     ip,
	}
	result := AuthTokens{
		Access:  "",
		Refresh: "",
	}
	//Проверка сущетсвует ли такой пользователь в бд
	exist, err := s.store.CheckUserExist(ctx, userID)
	if err != nil {
		return Tokens, errors.Wrap(err, "error checking user existence")
	}
	if exist {
		return result, ErrUserAlreadyExists
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token.MapToRefresh(s.cl))
	if err != nil {
		return result, errors.Wrap(err, "failed to make access token")
	}
	var refresh pkg.Refresh
	refreshStr, err := refreshToken.SignedString([]byte(s.secretKey))
	refreshStrBase64 := base64.StdEncoding.EncodeToString([]byte(refreshStr))
	if len(refreshStrBase64) > 64 {
		refreshStrBase64 = refreshStrBase64[len(refreshStrBase64)-66 : len(refreshStrBase64)-2] //bcrypt требует длинну токена не больше чем 64 байта, я не придумал как ещё можно влезть в это требование кроме как обрезать токен по изменяемой части
	}
	refresh = pkg.Refresh(refreshStrBase64)
	if err != nil {
		return result, errors.Wrap(err, "failed to make refresh token")
	}
	//генерируем аксес токен
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.MapToAccess(s.cl, refreshStr))
	var access pkg.Access
	accessStr, err := accessToken.SignedString([]byte(s.secretKey))
	access = pkg.Access(accessStr)
	err = s.store.Save(ctx, refresh, userID, ip) //сохраняем рефреш, userID и ip в бд
	if err != nil {
		return result, err
	}
	//формируем тело ответа
	result = AuthTokens{
		Access:  access,
		Refresh: refresh,
	}
	return result, nil
}

func (s *Service) Refresh(ctx context.Context, userID string, refresh pkg.Refresh, ip string) (newTokens AuthTokens, err error) {
	//храним в базе в виде хеша
	result := AuthTokens{ //регестрируем пустой результат
		Access:  "",
		Refresh: "",
	}
	//проверка существует ли такой токен в БД а так же извлечение значения старого ip адреса
	exists, oldIP, err := s.store.Get(ctx, userID, refresh)
	if err != nil {
		return result, err
	}
	if !exists {
		return result, ErrRefreshTokenNotFound
	}

	if ip != oldIP { // проверяем соответствует ли текущий ip старому ip при authorize
		err := s.notifier.NotifyNewLogin(ctx, userID, ip, oldIP)
		if err != nil {
			return result, errors.Wrapf(err, "failed to notify new login for user %s from %s, old IP - %s", userID, ip, oldIP)
		}
	}
	/*
		!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		из задания: Refresh токен должен быть защищён от повторного использования
		Я так понял нам нужно сгенерировать новый Refresh
		!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	*/
	//генерируем новые access и refresh если прошли все проверки выше
	err = s.store.Delete(ctx, userID) //удаляем старую запись
	if err != nil {
		return result, err
	}
	return s.Authorize(ctx, uuid.New().String(), userID, ip) //создаём новую запись с тем же userID и возвращаем новые токены
}
