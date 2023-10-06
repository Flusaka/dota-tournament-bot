package types

type Tier string

const (
	TierAmateur            Tier = "AMATEUR"
	TierProfessional       Tier = "PROFESSIONAL"
	TierMinor              Tier = "MINOR"
	TierMajor              Tier = "MAJOR"
	TierInternational      Tier = "INTERNATIONAL"
	TierDpcQualifier       Tier = "DPC_QUALIFIER"
	TierDpcLeagueQualifier Tier = "DPC_LEAGUE_QUALIFIER"
	TierDpcLeague          Tier = "DPC_LEAGUE"
	TierDpcLeagueFinals    Tier = "DPC_LEAGUE_FINALS"
)
