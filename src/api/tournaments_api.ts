import axios from 'axios';
import { ITournamentsAPI } from "../interfaces/tournaments/api";
import * as Requests from '../interfaces/tournaments/requests';
import * as Responses from '../interfaces/tournaments/responses';

class TournamentsAPI implements ITournamentsAPI {
    private readonly _axiosInstance = axios.create({
        baseURL: 'https://api.pandascore.co/dota2',
        headers: {
            common: {
                Authorization: 'Bearer 8FG9WnjcQBp9FkS8PA6bTQAEKYQefsBhWBjOG_hC7VYu4vWLxNM'
            }
        },
        responseType: 'json'
    });

    getUpcomingTournaments = async (request: Requests.UpcomingTournamentsRequest): Promise<Responses.UpcomingTournamentsResponse> | null => {
        try {
            const response = await this._axiosInstance.get<Responses.UpcomingTournamentsResponse>('/tournaments/upcoming', {
                params: request
            });
            return response.data;
        }
        catch (error) {
            console.log(error);
            return null;
        }
    }

    getRunningTournaments = async (request: Requests.RunningTournamentsRequest): Promise<Responses.RunningTournamentsResponse> | null => {
        try {
            const response = await this._axiosInstance.get<Responses.RunningTournamentsResponse>('/tournaments/running', {
                params: request
            });
            return response.data;
        }
        catch (error) {
            console.log(error);
            return null;
        }
    }
}

export default TournamentsAPI;