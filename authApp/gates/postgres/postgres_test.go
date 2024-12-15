package store

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	"github.com/stretchr/testify/require"
	"medodsTestovoe/auth/pkg"
	"testing"
)

// в этом тесте я стучусь в продовую дб, он может сломаться если кто-то зарегает пользователя "testUser"
func TestInsertGetDelete(t *testing.T) {
	//Подключаемся к дб
	conn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=medods_auth host=localhost sslmode=disable")
	if err != nil {
		require.NoError(t, err)
	}
	db := NewDB(conn)
	ctx := context.Background()
	//ищем того чего нет, надеемся получить ошибку что этого нет
	_, _, err = db.Get(ctx, "testUser")
	require.Equal(t, sql.ErrNoRows, err)

	//сохраняем тестового пользователя
	err = db.Save(ctx, "testHash", "testUser", "255.255.255.255")
	require.NoError(t, err) // проверяем отсутсвие ошибки

	hash, ip, err := db.Get(ctx, "testUser")     //получаем данные по userID "testUser" который мы только что записали
	require.NoError(t, err)                      //проверяем что ошибки нет
	require.Equal(t, ip, "255.255.255.255")      //сверяем ip
	require.Equal(t, pkg.Hash("testHash"), hash) //сверяем хеш

	err = db.Delete(ctx, "testUser") //удаляем пользователя
	require.NoError(t, err)          //ошибки нет

	_, ip, err = db.Get(ctx, "testUser") //ищем пользователя которого только что удалили
	require.Equal(t, sql.ErrNoRows, err)
}
