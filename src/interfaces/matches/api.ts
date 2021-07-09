import { RunningMatchesRequest, UpcomingMatchesRequest } from "./requests";
import { UpcomingMatchesResponse, RunningMatchesResponse } from "./responses";

interface IMatchesAPI {
    getRunningMatches: (request: RunningMatchesRequest) => Promise<RunningMatchesResponse>;
    getUpcomingMatches: (request: UpcomingMatchesRequest) => Promise<UpcomingMatchesResponse>;
}

export {
    IMatchesAPI
};