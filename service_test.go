package driplimit_test

import (
	"errors"
	"testing"

	"github.com/i4n-co/driplimit"
	"github.com/stretchr/testify/assert"
)

func TestErrNotFound(t *testing.T) {
	err := driplimit.ErrItemNotFound("item")
	assert.True(t, errors.Is(err, driplimit.ErrNotFound))

	target := driplimit.ErrItemNotFound("")
	assert.True(t, errors.As(driplimit.ErrItemNotFound("key"), &target))

	assert.Equal(t, "key not found", target.Error())
}
