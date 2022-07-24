import IDotaAPIClient from "../interfaces/api_client";
import { LeagueTier, League } from "../models/league";
import data from './fake_data.json';

export default class FakeClient implements IDotaAPIClient {
    getMatchesToday(tiers: LeagueTier[]): Promise<League[]> {
        const leagues = (data as League[]).filter(league => tiers.includes(league.tier));
        return Promise.resolve(leagues);
    }
}