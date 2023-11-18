package coordinators

import (
	"github.com/flusaka/dota-tournament-bot/cache"
	"github.com/flusaka/dota-tournament-bot/datasource"
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
)

type DefaultQueryCoordinator struct {
	dataSourceClient datasource.Client
	queryResultCache cache.QueryResultCache
}

func NewDefaultQueryCoordinator(dataSourceClient datasource.Client, queryResultCache cache.QueryResultCache) DefaultQueryCoordinator {
	return DefaultQueryCoordinator{
		dataSourceClient: dataSourceClient,
		queryResultCache: queryResultCache,
	}
}

func (receiver DefaultQueryCoordinator) GetLeagues(query *queries.GetLeagues) ([]*types.League, error) {
	if receiver.queryResultCache.Has(query) {
		cached := receiver.queryResultCache.Get(query)
		return cached.([]*types.League), nil
	}

	// TODO: Would be good to have a way to execute one query for multiple sources instead of executing per-source (i.e. channel)
	result, err := receiver.dataSourceClient.GetLeagues(query)
	if err != nil {
		return nil, err
	}
	receiver.queryResultCache.Set(query, result)
	return result, nil
}
