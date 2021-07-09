type Match = {
    // Game name
    readonly gameName: string;

    // Stream
    readonly streamLink: string;

    // Start time
    readonly startTime: Date;
}

type DailyMatchesMessage = {
    readonly matches: Match[];
}

export {
    DailyMatchesMessage,
    Match
};