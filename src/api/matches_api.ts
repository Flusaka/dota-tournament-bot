import axios from 'axios';
import { IMatchesAPI } from "../interfaces/matches/api";
import { RunningMatchesRequest, UpcomingMatchesRequest } from "../interfaces/matches/requests";
import { RunningMatchesResponse, UpcomingMatchesResponse } from "../interfaces/matches/responses";

class MatchesAPI implements IMatchesAPI {
    private readonly axiosInstance = axios.create({
        baseURL: 'https://api.pandascore.co/dota2/matches',
        headers: {
            common: {
                // TODO: Move into env variable
                Authorization: 'Bearer 8FG9WnjcQBp9FkS8PA6bTQAEKYQefsBhWBjOG_hC7VYu4vWLxNM'
            }
        },
        responseType: 'json'
    });

    getRunningMatches = async (request: RunningMatchesRequest): Promise<RunningMatchesResponse> => {
        const response = await this.axiosInstance.get<RunningMatchesResponse>('/running', {
            params: request
        });
        return response.data;
    }

    getUpcomingMatches = async (request: UpcomingMatchesRequest): Promise<UpcomingMatchesResponse> => {
        const response = await this.axiosInstance.get<UpcomingMatchesResponse>('/upcoming', {
            params: request
        });
        return response.data;
    }
}

export default MatchesAPI;