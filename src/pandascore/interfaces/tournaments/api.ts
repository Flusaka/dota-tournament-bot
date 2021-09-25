import * as TournamentRequests from './requests';
import * as TournamentResponses from './responses';

interface ITournamentsAPI {
    getUpcomingTournaments: (request: TournamentRequests.UpcomingTournamentsRequest) => Promise<TournamentResponses.UpcomingTournamentsResponse>;
    getRunningTournaments: (request: TournamentRequests.RunningTournamentsRequest) => Promise<TournamentResponses.RunningTournamentsResponse>;
    getPastTournaments: (request: TournamentRequests.PastTournamentsRequest) => Promise<TournamentResponses.PastTournamentsResponse>;
}

export {
    ITournamentsAPI
};