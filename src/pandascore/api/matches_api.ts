import axios from 'axios';
import { IMatchesAPI } from "../interfaces/matches/api";
import { PastMatchesRequest, RunningMatchesRequest, UpcomingMatchesRequest } from "../interfaces/matches/requests";
import { PastMatchesResponse, RunningMatchesResponse, UpcomingMatchesResponse } from "../interfaces/matches/responses";

class MatchesAPI implements IMatchesAPI {
    private readonly axiosInstance = axios.create({
        baseURL: 'https://api.pandascore.co/dota2/matches',
        headers: {
            common: {
                // TODO: Move into env variable
                Authorization: `Bearer ${process.env.PS_TOKEN}`
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

    getPastMatches = async (request: PastMatchesRequest): Promise<PastMatchesResponse> => {
        const response = await this.axiosInstance.get<UpcomingMatchesResponse>('/past', {
            params: request
        });
        return response.data;
    }
}

export default MatchesAPI;