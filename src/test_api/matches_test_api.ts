// import axios from 'axios';
import { IMatchesAPI } from "../pandascore/interfaces/matches/api";
import { PastMatchesRequest, RunningMatchesRequest, UpcomingMatchesRequest } from "../pandascore/interfaces/matches/requests";
import { PastMatchesResponse, RunningMatchesResponse, UpcomingMatchesResponse } from "../pandascore/interfaces/matches/responses";
import fs from 'fs';

class MatchesTestAPI implements IMatchesAPI {
    // private readonly axiosInstance = axios.create({
    //     baseURL: 'https://api.pandascore.co/dota2/matches',
    //     headers: {
    //         common: {
    //             // TODO: Move into env variable
    //             Authorization: 'Bearer 8FG9WnjcQBp9FkS8PA6bTQAEKYQefsBhWBjOG_hC7VYu4vWLxNM'
    //         }
    //     },
    //     responseType: 'json'
    // });

    getRunningMatches = (request: RunningMatchesRequest): Promise<RunningMatchesResponse> => {
        // const response = await this.axiosInstance.get<RunningMatchesResponse>('/running', {
        //     params: request
        // });
        // return response.data;
        const matches = fs.readFileSync("matches.json");
        const json = matches.toString();
        const response: RunningMatchesResponse = JSON.parse(json);
        return Promise.resolve(response);
    }

    getUpcomingMatches = (request: UpcomingMatchesRequest): Promise<UpcomingMatchesResponse> => {
        // const response = await this.axiosInstance.get<UpcomingMatchesResponse>('/upcoming', {
        //     params: request
        // });
        // return response.data;
        const matches = fs.readFileSync("matches.json");
        const json = matches.toString();
        const response: UpcomingMatchesResponse = JSON.parse(json);
        return Promise.resolve(response);
    }

    getPastMatches = (request: PastMatchesRequest): Promise<PastMatchesResponse> => {
        // const response = await this.axiosInstance.get<UpcomingMatchesResponse>('/past', {
        //     params: request
        // });
        // return response.data;
        const matches = fs.readFileSync("matches.json");
        const json = matches.toString();
        const response: PastMatchesResponse = JSON.parse(json);
        return Promise.resolve(response);
    }
}

export default MatchesTestAPI;