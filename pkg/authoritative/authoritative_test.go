package authoritative_test

import (
	"context"
	"testing"
	"time"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/authoritative"
	"github.com/i4n-co/driplimit/pkg/store"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	ctx := context.Background()
	dbHandler, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	sqlite, err := store.New(ctx, dbHandler)
	if err != nil {
		t.Fatal(err)
	}
	app := authoritative.NewService(sqlite)

	ks, err := app.KeyspaceCreate(ctx, driplimit.KeyspaceCreatePayload{
		Name: "test key space",
		Ratelimit: driplimit.RatelimitPayload{
			Limit:          10,
			RefillRate:     1,
			RefillInterval: 60 * 1000,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, token, err := app.KeyCreate(ctx, driplimit.KeyCreatePayload{
		KSID:      ks.KSID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30 * 12 * 10),
		Ratelimit: driplimit.RatelimitPayload{
			Limit:          100,
			RefillRate:     1,
			RefillInterval: 10,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	key, err := app.KeyGet(ctx, driplimit.KeyGetPayload{KSID: ks.KSID, Token: *token})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(100), key.Ratelimit.State.Remaining)
	if key.Ratelimit.State.Remaining > 0 {
		_, err = app.KeyCheck(ctx, driplimit.KeysCheckPayload{KSID: key.KSID, Token: *token})
		if err != nil {
			t.Fatal(err)
		}
	}

	key, err = app.KeyGet(ctx, driplimit.KeyGetPayload{KSID: ks.KSID, Token: *token})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(99), key.Ratelimit.State.Remaining)

	time.Sleep(10 * time.Millisecond)

	key, err = app.KeyGet(ctx, driplimit.KeyGetPayload{KSID: ks.KSID, Token: *token})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(100), key.Ratelimit.State.Remaining)

	_, token, err = app.KeyCreate(ctx, driplimit.KeyCreatePayload{
		KSID:      ks.KSID,
		ExpiresAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	key, _ = app.KeyGet(ctx, driplimit.KeyGetPayload{KSID: ks.KSID, Token: *token})
	assert.True(t, key.Expired())
}

func TestUnconfiguredRateLimit(t *testing.T) {
	ctx := context.Background()
	dbHandler, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	sqlite, err := store.New(ctx, dbHandler)
	if err != nil {
		t.Fatal(err)
	}
	app := authoritative.NewService(sqlite)

	ks, err := app.KeyspaceCreate(ctx, driplimit.KeyspaceCreatePayload{
		Name: "test key space",
	})
	if err != nil {
		t.Fatal(err)
	}

	key, token, err := app.KeyCreate(ctx, driplimit.KeyCreatePayload{
		KSID:      ks.KSID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30 * 12 * 10),
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = app.KeyCheck(ctx, driplimit.KeysCheckPayload{KSID: key.KSID, Token: *token})
	assert.NoError(t, err)
}
