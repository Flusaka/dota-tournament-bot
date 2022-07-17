import { League } from "../models/league";

interface IDotaAPIClient {
    getActiveLeagues(): Promise<League[]>;
}

export default IDotaAPIClient;