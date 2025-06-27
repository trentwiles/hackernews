package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresEnvs(t *testing.T) {
	LoadEnv()
	assert.Equal(t, GetEnv("POSTGRES_USERNAME"), "postgres", "default username from .env")
	assert.Equal(t, GetEnv("POSTGRES_DB"), "hn", "default database name from .env")
	assert.Equal(t, GetEnv("nonsense"), "", "bogus .env test")
}