package store

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInsertGetCheckExistDelete(t *testing.T) {
	conn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=medods_auth host=localhost sslmode=disable")
	if err != nil {
		require.NoError(t, err)
	}
	db := NewDB(conn)
	ctx := context.Background()

	err = db.Save(ctx, "testToken", "testUser", "255.255.255.255")
	require.NoError(t, err)

	got, ip, err := db.Get(ctx, "testUser", "testToken2")
	require.Error(t, err)
	got, ip, err = db.Get(ctx, "testUser2", "testToken")
	require.Error(t, err)
	got, ip, err = db.Get(ctx, "testUser", "testToken")
	require.Equal(t, true, got)
	require.Equal(t, ip, "255.255.255.255")

	exist, err := db.CheckUserExist(ctx, "testUser")
	require.NoError(t, err)
	require.True(t, exist)

	err = db.Delete(ctx, "testUser")
	require.NoError(t, err)

	got, ip, err = db.Get(ctx, "testUser", "testToken")
	require.Error(t, err)
	require.Equal(t, false, got)
	exist, err = db.CheckUserExist(ctx, "testUser")
	require.NoError(t, err)
	require.False(t, exist)
}
