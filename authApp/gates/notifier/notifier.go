package notify

import (
	"context"
	"log"
)

// интерфейс для уведомлений
type Notifier interface {
	NotifyNewLogin(ctx context.Context, userID string) error
}

// моковая реализация интерфейса нотифаера
type mockNotifier struct{}

// NotifyNewLogin - метод мокового уведомления о новом логине.
func (m *mockNotifier) NotifyNewLogin(ctx context.Context, userID string) error {
	log.Printf("Mock notification: new login for user %s", userID)
	return nil
}

// инициализация мок нотифаера
func InitNotifier() Notifier {
	return &mockNotifier{}
}
