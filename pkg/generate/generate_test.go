package generate_test

import (
	"github.com/i4n-co/driplimit/pkg/generate"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerated(t *testing.T) {
	id := generate.ID()
	assert.Len(t, id, 22)
}

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generate.ID()
	}
}
