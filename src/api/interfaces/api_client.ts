import { League } from "../models/league";

interface IDotaAPIClient {
    getMatchesToday(): Promise<League[]>;
}

export default IDotaAPIClient;