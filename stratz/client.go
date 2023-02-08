package stratz

import (
	"context"
	"github.com/Khan/genqlient/graphql"
	"github.com/flusaka/dota-tournament-bot/stratz/schema"
	"net/http"
)

type Client struct {
	token     string
	gqlClient graphql.Client
}

type authenticatedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (at *authenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+at.token)
	return at.wrapped.RoundTrip(req)
}

func NewClient(token string) *Client {
	client := new(Client)
	client.token = token
	return client
}

func (c *Client) Initialise() {
	httpClient := http.Client{
		Transport: &authenticatedTransport{
			c.token,
			http.DefaultTransport,
		},
	}
	c.gqlClient = graphql.NewClient("https://api.stratz.com/graphql", &httpClient)
}

func (c *Client) GetLeagues(tiers []schema.LeagueTier, finished bool) ([]schema.GetLeaguesLeaguesLeagueType, error) {
	response, err := schema.GetLeagues(context.Background(), c.gqlClient, tiers, finished)
	if err != nil {
		return nil, err
	}
	return response.Leagues, nil
}

func (c *Client) GetActiveLeagues(tiers []schema.LeagueTier) ([]schema.GetLeaguesLeaguesLeagueType, error) {
	return c.GetLeagues(tiers, false)
}

func (c *Client) GetMatchesInActiveLeagues(tiers []schema.LeagueTier) ([]schema.GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType, error) {
	leagues, err := c.GetActiveLeagues(tiers)
	if err != nil {
		return nil, err
	}

	var matches []schema.GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType
	for _, league := range leagues {
		for _, nodeGroup := range league.NodeGroups {
			matches = append(matches, nodeGroup.Nodes...)
		}
	}
	return matches, nil
}
