import axios from 'axios';
import { ITournamentsAPI } from "../interfaces/tournaments/api";
import * as Requests from '../interfaces/tournaments/requests';
import * as Responses from '../interfaces/tournaments/responses';

class TournamentsAPI implements ITournamentsAPI {
    private readonly axiosInstance = axios.create({
        baseURL: 'https://api.pandascore.co/dota2',
        headers: {
            common: {
                Authorization: 'Bearer 8FG9WnjcQBp9FkS8PA6bTQAEKYQefsBhWBjOG_hC7VYu4vWLxNM'
            }
        },
        responseType: 'json'
    });

    getUpcomingTournaments = async (request: Requests.UpcomingTournamentsRequest): Promise<Responses.UpcomingTournamentsResponse> => {
        const response = await this.axiosInstance.get<Responses.UpcomingTournamentsResponse>('/tournaments/upcoming', {
            params: request
        });
        return response.data;
    }

    getRunningTournaments = async (request: Requests.RunningTournamentsRequest): Promise<Responses.RunningTournamentsResponse> => {
        const response = await this.axiosInstance.get<Responses.RunningTournamentsResponse>('/tournaments/running', {
            params: request
        });
        return response.data;
    }
}

export default TournamentsAPI;