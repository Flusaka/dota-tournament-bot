type MatchDetails = {
    readonly matchId: number;

    // Match title (e.g. OG vs Nigma)
    readonly matchTitle: string;

    // Stream link
    readonly streamLink: string;

    // Start time
    readonly startTime: Date;
}

type DailyMatchesMessage = {
    readonly tournamentName: string;
    readonly matches: MatchDetails[];
}

export {
    DailyMatchesMessage,
    MatchDetails
};