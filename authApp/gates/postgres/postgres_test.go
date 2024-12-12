package store

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInsertGetDelete(t *testing.T) {
	conn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=medods_auth host=localhost sslmode=disable")
	if err != nil {
		require.NoError(t, err)
	}
	db := NewDB(conn)
	ctx := context.Background()

	err = db.Save(ctx, "testToken", "testUser")
	require.NoError(t, err)

	got, err := db.Get(ctx, "testToken")
	require.NoError(t, err)
	require.Equal(t, true, got)

	err = db.Delete(ctx, "testToken")
	require.NoError(t, err)

	got, err = db.Get(ctx, "testToken")
	require.NoError(t, err)
	require.Equal(t, false, got)
}
