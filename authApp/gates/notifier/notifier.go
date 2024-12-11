package notify

import (
	"context"
	"log"
)

// интерфейс для уведомлений
type Notifier interface {
	NotifyNewLogin(ctx context.Context, userID string, oldIP string, newIP string) error
}

// моковая реализация интерфейса нотифаера
type mockNotifier struct{}

// NotifyNewLogin - метод мокового уведомления о новом логине.
func (m *mockNotifier) NotifyNewLogin(ctx context.Context, userID string, oldIP string, newIP string) error {
	log.Printf("Mock notification: new login for user %s, new IP - %s, origin IP - %s", userID, newIP, oldIP)
	return nil
}

// инициализация мок нотифаера
func InitNotifier() Notifier {
	return &mockNotifier{}
}
