package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"medodsTestovoe/auth/pkg"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStore struct{}

func (m *mockStore) Save(ctx context.Context, token pkg.Refresh, userID string, ip string) error {
	return nil
}

func (m *mockStore) Get(ctx context.Context, userID string, token pkg.Refresh) (bool, string, error) {
	if userID == "testUser" {
		return true, "255.255.255.255", nil
	}
	return false, "", errors.New("not found")
}

func (m *mockStore) Delete(ctx context.Context, userID string) error {
	return nil
}

func (m *mockStore) CheckUserExist(ctx context.Context, userID string) (bool, error) {
	return userID == "testUser", nil
}

type mockNotifier struct{}

func (m *mockNotifier) NotifyNewLogin(ctx context.Context, userID string, oldIP string, newIP string) error {
	return nil
}

func TestServer_LoginHandler(t *testing.T) {
	//заполняем моковую структуру сервера
	logger := zap.NewNop()
	mockDb := &mockStore{}
	mockNotifier := &mockNotifier{}
	r := chi.NewRouter()
	_ = NewServer(mockDb, r, logger, mockNotifier)

	//стучимся с неправильным методом
	req := httptest.NewRequest(http.MethodGet, "/login?user_id=123testUser123&secret=789", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	//Получаем 405 тк такого метода нет
	require.Equal(t, http.StatusMethodNotAllowed, rec.Code, "expected status 405")

	//стучимся просто в логин без всего
	req = httptest.NewRequest(http.MethodPost, "/login", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	//должны получить статусАнотхарайзд
	require.Equal(t, http.StatusUnauthorized, rec.Code, "expected status 401")

	//Стучимся с правильным запросом
	req = httptest.NewRequest(http.MethodPost, "/login?user_id=123testUser123&secret=789", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	//Должны получить 200 тк запрос корректный
	require.Equal(t, http.StatusOK, rec.Code, "expected status 200")
}

func TestServer_RefreshHandler(t *testing.T) {
	//заполняем моковую структуру сервера
	logger := zap.NewNop()
	mockDb := &mockStore{}
	mockNotifier := &mockNotifier{}
	r := chi.NewRouter()
	_ = NewServer(mockDb, r, logger, mockNotifier)

	//стучимся с неправильным методом
	req := httptest.NewRequest(http.MethodGet, "/refresh?refresh_token=VpZlEuSkZMVFQwRnVQblJkMnI4SGxxNU1FdVZtVHhmWFEydEF0em9odzNoaVhCTQ&user_id=1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	//Получаем 405 тк такого метода нет
	require.Equal(t, http.StatusMethodNotAllowed, rec.Code, "expected status 405")

	//стучимся просто рефреш без всего
	req = httptest.NewRequest(http.MethodPost, "/refresh", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	//ловим StatusUnauthorized
	require.Equal(t, http.StatusUnauthorized, rec.Code, "expected status 401")
	//отправляем нормальный запрос
	req = httptest.NewRequest(http.MethodPost, "/refresh?refresh_token=VpZlEuSkZMVFQwRnVQblJkMnI4SGxxNU1FdVZtVHhmWFEydEF0em9odzNoaVhCTQ&user_id=1", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	//но всё ещё получаем ошибку тк у нас пустая тестовая бд и неоткуда взяться refresh и никак не провести операцию
	require.Equal(t, http.StatusInternalServerError, rec.Code, "expected status 500")
}
