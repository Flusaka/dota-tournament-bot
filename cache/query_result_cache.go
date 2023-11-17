package cache

import (
	"github.com/flusaka/dota-tournament-bot/queries"
)

type QueryResultCache interface {
	Get(query queries.Query) interface{}
	Has(query queries.Query) bool
	Set(query queries.Query, result interface{})
}
