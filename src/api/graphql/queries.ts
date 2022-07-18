export const LEAGUES_QUERY = `
    query {
        leagues(request: {tiers: DPC_LEAGUE, leagueEnded: false}) {
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