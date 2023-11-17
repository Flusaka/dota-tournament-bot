package cache

import (
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/jellydator/ttlcache/v3"
	"time"
)

type DefaultQueryResultCache struct {
	cache *ttlcache.Cache[uint64, interface{}]
}

func NewDefaultQueryResultCache(cacheTime time.Duration) *DefaultQueryResultCache {
	return &DefaultQueryResultCache{
		cache: ttlcache.New[uint64, interface{}](
			ttlcache.WithTTL[uint64, interface{}](cacheTime),
			ttlcache.WithDisableTouchOnHit[uint64, interface{}](),
		),
	}
}

func (q *DefaultQueryResultCache) Get(query queries.Query) interface{} {
	hash, err := query.HashCode()
	if err != nil {
		return nil
	}
	return q.cache.Get(hash).Value()
}

func (q *DefaultQueryResultCache) Has(query queries.Query) bool {
	hash, err := query.HashCode()
	if err != nil {
		return false
	}
	return q.cache.Has(hash)
}

func (q *DefaultQueryResultCache) Set(query queries.Query, result interface{}) {
	hash, err := query.HashCode()
	if err == nil {
		q.cache.Set(hash, result, ttlcache.DefaultTTL)
	}
}
