package proxycache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/config"
	"github.com/i4n-co/driplimit/pkg/generate"
)

// proxyCache is a driplimit proxy that caches keys and predicts checks.
type proxyCache struct {
	cache         *cache
	refreshOrders chan refreshOrder
	upstream      driplimit.ServiceWithToken
	cfg           *config.Config
	logger        *slog.Logger
}

// NewServiceProxyCache creates a new proxy cache service.
func NewServiceProxyCache(ctx context.Context, cfg *config.Config, upstream driplimit.ServiceWithToken) driplimit.ServiceWithToken {
	proxy := &proxyCache{
		upstream:      upstream,
		cache:         newCache(cfg),
		refreshOrders: make(chan refreshOrder),
		cfg:           cfg,
		logger:        cfg.Logger().With("component", "proxycache"),
	}

	proxy.logger.Info("starting proxy cache...")
	go proxy.cacheRefresher(ctx)
	return proxy
}

// proxyCacheWithToken is a proxy cache with a token.
type proxyCacheWithToken struct {
	*proxyCache
	token string
}

// WithToken creates a new proxy cache with a token.
func (proxy *proxyCache) WithToken(token string) driplimit.Service {
	return &proxyCacheWithToken{
		proxyCache: proxy,
		token:      token,
	}
}

// KeyCheck checks if a key is valid and predicts the next check.
// this method tries to be as asynchronous as possible.
func (proxy *proxyCacheWithToken) KeyCheck(ctx context.Context, payload driplimit.KeysCheckPayload) (key *driplimit.Key, err error) {
	sk, err := proxy.ServiceKeyGet(ctx, driplimit.ServiceKeyGetPayload{
		Token: generate.Hash(proxy.token),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get service key: %w", err)
	}

	// cache can be populated by multiple service keys. Therfore, we need to check if the
	// service key is allowed to check the key in the local cache.
	if !sk.Admin && !sk.KeyspacesPolicies.Can(driplimit.Read, payload.KSID) {
		return nil, driplimit.ErrUnauthorized
	}

	refreshOrder := refreshOrder{payload, proxy.token}
	refreshErr, _ := proxy.cache.Errors.Get(refreshOrder.CacheKey())

	if errors.Is(refreshErr, driplimit.ErrKeyExpired) {
		return nil, driplimit.ErrKeyExpired
	}

	key, found := proxy.cache.Keys.Get(refreshOrder.CacheKey())
	if !found {
		err = proxy.refreshCache(ctx, refreshOrder)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh cache: %w", err)
		}
		key, found = proxy.cache.Keys.Get(refreshOrder.CacheKey())
		if !found {
			return nil, fmt.Errorf("key not found in cache after synchronous refresh")
		}
		return key, nil
	}
	// notify ahead the cache refresher to refresh the cache asynchronously
	proxy.refreshOrders <- refreshOrder

	if !key.ConfiguredRatelimit() {
		return key, nil
	}

	if key.UpdateRemaining() && errors.Is(refreshErr, driplimit.ErrRateLimitExceeded) {
		proxy.cache.Errors.Remove(refreshOrder.CacheKey())
	}

	if key.Ratelimit.State.Remaining <= 0 {
		return nil, driplimit.ErrRateLimitExceeded
	}

	key.Ratelimit.State.Remaining--
	if key.Ratelimit.State.Remaining < 0 {
		key.Ratelimit.State.Remaining = 0
	}
	key.LastUsed = time.Now()
	proxy.cache.Errors.Remove(refreshOrder.CacheKey())

	return key, nil
}

func (proxy *proxyCacheWithToken) KeyCreate(ctx context.Context, payload driplimit.KeyCreatePayload) (key *driplimit.Key, token *string, err error) {
	return proxy.upstream.WithToken(proxy.token).KeyCreate(ctx, payload)
}

func (proxy *proxyCacheWithToken) KeyGet(ctx context.Context, payload driplimit.KeyGetPayload) (key *driplimit.Key, err error) {
	return proxy.upstream.WithToken(proxy.token).KeyGet(ctx, payload)
}

func (proxy *proxyCacheWithToken) KeyList(ctx context.Context, payload driplimit.KeyListPayload) (klist *driplimit.KeyList, err error) {
	return proxy.upstream.WithToken(proxy.token).KeyList(ctx, payload)
}

func (proxy *proxyCacheWithToken) KeyDelete(ctx context.Context, payload driplimit.KeyDeletePayload) (err error) {
	return proxy.upstream.WithToken(proxy.token).KeyDelete(ctx, payload)
}

func (proxy *proxyCacheWithToken) KeyspaceCreate(ctx context.Context, payload driplimit.KeyspaceCreatePayload) (keyspace *driplimit.Keyspace, err error) {
	return proxy.upstream.WithToken(proxy.token).KeyspaceCreate(ctx, payload)
}

func (proxy *proxyCacheWithToken) KeyspaceGet(ctx context.Context, payload driplimit.KeyspaceGetPayload) (keyspace *driplimit.Keyspace, err error) {
	return proxy.upstream.WithToken(proxy.token).KeyspaceGet(ctx, payload)
}

func (proxy *proxyCacheWithToken) KeyspaceList(ctx context.Context, payload driplimit.KeyspaceListPayload) (kslist *driplimit.KeyspaceList, err error) {
	return proxy.upstream.WithToken(proxy.token).KeyspaceList(ctx, payload)
}

func (proxy *proxyCacheWithToken) KeyspaceDelete(ctx context.Context, payload driplimit.KeyspaceDeletePayload) (err error) {
	return proxy.upstream.WithToken(proxy.token).KeyspaceDelete(ctx, payload)
}

func (proxy *proxyCacheWithToken) ServiceKeyGet(ctx context.Context, payload driplimit.ServiceKeyGetPayload) (sk *driplimit.ServiceKey, err error) {
	sk, found := proxy.cache.ServiceKeys.Get(generate.Hash(proxy.token))
	if found {
		return sk, nil
	}
	sk, err = proxy.upstream.WithToken(proxy.token).ServiceKeyGet(ctx, payload)
	if err != nil {
		return nil, err
	}
	proxy.cache.ServiceKeys.Add(generate.Hash(proxy.token), sk)
	return sk, nil
}

func (proxy *proxyCacheWithToken) ServiceKeyCreate(ctx context.Context, payload driplimit.ServiceKeyCreatePayload) (sk *driplimit.ServiceKey, err error) {
	return proxy.upstream.WithToken(proxy.token).ServiceKeyCreate(ctx, payload)
}

func (proxy *proxyCacheWithToken) ServiceKeyList(ctx context.Context, payload driplimit.ServiceKeyListPayload) (sklist *driplimit.ServiceKeyList, err error) {
	return proxy.upstream.WithToken(proxy.token).ServiceKeyList(ctx, payload)
}

func (proxy *proxyCacheWithToken) ServiceKeyDelete(ctx context.Context, payload driplimit.ServiceKeyDeletePayload) (err error) {
	return proxy.upstream.WithToken(proxy.token).ServiceKeyDelete(ctx, payload)
}
