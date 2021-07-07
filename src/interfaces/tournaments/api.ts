import * as TournamentRequests from './requests';
import * as TournamentResponses from './responses';

interface ITournamentsAPI {
    getUpcomingTournaments: (request: TournamentRequests.UpcomingTournamentsRequest) => Promise<TournamentResponses.UpcomingTournamentsResponse> | null;
    getRunningTournaments: (request: TournamentRequests.RunningTournamentsRequest) => Promise<TournamentResponses.RunningTournamentsResponse> | null;
}

export {
    ITournamentsAPI
};