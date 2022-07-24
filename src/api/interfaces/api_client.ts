import { League, LeagueTier } from "../models/league";

interface IDotaAPIClient {
    getMatchesToday(tiers: LeagueTier[]): Promise<League[]>;
}

export default IDotaAPIClient;