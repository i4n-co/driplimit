package driplimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyUpdateRemaining(t *testing.T) {
	// override global time functions, considering
	// time.Now() as 2024-01-01 10:30:00

	testClock := 0 * time.Minute
	now = func() time.Time {
		return time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC).Add(testClock)
	}
	since = func(t time.Time) time.Duration {
		return now().Sub(t)
	}

	key := Key{
		Ratelimit: &Ratelimit{
			State: &RatelimitState{
				LastRefilled: now(),
				Remaining:    10,
			},
			Limit:          10,
			RefillRate:     1,
			RefillInterval: Milliseconds{Duration: 1 * time.Minute},
		},
	}

	updated := key.UpdateRemaining()
	assert.False(t, updated)
	assert.Equal(t, int64(10), key.Ratelimit.State.Remaining)

	key.Ratelimit.State.Remaining = 5
	updated = key.UpdateRemaining()
	assert.False(t, updated)
	assert.Equal(t, int64(5), key.Ratelimit.State.Remaining)

	testClock += 1 * time.Minute
	updated = key.UpdateRemaining()
	assert.True(t, updated)
	assert.Equal(t, int64(6), key.Ratelimit.State.Remaining)

	testClock += 2 * time.Minute
	updated = key.UpdateRemaining()
	assert.True(t, updated)
	assert.Equal(t, int64(8), key.Ratelimit.State.Remaining)

	testClock += 5 * time.Minute
	updated = key.UpdateRemaining()
	assert.True(t, updated)
	assert.Equal(t, int64(10), key.Ratelimit.State.Remaining)
}
