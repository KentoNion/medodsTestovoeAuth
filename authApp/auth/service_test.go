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

	expectRefresh := gomock.Any()
	store.EXPECT().Save(gomock.Any(), expectRefresh, "test_user").Return(nil)

	ctx := context.Background()
	tok, err := svc.Authorize(ctx, "password", "test_user", "123.123.123.123")
	require.NoError(t, err)
	require.Equal(t, pkg.Access("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjcxMzYwMTAsImlwIjoiMTIzLjEyMy4xMjMuMTIzIiwicmVmcmVzaCI6ImV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpsZUhBaU9qRTNNamszTWpnd01EQXNJbWx3SWpvaU1USXpMakV5TXk0eE1qTXVNVEl6SWl3aWMyVmpjbVYwSWpvaWNHRnpjM2R2Y21RaUxDSjFjMlZ5WDJsa0lqb2lkR1Z6ZEY5MWMyVnlJbjAuc3pReUxEcHo2Vm0yLVFMOU5BMmVuQkE1dWNhT0JvTjBJcmJvTV9SUERqQllyWTJ3MC0zQmhFTC01aXVQeGtBeVBzczhJUWpQVmJhdjQzdldXdE5SQmciLCJ1c2VyX2lkIjoidGVzdF91c2VyIn0.8AqHN2W5zka-oBOyqg52xrGd4qp1NF50pWdrUva7Wcc9hTuLKYfDKed7tbU5mlfBT1nXEJHrj2znoyJByusekw"), tok.Access)
	require.Equal(t, pkg.Refresh("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk3MjgwMDAsImlwIjoiMTIzLjEyMy4xMjMuMTIzIiwic2VjcmV0IjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoidGVzdF91c2VyIn0.szQyLDpz6Vm2-QL9NA2enBA5ucaOBoN0IrboM_RPDjBYrY2w0-3BhEL-5iuPxkAyPss8IQjPVbav43vWWtNRBg"), tok.Refresh)
}
