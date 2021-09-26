import { IMatchesAPI } from "../../pandascore/interfaces/matches/api";
import { PastMatchesRequest, RunningMatchesRequest, UpcomingMatchesRequest } from "../../pandascore/interfaces/matches/requests";
import { PastMatchesResponse, RunningMatchesResponse, UpcomingMatchesResponse } from "../../pandascore/interfaces/matches/responses";
import fs from 'fs';

class MatchesTestAPI implements IMatchesAPI {
    getRunningMatches = (request: RunningMatchesRequest): Promise<RunningMatchesResponse> => {
        const matches = fs.readFileSync("matches.json");
        const json = matches.toString();
        const response: RunningMatchesResponse = JSON.parse(json);
        return Promise.resolve(response);
    }

    getUpcomingMatches = (request: UpcomingMatchesRequest): Promise<UpcomingMatchesResponse> => {
        const matches = fs.readFileSync("matches.json");
        const json = matches.toString();
        const response: UpcomingMatchesResponse = JSON.parse(json);
        return Promise.resolve(response);
    }

    getPastMatches = (request: PastMatchesRequest): Promise<PastMatchesResponse> => {
        const matches = fs.readFileSync("matches.json");
        const json = matches.toString();
        const response: PastMatchesResponse = JSON.parse(json);
        return Promise.resolve(response);
    }
}

export default MatchesTestAPI;