type BaseTournamentsRequest = {
    readonly sort: 'begin_at' | '-begin_at' | 'end_at' |
    '-end_at' | 'id' | '-id' | 'modified_at' |
    '-modified_at' | 'name' | '-name' | 'prizepool' |
    '-prizepool' | 'serie_id' | '-serie_id' | 'slug' |
    '-slug' | 'winner_id' | '-winner_id' | 'winner_type' | '-winner_type';
}

type UpcomingTournamentsRequest = BaseTournamentsRequest;
type RunningTournamentsRequest = BaseTournamentsRequest;

export {
    UpcomingTournamentsRequest,
    RunningTournamentsRequest
};