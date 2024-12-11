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

	svc := service{
		secretKey: "test_key",
		store:     store,
		notifier:  notifier,
		cl:        pkg.StubClock{time.Date(2024, time.September, 24, 0, 0, 0, 0, time.UTC)},
	}

	expectRefresh := gomock.Any()
	store.EXPECT().Save(gomock.Any(), "test_user", expectRefresh).Return(nil)

	ctx := context.Background()
	tok, err := svc.Authorize(ctx, "password", "test_user", "123.123.123.123")
	require.NoError(t, err)
	require.Equal(t, "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjcyMjI0MDAsImlwIjoiMTIzLjEyMy4xMjMuMTIzIiwicmVmcmVzaCI6ImV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpsZUhBaU9pSXlNREkwTFRFd0xUSTBWREF3T2pBd09qQXdXaUlzSW1sd0lqb2lNVEl6TGpFeU15NHhNak11TVRJeklpd2ljMlZqY21WMElqb2ljR0Z6YzNkdmNtUWlMQ0oxYzJWeVgybGtJam9pZEdWemRGOTFjMlZ5SW4wLkR2U0l1cjctSjZZQThrUWV2Rkp5dXJiNjVKNFNVOUZMYTlGeG9wV19CMVZEazFRTWQ0OEx6VUp0N0V5dTdqRUlrbGU2T25rLWdQT0lDeTlfZ3U3SXVnIiwidXNlcl9pZCI6InRlc3RfdXNlciJ9.OBtA1D50OchHr-kb9Tb_587kPJMFZ3M-0BQFAsbVCYVvwK0URQfZ-ZmQxDn8ITeH1DaMNTOqUWlu5dZkxSnbXQ", tok.Access)
	require.Equal(t, "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI0LTEwLTI0VDAwOjAwOjAwWiIsImlwIjoiMTIzLjEyMy4xMjMuMTIzIiwic2VjcmV0IjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoidGVzdF91c2VyIn0.DvSIur7-J6YA8kQevFJyurb65J4SU9FLa9FxopW_B1VDk1QMd48LzUJt7Eyu7jEIkle6Onk-gPOICy9_gu7Iug", tok.Refresh)
}
