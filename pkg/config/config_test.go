package config_test

import (
	"context"
	"strings"
	"testing"

	"github.com/i4n-co/driplimit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	cfg := new(config.Config)
	deflt, err := cfg.Defaults()
	assert.NoError(t, err)
	assert.Equal(t, config.ConfigHelp{Default: "127.0.0.1", Description: "address to listen on"}, deflt["ADDR"])
}

func TestFromEnvFile(t *testing.T) {
	file := "PORT=7131\n"
	file += "# comment\n"
	file += "\n" // empty line
	file += "LOG_SEVERITY=debug\n"

	cfg, err := config.FromEnvFile(context.Background(), strings.NewReader(file))
	assert.NoError(t, err)
	assert.Equal(t, "debug", cfg.LogSeverity)
	assert.Equal(t, 7131, cfg.Port)
}
