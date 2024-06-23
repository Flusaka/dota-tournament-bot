package types

type BaseMatch struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	TeamOne       *Team  `json:"teamOne"`
	TeamTwo       *Team  `json:"teamTwo"`
	ScheduledTime int64  `json:"scheduledTime"`
	StreamUrl     string `json:"streamUrl"`
}

type Match struct {
	BaseMatch
	Tournament BaseTournament
	Serie      BaseSerie
	League     BaseLeague
}
