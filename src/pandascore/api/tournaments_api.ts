import axios from 'axios';
import { ITournamentsAPI } from '../interfaces/tournaments/api';
import { PastTournamentsRequest, RunningTournamentsRequest, UpcomingTournamentsRequest } from '../interfaces/tournaments/requests';
import { PastTournamentsResponse, RunningTournamentsResponse, UpcomingTournamentsResponse } from '../interfaces/tournaments/responses';

class TournamentsAPI implements ITournamentsAPI {
    private readonly axiosInstance = axios.create({
        baseURL: 'https://api.pandascore.co/dota2/tournaments',
        headers: {
            common: {
                Authorization: `Bearer ${process.env.PS_TOKEN}`
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