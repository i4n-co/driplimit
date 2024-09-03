package api_test

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/api"
	"github.com/i4n-co/driplimit/pkg/authoritative"
	"github.com/i4n-co/driplimit/pkg/client"
	"github.com/i4n-co/driplimit/pkg/config"
	"github.com/i4n-co/driplimit/pkg/store"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

var cli *client.HTTP

func init() {
	ctx := context.Background()

	cfg, err := config.FromEnv(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	store, err := store.New(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	err = store.InitRootServiceKeyToken(ctx, "t0k3n")
	if err != nil {
		log.Fatal(err)
	}

	authoritative := authoritative.NewService(store)
	authorizer := driplimit.NewAuthorizer(authoritative)
	validator := driplimit.NewServiceValidator(authorizer)

	server := api.New(cfg, validator)

	cli = client.New("http://localhost.test").WithSendRequestFunc(server.Test)
}

func TestDriplimitAPI(t *testing.T) {
	ctx := context.Background()
	// SERVICE KEYS
	// current shoud return unauthorize as the token is uninitialized yet

	_, err := cli.ServiceKeyCurrent(ctx)
	assert.ErrorIs(t, err, driplimit.ErrUnauthorized)

	// Set default service key token
	cli = cli.WithServiceToken("t0k3n")

	// Should returns sk_root
	sk, err := cli.ServiceKeyCurrent(ctx)
	assert.NoError(t, err)
	assert.Equal(t, sk.SKID, "sk_root")

	// KEYSPACES
	// create keyspace without ratelimit configuration
	ks1, err := cli.KeyspaceCreate(ctx, driplimit.KeyspaceCreatePayload{
		Name:       "test",
		KeysPrefix: "test_",
	})
	assert.NoError(t, err)
	assert.Equal(t, "test", ks1.Name)

	// create keyspace with ratelimit configuration
	withRateLimitKS, err := cli.KeyspaceCreate(ctx, driplimit.KeyspaceCreatePayload{
		Name:       "test_with_ratelimit",
		KeysPrefix: "test_wrl_",
		Ratelimit: driplimit.RatelimitPayload{
			Limit:          10,
			RefillRate:     1,
			RefillInterval: driplimit.Milliseconds{Duration: time.Second},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "test_with_ratelimit", withRateLimitKS.Name)
	assert.Equal(t, int64(10), withRateLimitKS.Ratelimit.Limit)

	// KEYS
	// create should fail as expiration is not configured
	_, err = cli.KeyCreate(ctx, driplimit.KeyCreatePayload{
		KSID: withRateLimitKS.KSID,
	})
	assert.ErrorIs(t, err, driplimit.ErrInvalidExpiration)

	// should not fail as key creation is valid
	expiresAt := time.Date(2024, 01, 01, 01, 01, 0, 0, time.UTC)
	k, err := cli.KeyCreate(ctx, driplimit.KeyCreatePayload{
		KSID:      withRateLimitKS.KSID,
		ExpiresAt: expiresAt,
	})
	assert.NotEmpty(t, k.Token)
	assert.Equal(t, k.ExpiresAt, expiresAt)
	assert.True(t, strings.HasPrefix(k.Token, "test_wrl_"))
	assert.NoError(t, err)

	// should fail as key is expired
	_, err = cli.KeyCheck(ctx, driplimit.KeysCheckPayload{KSID: withRateLimitKS.KSID, Token: k.Token})
	assert.ErrorIs(t, err, driplimit.ErrKeyExpired)

	expiresAt = time.Date(2050, 01, 01, 01, 01, 0, 0, time.UTC)
	k, err = cli.KeyCreate(ctx, driplimit.KeyCreatePayload{
		KSID:      withRateLimitKS.KSID,
		ExpiresAt: expiresAt,
	})
	assert.NotEmpty(t, k.Token)
	assert.Equal(t, k.ExpiresAt, expiresAt)
	assert.NoError(t, err)

	token := k.Token

	k, err = cli.KeyGet(ctx, driplimit.KeyGetPayload{KSID: k.KSID, KID: k.KID})
	assert.NoError(t, err)
	assert.True(t, !k.Expired())

	// should succeed as key is not expired
	k, err = cli.KeyCheck(ctx, driplimit.KeysCheckPayload{KSID: k.KSID, Token: token})
	assert.NoError(t, err)
	assert.Equal(t, k.Ratelimit.State.Remaining, int64(9))

	lkeys, err := cli.KeyList(ctx, driplimit.KeyListPayload{
		KSID: withRateLimitKS.KSID,
	})
	assert.NoError(t, err)
	assert.Len(t, lkeys.Keys, 2)

	err = cli.KeyDelete(ctx, driplimit.KeyDeletePayload{
		KSID: withRateLimitKS.KSID,
		KID:  lkeys.Keys[1].KID,
	})
	assert.NoError(t, err)

	lkeys, err = cli.KeyList(ctx, driplimit.KeyListPayload{
		KSID: withRateLimitKS.KSID,
	})
	assert.NoError(t, err)
	assert.Len(t, lkeys.Keys, 1)

	nsk, err := cli.ServiceKeyCreate(ctx, driplimit.ServiceKeyCreatePayload{
		KeyspacesPolicies: driplimit.Policies{
			withRateLimitKS.KSID: driplimit.Policy{
				Read: true,
			},
		},
		Description: "read restricted service key on test_with_ratelimit keyspace",
	})
	assert.NoError(t, err)
	assert.False(t, nsk.Admin)

	k, err = cli.WithServiceToken(nsk.Token).KeyCheck(ctx, driplimit.KeysCheckPayload{
		KSID:  withRateLimitKS.KSID,
		Token: token,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(8), k.Ratelimit.State.Remaining)

	// Should fail as nsk is read only
	err = cli.WithServiceToken(nsk.Token).KeyDelete(ctx, driplimit.KeyDeletePayload{
		KSID: withRateLimitKS.KSID,
		KID:  k.KID,
	})
	assert.ErrorIs(t, driplimit.ErrUnauthorized, err)
}
