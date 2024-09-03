package proxycache

import (
	"context"
	"fmt"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/i4n-co/driplimit"
	"github.com/i4n-co/driplimit/pkg/config"
	"github.com/i4n-co/driplimit/pkg/generate"
)

// cache can store service keys, keys, and errors.
type cache struct {
	ServiceKeys *expirable.LRU[string, *driplimit.ServiceKey]
	Keys        *expirable.LRU[string, *driplimit.Key]
	Errors      *expirable.LRU[string, error]
}

func newCache(cfg *config.Config) *cache {
	return &cache{
		ServiceKeys: expirable.NewLRU[string, *driplimit.ServiceKey](cfg.ServiceKeysCacheSize, nil, cfg.CacheDuration),
		Keys:        expirable.NewLRU[string, *driplimit.Key](cfg.KeysCacheSize, nil, cfg.CacheDuration),
		Errors:      expirable.NewLRU[string, error](cfg.KeysCacheSize, nil, cfg.CacheDuration),
	}
}

// cacheRefresher refreshes the cache with the upstream asynchronously.
func (proxy *proxyCache) cacheRefresher(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			proxy.logger.Info("shutting down cache refresher...")
			return
		case order := <-proxy.refreshOrders:
			go func() {
				err := proxy.refreshCache(ctx, order)
				if err != nil {
					proxy.logger.Warn("failed to refresh cache", "err", err)
				}
			}()
		}
	}
}

type refreshOrder struct {
	driplimit.KeysCheckPayload
}

func (order refreshOrder) CacheKey() string {
	return order.KSID + generate.Hash(order.Token)
}

// refreshCache refreshes the cache with the upstream synchronously.
func (proxy *proxyCache) refreshCache(ctx context.Context, order refreshOrder) error {
	key, err := proxy.upstream.KeyCheck(ctx, order.KeysCheckPayload)
	if err != nil {
		proxy.cache.Errors.Add(order.CacheKey(), err)
		return fmt.Errorf("failed to check key: %w", err)
	}
	proxy.cache.Errors.Remove(order.CacheKey())
	proxy.cache.Keys.Add(order.CacheKey(), key)
	return nil
}
