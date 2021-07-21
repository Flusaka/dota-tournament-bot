import { PastMatchesRequest, RunningMatchesRequest, UpcomingMatchesRequest } from "./requests";
import { UpcomingMatchesResponse, RunningMatchesResponse, PastMatchesResponse } from "./responses";

interface IMatchesAPI {
    getRunningMatches: (request: RunningMatchesRequest) => Promise<RunningMatchesResponse>;
    getUpcomingMatches: (request: UpcomingMatchesRequest) => Promise<UpcomingMatchesResponse>;
    getPastMatches: (request: PastMatchesRequest) => Promise<PastMatchesResponse>;
}

export {
    IMatchesAPI
};