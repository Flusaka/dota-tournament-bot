type BaseMatchesRequest = {
    readonly sort: "begin_at" | "-begin_at" | "detailed_stats" |
    "-detailed_stats" | "draw" | "-draw" | "end_at" | "-end_at" |
    "forfeit" | "-forfeit" | "id" | "-id" | "match_type" |
    "-match_type" | "modified_at" | "-modified_at" | "name" |
    "-name" | "number_of_games" | "-number_of_games" | "scheduled_at" |
    "-scheduled_at" | "slug" | "-slug" | "status" | "-status" |
    "tournament_id" | "-tournament_id" | "winner_id" | "-winner_id";
}

type UpcomingMatchesRequest = BaseMatchesRequest;
type RunningMatchesRequest = BaseMatchesRequest;

export {
    UpcomingMatchesRequest,
    RunningMatchesRequest
};