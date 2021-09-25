import axios from 'axios';
import { ITournamentsAPI } from '../interfaces/tournaments/api';
import { PastTournamentsRequest, RunningTournamentsRequest, UpcomingTournamentsRequest } from '../interfaces/tournaments/requests';
import { PastTournamentsResponse, RunningTournamentsResponse, UpcomingTournamentsResponse } from '../interfaces/tournaments/responses';
// import { PastMatchesRequest, RunningMatchesRequest, UpcomingMatchesRequest } from "../interfaces/matches/requests";
// import { PastMatchesResponse, RunningMatchesResponse, UpcomingMatchesResponse } from "../interfaces/matches/responses";

class TournamentsAPI implements ITournamentsAPI {
    private readonly axiosInstance = axios.create({
        baseURL: 'https://api.pandascore.co/dota2/tournaments',
        headers: {
            common: {
                // TODO: Move into env variable
                Authorization: 'Bearer 8FG9WnjcQBp9FkS8PA6bTQAEKYQefsBhWBjOG_hC7VYu4vWLxNM'
            }
        },
        responseType: 'json'
    });

    getUpcomingTournaments = async (request: UpcomingTournamentsRequest): Promise<UpcomingTournamentsResponse> => {
        const response = await this.axiosInstance.get<UpcomingTournamentsResponse>('/upcoming', {
            params: request
        });
        return response.data;
    }

    getRunningTournaments = async (request: RunningTournamentsRequest): Promise<RunningTournamentsResponse> => {
        const response = await this.axiosInstance.get<RunningTournamentsResponse>('/running', {
            params: request
        });
        return response.data;
    }

    getPastTournaments = async (request: PastTournamentsRequest): Promise<PastTournamentsResponse> => {
        const response = await this.axiosInstance.get<PastTournamentsResponse>('/past', {
            params: request
        });
        return response.data;
    }
}

export default TournamentsAPI;