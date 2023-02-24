package config_test

import (
	"os"
	"testing"

	"github.com/polyxia-org/morty-registry/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_Parse_Environment(t *testing.T) {
	t.Parallel()

	os.Setenv("REGISTRY_STORAGE_S3_ENDPOINT", "http://example.com")
	os.Setenv("REGISTRY_PORT", "8080")

	cfg, err := config.Parse()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, cfg.Port, 8080)
	assert.Equal(t, cfg.Storage.S3.Endpoint, "http://example.com")
}
