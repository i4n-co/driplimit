package store_test

import (
	"github.com/i4n-co/driplimit/pkg/store"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeNano(t *testing.T) {
	tn := new(store.TimeNano)
	err := tn.Scan("2020-01-01T00:00:00Z")
	assert.Error(t, err)

	rawnano := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()
	err = tn.Scan(rawnano)
	if err != nil {
		t.Fatal(err)
	}

	v, err := tn.Value()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rawnano, v)

}
