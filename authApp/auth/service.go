package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"medodsTestovoe/auth/pkg"
)

type AuthStore interface {
	Save(ctx context.Context, token pkg.Hash, userID string, ip string) error
	Get(ctx context.Context, userID string) (pkg.Hash, string, error)
	Delete(ctx context.Context, userID string) error
}

type notifier interface {
	NotifyNewLogin(ctx context.Context, userID string, oldIP string, newIP string) error
}

type Service struct {
	secretKey string
	store     AuthStore
	notifier  notifier
	cl        pkg.Clock
}

func NewService(secretKey string, store AuthStore, notifier notifier, cl pkg.Clock) *Service {
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
		Secret: secret, //Здесь должен быть uuid.New().String(), но для того что бы тест Authorize правильно работал я вынес secret за пределы service и он задаётся при вызове функции authorize
		IP:     ip,
	}
	result := AuthTokens{
		Access:  "",
		Refresh: "",
	}
	//Проверка сущетсвует ли такой пользователь в бд
	_, _, err = s.store.Get(ctx, userID)
	if err == nil { //nil будет если такой уже есть в бд
		return result, ErrGUIDAlreadyExists
	} else if err == sql.ErrNoRows { //если ошибка "пользователя нет"
		err = nil //значит никакой ошибки нет, едем дальше
	}
	if err != nil { //если ошибка другая, значит что-то не так
		return Tokens, errors.Wrap(err, "error checking user existence")
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
	hash, err := pkg.HashToken(refresh) //хешируем в bcrypt
	if err != nil {
		return result, err
	}
	err = s.store.Save(ctx, hash, userID, ip) //сохраняем рефреш, userID и ip в бд
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
	storedHash, oldIP, err := s.store.Get(ctx, userID)
	if err == sql.ErrNoRows {
		return result, ErrGUIDNotFound
	}
	if err != nil {
		return result, err
	}

	if ip != oldIP { // проверяем соответствует ли текущий ip старому ip при authorize
		err := s.notifier.NotifyNewLogin(ctx, userID, ip, oldIP)
		if err != nil {
			return result, errors.Wrapf(err, "failed to notify new login for user %s from %s, old IP - %s", userID, ip, oldIP)
		}
	}

	// Сравнение токена с сохранённым хэшем
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(refresh))
	if err != nil {
		return result, errors.New("invalid token, not found")
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
