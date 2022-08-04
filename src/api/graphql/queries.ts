export const LEAGUES_QUERY = `
    query Leagues($tiers: [LeagueTier], $leagueEnded: Boolean) {
        leagues(request: {tiers: $tiers, leagueEnded: $leagueEnded}) {
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
`;