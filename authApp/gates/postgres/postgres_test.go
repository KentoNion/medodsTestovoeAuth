package store

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	"github.com/stretchr/testify/require"
	"testing"
)

// в этом тесте я стучусь в продовую дб, он может сломаться если кто-то зарегает пользователя "testUser"
func TestInsertGetCheckExistDelete(t *testing.T) {
	//Подключаемся к дб
	conn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=medods_auth host=localhost sslmode=disable")
	if err != nil {
		require.NoError(t, err)
	}
	db := NewDB(conn)
	ctx := context.Background()
	//сохраняем тестового пользователя
	err = db.Save(ctx, "testToken", "testUser", "255.255.255.255")
	require.NoError(t, err) // проверяем отсутсвие ошибки

	got, ip, err := db.Get(ctx, "testUser2", "testToken") //получаем данные по несуществующему юзеру
	require.Error(t, err)                                 //получаем ожидаемую ошибку
	got, ip, err = db.Get(ctx, "testUser", "testToken")   ////получаем данные по userID "testUser" который мы только что записали
	require.Error(t, err)                                 //проверяем что ошибки нет
	require.Equal(t, true, got)                           //пользователь есть
	require.Equal(t, ip, "255.255.255.255")               //сверяем ip

	exist, err := db.CheckUserExist(ctx, "testUser") //проврека что пользователь есть
	require.NoError(t, err)                          //ошибки нет
	require.True(t, exist)                           //пользователь есть

	err = db.Delete(ctx, "testUser") //удаляем пользователя
	require.NoError(t, err)          //ошибки нет

	got, ip, err = db.Get(ctx, "testUser", "testToken") //ищем пользователя которого только что удалили
	//ловим ошибки и что его нет, двумя способами get и checkExist
	require.Error(t, err)
	require.Equal(t, false, got)
	exist, err = db.CheckUserExist(ctx, "testUser")
	require.NoError(t, err)
	require.False(t, exist)
}
