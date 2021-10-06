import { ITournamentsAPI } from '../../pandascore/interfaces/tournaments/api';
import { PastTournamentsRequest, RunningTournamentsRequest, UpcomingTournamentsRequest } from '../../pandascore/interfaces/tournaments/requests';
import { PastTournamentsResponse, RunningTournamentsResponse, UpcomingTournamentsResponse } from '../../pandascore/interfaces/tournaments/responses';
import fs from 'fs';
import { Tournament } from '../../pandascore/interfaces/tournaments/types';

class TournamentsTestAPI implements ITournamentsAPI {
    getUpcomingTournaments = async (request: UpcomingTournamentsRequest): Promise<UpcomingTournamentsResponse> => {
        const matches = fs.readFileSync("tournaments.json");
        const json = matches.toString();
        const response: UpcomingTournamentsResponse = JSON.parse(json);
        const altered: Tournament = {
            ...response[0],
            begin_at: new Date(Date.now()),
            end_at: new Date(Date.now() + 5 * 1000 * 60),
            matches: response[0].matches.map((match, index) => {
                const fakeStart = index * 5 * 1000 * 60;
                const fakeEnd = fakeStart + (5 * 1000 * 60);
                return {
                    ...match,
                    // begin_at: new Date(Date.now() + fakeStart),
                    // end_at: new Date(Date.now() + fakeEnd)
                };
            })
        }
        return Promise.resolve([altered]);
    }

    getRunningTournaments = async (request: RunningTournamentsRequest): Promise<RunningTournamentsResponse> => {
        const matches = fs.readFileSync("tournaments.json");
        const json = matches.toString();
        const response: RunningTournamentsResponse = JSON.parse(json);
        return Promise.resolve(response);
    }

    getPastTournaments = async (request: PastTournamentsRequest): Promise<PastTournamentsResponse> => {
        const matches = fs.readFileSync("tournaments.json");
        const json = matches.toString();
        const response: PastTournamentsResponse = JSON.parse(json);
        return Promise.resolve(response);
    }
}

export default TournamentsTestAPI;