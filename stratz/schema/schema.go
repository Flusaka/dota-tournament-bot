// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package schema

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// GetLeaguesLeaguesLeagueType includes the requested fields of the GraphQL type LeagueType.
type GetLeaguesLeaguesLeagueType struct {
	Id          *int                                                        `json:"id"`
	DisplayName *string                                                     `json:"displayName"`
	Region      *LeagueRegion                                               `json:"region"`
	Tier        *LeagueTier                                                 `json:"tier"`
	Description *string                                                     `json:"description"`
	NodeGroups  []*GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType `json:"nodeGroups"`
}

// GetId returns GetLeaguesLeaguesLeagueType.Id, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueType) GetId() *int { return v.Id }

// GetDisplayName returns GetLeaguesLeaguesLeagueType.DisplayName, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueType) GetDisplayName() *string { return v.DisplayName }

// GetRegion returns GetLeaguesLeaguesLeagueType.Region, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueType) GetRegion() *LeagueRegion { return v.Region }

// GetTier returns GetLeaguesLeaguesLeagueType.Tier, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueType) GetTier() *LeagueTier { return v.Tier }

// GetDescription returns GetLeaguesLeaguesLeagueType.Description, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueType) GetDescription() *string { return v.Description }

// GetNodeGroups returns GetLeaguesLeaguesLeagueType.NodeGroups, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueType) GetNodeGroups() []*GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType {
	return v.NodeGroups
}

// GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType includes the requested fields of the GraphQL type LeagueNodeGroupType.
type GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType struct {
	Id            *int16                                                                         `json:"id"`
	Name          *string                                                                        `json:"name"`
	NodeGroupType *LeagueNodeGroupTypeEnum                                                       `json:"nodeGroupType"`
	Round         *byte                                                                          `json:"round"`
	Nodes         []*GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType `json:"nodes"`
}

// GetId returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType.Id, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType) GetId() *int16 { return v.Id }

// GetName returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType.Name, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType) GetName() *string { return v.Name }

// GetNodeGroupType returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType.NodeGroupType, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType) GetNodeGroupType() *LeagueNodeGroupTypeEnum {
	return v.NodeGroupType
}

// GetRound returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType.Round, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType) GetRound() *byte { return v.Round }

// GetNodes returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType.Nodes, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupType) GetNodes() []*GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType {
	return v.Nodes
}

// GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType includes the requested fields of the GraphQL type LeagueNodeType.
type GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType struct {
	Id            *int16                                                                                                `json:"id"`
	ScheduledTime *int64                                                                                                `json:"scheduledTime"`
	ActualTime    *int64                                                                                                `json:"actualTime"`
	NodeType      *LeagueNodeDefaultGroupEnum                                                                           `json:"nodeType"`
	HasStarted    *bool                                                                                                 `json:"hasStarted"`
	IsCompleted   *bool                                                                                                 `json:"isCompleted"`
	WinningNodeId *int16                                                                                                `json:"winningNodeId"`
	LosingNodeId  *int16                                                                                                `json:"losingNodeId"`
	Streams       []*GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType `json:"streams"`
	TeamOne       *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType           `json:"teamOne"`
	TeamTwo       *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType           `json:"teamTwo"`
}

// GetId returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.Id, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetId() *int16 {
	return v.Id
}

// GetScheduledTime returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.ScheduledTime, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetScheduledTime() *int64 {
	return v.ScheduledTime
}

// GetActualTime returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.ActualTime, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetActualTime() *int64 {
	return v.ActualTime
}

// GetNodeType returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.NodeType, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetNodeType() *LeagueNodeDefaultGroupEnum {
	return v.NodeType
}

// GetHasStarted returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.HasStarted, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetHasStarted() *bool {
	return v.HasStarted
}

// GetIsCompleted returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.IsCompleted, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetIsCompleted() *bool {
	return v.IsCompleted
}

// GetWinningNodeId returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.WinningNodeId, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetWinningNodeId() *int16 {
	return v.WinningNodeId
}

// GetLosingNodeId returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.LosingNodeId, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetLosingNodeId() *int16 {
	return v.LosingNodeId
}

// GetStreams returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.Streams, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetStreams() []*GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType {
	return v.Streams
}

// GetTeamOne returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.TeamOne, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetTeamOne() *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType {
	return v.TeamOne
}

// GetTeamTwo returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType.TeamTwo, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeType) GetTeamTwo() *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType {
	return v.TeamTwo
}

// GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType includes the requested fields of the GraphQL type LeagueStreamType.
type GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType struct {
	Id         *int      `json:"id"`
	LanguageId *Language `json:"languageId"`
	Name       *string   `json:"name"`
	StreamUrl  *string   `json:"streamUrl"`
}

// GetId returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType.Id, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType) GetId() *int {
	return v.Id
}

// GetLanguageId returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType.LanguageId, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType) GetLanguageId() *Language {
	return v.LanguageId
}

// GetName returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType.Name, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType) GetName() *string {
	return v.Name
}

// GetStreamUrl returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType.StreamUrl, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeStreamsLeagueStreamType) GetStreamUrl() *string {
	return v.StreamUrl
}

// GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType includes the requested fields of the GraphQL type TeamType.
type GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType struct {
	Id   int     `json:"id"`
	Name *string `json:"name"`
	Tag  *string `json:"tag"`
}

// GetId returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType.Id, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType) GetId() int {
	return v.Id
}

// GetName returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType.Name, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType) GetName() *string {
	return v.Name
}

// GetTag returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType.Tag, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamOneTeamType) GetTag() *string {
	return v.Tag
}

// GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType includes the requested fields of the GraphQL type TeamType.
type GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType struct {
	Id   int     `json:"id"`
	Name *string `json:"name"`
	Tag  *string `json:"tag"`
}

// GetId returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType.Id, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType) GetId() int {
	return v.Id
}

// GetName returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType.Name, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType) GetName() *string {
	return v.Name
}

// GetTag returns GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType.Tag, and is useful for accessing the field via an interface.
func (v *GetLeaguesLeaguesLeagueTypeNodeGroupsLeagueNodeGroupTypeNodesLeagueNodeTypeTeamTwoTeamType) GetTag() *string {
	return v.Tag
}

// GetLeaguesResponse is returned by GetLeagues on success.
type GetLeaguesResponse struct {
	// Find league details by searching for leagues using a LeagueRequest.
	Leagues []*GetLeaguesLeaguesLeagueType `json:"leagues"`
}

// GetLeagues returns GetLeaguesResponse.Leagues, and is useful for accessing the field via an interface.
func (v *GetLeaguesResponse) GetLeagues() []*GetLeaguesLeaguesLeagueType { return v.Leagues }

type Language string

const (
	LanguageEnglish    Language = "ENGLISH"
	LanguageBrazilian  Language = "BRAZILIAN"
	LanguageBulgarian  Language = "BULGARIAN"
	LanguageCzech      Language = "CZECH"
	LanguageDanish     Language = "DANISH"
	LanguageDutch      Language = "DUTCH"
	LanguageFinnish    Language = "FINNISH"
	LanguageFrench     Language = "FRENCH"
	LanguageGerman     Language = "GERMAN"
	LanguageGreek      Language = "GREEK"
	LanguageHungarian  Language = "HUNGARIAN"
	LanguageItalian    Language = "ITALIAN"
	LanguageJapanese   Language = "JAPANESE"
	LanguageKorean     Language = "KOREAN"
	LanguageKoreana    Language = "KOREANA"
	LanguageNorwegian  Language = "NORWEGIAN"
	LanguagePolish     Language = "POLISH"
	LanguagePortuguese Language = "PORTUGUESE"
	LanguageRomanian   Language = "ROMANIAN"
	LanguageRussian    Language = "RUSSIAN"
	LanguageSChinese   Language = "S_CHINESE"
	LanguageSpanish    Language = "SPANISH"
	LanguageSwedish    Language = "SWEDISH"
	LanguageTChinese   Language = "T_CHINESE"
	LanguageThai       Language = "THAI"
	LanguageTurkish    Language = "TURKISH"
	LanguageUkrainian  Language = "UKRAINIAN"
)

type LeagueNodeDefaultGroupEnum string

const (
	LeagueNodeDefaultGroupEnumInvalid     LeagueNodeDefaultGroupEnum = "INVALID"
	LeagueNodeDefaultGroupEnumBestOfOne   LeagueNodeDefaultGroupEnum = "BEST_OF_ONE"
	LeagueNodeDefaultGroupEnumBestOfThree LeagueNodeDefaultGroupEnum = "BEST_OF_THREE"
	LeagueNodeDefaultGroupEnumBestOfFive  LeagueNodeDefaultGroupEnum = "BEST_OF_FIVE"
	LeagueNodeDefaultGroupEnumBestOfTwo   LeagueNodeDefaultGroupEnum = "BEST_OF_TWO"
)

type LeagueNodeGroupTypeEnum string

const (
	LeagueNodeGroupTypeEnumInvalid                LeagueNodeGroupTypeEnum = "INVALID"
	LeagueNodeGroupTypeEnumOrganizational         LeagueNodeGroupTypeEnum = "ORGANIZATIONAL"
	LeagueNodeGroupTypeEnumRoundRobin             LeagueNodeGroupTypeEnum = "ROUND_ROBIN"
	LeagueNodeGroupTypeEnumSwiss                  LeagueNodeGroupTypeEnum = "SWISS"
	LeagueNodeGroupTypeEnumBracketSingle          LeagueNodeGroupTypeEnum = "BRACKET_SINGLE"
	LeagueNodeGroupTypeEnumBracketDoubleSeedLoser LeagueNodeGroupTypeEnum = "BRACKET_DOUBLE_SEED_LOSER"
	LeagueNodeGroupTypeEnumBracketDoubleAllWinner LeagueNodeGroupTypeEnum = "BRACKET_DOUBLE_ALL_WINNER"
	LeagueNodeGroupTypeEnumShowmatch              LeagueNodeGroupTypeEnum = "SHOWMATCH"
	LeagueNodeGroupTypeEnumGsl                    LeagueNodeGroupTypeEnum = "GSL"
)

type LeagueRegion string

const (
	LeagueRegionUnset  LeagueRegion = "UNSET"
	LeagueRegionNa     LeagueRegion = "NA"
	LeagueRegionSa     LeagueRegion = "SA"
	LeagueRegionEurope LeagueRegion = "EUROPE"
	LeagueRegionCis    LeagueRegion = "CIS"
	LeagueRegionChina  LeagueRegion = "CHINA"
	LeagueRegionSea    LeagueRegion = "SEA"
)

type LeagueTier string

const (
	LeagueTierUnset              LeagueTier = "UNSET"
	LeagueTierAmateur            LeagueTier = "AMATEUR"
	LeagueTierProfessional       LeagueTier = "PROFESSIONAL"
	LeagueTierMinor              LeagueTier = "MINOR"
	LeagueTierMajor              LeagueTier = "MAJOR"
	LeagueTierInternational      LeagueTier = "INTERNATIONAL"
	LeagueTierDpcQualifier       LeagueTier = "DPC_QUALIFIER"
	LeagueTierDpcLeagueQualifier LeagueTier = "DPC_LEAGUE_QUALIFIER"
	LeagueTierDpcLeague          LeagueTier = "DPC_LEAGUE"
	LeagueTierDpcLeagueFinals    LeagueTier = "DPC_LEAGUE_FINALS"
)

// __GetLeaguesInput is used internally by genqlient
type __GetLeaguesInput struct {
	Tiers       []*LeagueTier `json:"tiers"`
	LeagueEnded *bool         `json:"leagueEnded"`
}

// GetTiers returns __GetLeaguesInput.Tiers, and is useful for accessing the field via an interface.
func (v *__GetLeaguesInput) GetTiers() []*LeagueTier { return v.Tiers }

// GetLeagueEnded returns __GetLeaguesInput.LeagueEnded, and is useful for accessing the field via an interface.
func (v *__GetLeaguesInput) GetLeagueEnded() *bool { return v.LeagueEnded }

// The query or mutation executed by GetLeagues.
const GetLeagues_Operation = `
query GetLeagues ($tiers: [LeagueTier], $leagueEnded: Boolean) {
	leagues(request: {tiers:$tiers,leagueEnded:$leagueEnded}) {
		id
		displayName
		region
		tier
		description
		nodeGroups {
			id
			name
			nodeGroupType
			round
			nodes {
				id
				scheduledTime
				actualTime
				nodeType
				hasStarted
				isCompleted
				winningNodeId
				losingNodeId
				streams {
					id
					languageId
					name
					streamUrl
				}
				teamOne {
					id
					name
					tag
				}
				teamTwo {
					id
					name
					tag
				}
			}
		}
	}
}
`

func GetLeagues(
	ctx context.Context,
	client graphql.Client,
	tiers []*LeagueTier,
	leagueEnded *bool,
) (*GetLeaguesResponse, error) {
	req := &graphql.Request{
		OpName: "GetLeagues",
		Query:  GetLeagues_Operation,
		Variables: &__GetLeaguesInput{
			Tiers:       tiers,
			LeagueEnded: leagueEnded,
		},
	}
	var err error

	var data GetLeaguesResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
