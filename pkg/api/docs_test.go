package api

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/i4n-co/driplimit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDocs(t *testing.T) {
	cfg, err := config.FromEnv(context.Background())
	assert.NoError(t, err)
	api := New(cfg, nil)
	docs, err := api.GenerateDocs()
	assert.NoError(t, err)

	_, err = json.MarshalIndent(docs, "", "  ")
	assert.NoError(t, err)
}
