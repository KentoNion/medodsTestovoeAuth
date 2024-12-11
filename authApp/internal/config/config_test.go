package config

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Setenv("CONFIG_PATH", "../../config.yaml")
	cfg, err := MustLoad()
	require.NoError(t, err)
	fmt.Println(cfg.DB.DbName)
	require.NotNil(t, cfg.DB.DbName)
}
