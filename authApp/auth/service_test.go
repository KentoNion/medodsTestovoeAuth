package auth

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	mock "medodsTestovoe/auth/mock"
	"medodsTestovoe/auth/pkg"
	"testing"
	"time"
)

func TestAuthorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockAuthStore(ctrl)
	notifier := mock.NewMocknotifier(ctrl)

	svc := Service{
		secretKey: "test_key",
		store:     store,
		notifier:  notifier,
		cl:        pkg.StubClock{time.Date(2024, time.September, 24, 0, 0, 0, 0, time.UTC)},
	}

	store.EXPECT().Get(gomock.Any(), "testUser").Return(pkg.Hash(""), "", sql.ErrNoRows) //говорим что такой записи нет
	store.EXPECT().Save(gomock.Any(), gomock.Any(), "testUser", "123.123.123.123").Return(nil)

	ctx := context.Background()
	tok, err := svc.Authorize(ctx, "password", "testUser", "123.123.123.123")
	require.NoError(t, err)
	require.Equal(t, pkg.Access("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjcxMzYwMTAsImlwIjoiMTIzLjEyMy4xMjMuMTIzIiwicmVmcmVzaCI6ImV5SmhiR2NpT2lKSVV6STFOaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpsZUhBaU9qRTNNamszTWpnd01EQXNJbWx3SWpvaU1USXpMakV5TXk0eE1qTXVNVEl6SWl3aWMyVmpjbVYwSWpvaWNHRnpjM2R2Y21RaUxDSjFjMlZ5WDJsa0lqb2lkR1Z6ZEZWelpYSWlmUS5QQUphWmo4NmEtRF9WMTdQM05oSU5MOVAwTlV6MG9raTdKU2NHZ0F2V2drIiwidXNlcl9pZCI6InRlc3RVc2VyIn0.H1BIo6XQLQCfa2LK2XEq1aisrpE_T0gqV4WDaFKLMfcc2_C0-r7XgkZTgSO0tgKQWywDzBnsODLpdGUAYu4rqA"), tok.Access)
	require.Equal(t, pkg.Refresh("hJaWZRLlBBSmFaajg2YS1EX1YxN1AzTmhJTkw5UDBOVXowb2tpN0pTY0dnQXZXZ2"), tok.Refresh)
}

func TestRefresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockAuthStore(ctrl)
	notifier := mock.NewMocknotifier(ctrl)

	svc := Service{
		secretKey: "test_key",
		store:     store,
		notifier:  notifier,
		cl:        pkg.StubClock{time.Date(2024, time.September, 24, 0, 0, 0, 0, time.UTC)},
	}
	//проверка правильности выдаваемой ошибки если такого рефреш токена не сущетсвует
	store.EXPECT().Get(gomock.Any(), "testUser").Return(pkg.Hash(""), "", sql.ErrNoRows) //типа у нас нет такого пользователя

	ctx := context.Background()
	_, err := svc.Refresh(ctx, "testUser", "123", "123.123.123.123") //делаем запрос refresh с ожидаемым ответом что нет такого пользователя

	require.Equal(t, err, ErrGUIDNotFound)

	hash, err := pkg.HashToken(pkg.Refresh("123"))
	require.NoError(t, err)

	store.EXPECT().Get(gomock.Any(), "testUser").Return(hash, "123.123.123.123", nil)
	store.EXPECT().Delete(gomock.Any(), "testUser").Return(nil)
	store.EXPECT().Get(gomock.Any(), "testUser").Return(pkg.Hash(""), "", sql.ErrNoRows) //говорим что такой записи нет
	store.EXPECT().Save(gomock.Any(), gomock.Any(), "testUser", "123.123.123.123").Return(nil)

	tokens, err := svc.Refresh(ctx, "testUser", "123", "123.123.123.123")
	require.NoError(t, err)
	require.NotEmpty(t, tokens)

	//проверка на уведомление если мы подсовываем другой ip
	store.EXPECT().Get(gomock.Any(), "testUser").Return(hash, "123.123.123.124", nil)
	notifier.EXPECT().NotifyNewLogin(gomock.Any(), "testUser", "123.123.123.123", "123.123.123.124") //проверяем что сервис пытается отправить уведомление что ip изменился
	store.EXPECT().Delete(gomock.Any(), "testUser").Return(nil)
	store.EXPECT().Get(gomock.Any(), "testUser").Return(pkg.Hash(""), "", sql.ErrNoRows)
	store.EXPECT().Save(gomock.Any(), gomock.Any(), "testUser", "123.123.123.123").Return(nil)

	tokens, err = svc.Refresh(ctx, "testUser", "123", "123.123.123.123")
	require.NoError(t, err)
	require.NotEmpty(t, tokens)
}
