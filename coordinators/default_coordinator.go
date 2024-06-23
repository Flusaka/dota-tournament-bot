package coordinators

import (
	"github.com/flusaka/dota-tournament-bot/bot"
	"github.com/flusaka/dota-tournament-bot/cache"
	"github.com/flusaka/dota-tournament-bot/queries"
	"github.com/flusaka/dota-tournament-bot/types"
)

type DefaultQueryCoordinator struct {
	dataSource       bot.DataSource
	queryResultCache cache.QueryResultCache
}

func NewDefaultQueryCoordinator(dataSource bot.DataSource, queryResultCache cache.QueryResultCache) *DefaultQueryCoordinator {
	return &DefaultQueryCoordinator{
		dataSource:       dataSource,
		queryResultCache: queryResultCache,
	}
}

func (c *DefaultQueryCoordinator) GetTournaments(query *queries.GetTournaments) ([]types.Tournament, error) {
	if c.queryResultCache.Has(query) {
		cached := c.queryResultCache.Get(query)
		return cached.([]types.Tournament), nil
	}

	// TODO: Would be good to have a way to execute one query for multiple sources instead of executing per-source (i.e. channel)
	result, err := c.dataSource.GetTournaments(query)
	if err != nil {
		return nil, err
	}
	c.queryResultCache.Set(query, result)
	return result, nil
}

func (c *DefaultQueryCoordinator) GetMatches(query *queries.GetMatches) ([]types.Match, error) {
	if c.queryResultCache.Has(query) {
		cached := c.queryResultCache.Get(query)
		return cached.([]types.Match), nil
	}

	result, err := c.dataSource.GetMatches(query)
	if err != nil {
		return nil, err
	}
	c.queryResultCache.Set(query, result)
	return result, nil
}

func (c *DefaultQueryCoordinator) GetUpcomingMatches(query *queries.GetUpcomingMatches) ([]types.Match, error) {
	if c.queryResultCache.Has(query) {
		cached := c.queryResultCache.Get(query)
		return cached.([]types.Match), nil
	}

	result, err := c.dataSource.GetUpcomingMatches(query)
	if err != nil {
		return nil, err
	}
	c.queryResultCache.Set(query, result)
	return result, nil
}
