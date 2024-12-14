package auth

import (
	"context"
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

	store := mock.NewMockauthStore(ctrl)
	notifier := mock.NewMocknotifier(ctrl)

	svc := Service{
		secretKey: "test_key",
		store:     store,
		notifier:  notifier,
		cl:        pkg.StubClock{time.Date(2024, time.September, 24, 0, 0, 0, 0, time.UTC)},
	}

	store.EXPECT().Save(gomock.Any(), gomock.Any(), "test_user", "123.123.123.123").Return(nil)
	store.EXPECT().CheckUserExist(gomock.Any(), gomock.Any()).Return(false, nil)

	ctx := context.Background()
	tok, err := svc.Authorize(ctx, "password", "test_user", "123.123.123.123")
	require.NoError(t, err)
	require.Equal(t, pkg.Access("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjcxMzYwMTAsImlwIjoiMTIzLjEyMy4xMjMuMTIzIiwicmVmcmVzaCI6ImV5SmhiR2NpT2lKSVV6STFOaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpsZUhBaU9qRTNNamszTWpnd01EQXNJbWx3SWpvaU1USXpMakV5TXk0eE1qTXVNVEl6SWl3aWMyVmpjbVYwSWpvaWNHRnpjM2R2Y21RaUxDSjFjMlZ5WDJsa0lqb2lkR1Z6ZEY5MWMyVnlJbjAuM2JYcVZhNFcwdVlmYk9JYzBvTm9XdkFuVHMwamJCWU9Jd0hZYWo5NGhIbyIsInVzZXJfaWQiOiJ0ZXN0X3VzZXIifQ.CRbUkxRi7tA_Sos8CpvEXu53OE1qYE0FH-4M3uEbpH-iOruMk8MJtTz3Mh-9sesdtCvzW8DtHBa3Kpl8OCqfzA"), tok.Access)
	require.Equal(t, pkg.Refresh("JWeUluMC4zYlhxVmE0VzB1WWZiT0ljMG9Ob1d2QW5UczBqYkJZT0l3SFlhajk0aE"), tok.Refresh)
}

func TestRefresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockauthStore(ctrl)
	notifier := mock.NewMocknotifier(ctrl)

	svc := Service{
		secretKey: "test_key",
		store:     store,
		notifier:  notifier,
		cl:        pkg.StubClock{time.Date(2024, time.September, 24, 0, 0, 0, 0, time.UTC)},
	}
	//проверка правильности выдаваемой ошибки если такого рефреш токена не сущетсвует
	store.EXPECT().Get(gomock.Any(), "testUser", pkg.Refresh("123")).Return(false, "123.123.123.123", nil)

	ctx := context.Background()
	_, err := svc.Refresh(ctx, "testUser", "123", "123.123.123.123")

	require.Equal(t, err, ErrRefreshTokenNotFound)

	store.EXPECT().Get(gomock.Any(), "testUser", pkg.Refresh("123")).Return(true, "123.123.123.123", nil)
	store.EXPECT().Delete(gomock.Any(), "testUser").Return(nil)
	store.EXPECT().CheckUserExist(gomock.Any(), gomock.Any()).Return(false, nil)
	store.EXPECT().Save(gomock.Any(), gomock.Any(), "testUser", "123.123.123.123").Return(nil)

	tokens, err := svc.Refresh(ctx, "testUser", "123", "123.123.123.123")
	require.NoError(t, err)
	require.NotEmpty(t, tokens)
}
