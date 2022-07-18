import moment from "moment";

type MatchDetails = {
    readonly matchId: number;

    // Match title (e.g. OG vs Nigma)
    readonly matchTitle: string;

    // Stream link
    readonly streamLink: string;

    // Start time
    readonly startTime: moment.Moment;
}

type DailyMatchesMessage = {
    readonly leagueName: string;
    readonly matches: MatchDetails[];
}

export {
    DailyMatchesMessage,
    MatchDetails
};