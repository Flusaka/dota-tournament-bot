import IDotaAPIClient from "../interfaces/api_client";
import fetch from 'cross-fetch';
import { League, LeagueTier } from "../models/league";
import { LEAGUES_QUERY } from "./queries";

class DotaGraphQLClient implements IDotaAPIClient {
    API_URL = 'https://api.stratz.com/graphql';
    API_HEADERS = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': `Bearer ${process.env.STRATZ_TOKEN}`
    };

    getMatchesToday(tiers: LeagueTier[]): Promise<League[]> {
        return new Promise((resolve, reject) => {
            this._query(LEAGUES_QUERY, {
                tiers,
                leagueEnded: false
            }).then(result => {
                resolve(result['leagues']);
            })
                .catch(error => reject(error));
        });
    }

    private _query(query: string, variables?: object): Promise<object> {
        return new Promise((resolve, reject) => {
            const body = JSON.stringify({
                query,
                variables
            });

            fetch(this.API_URL, {
                method: 'POST',
                headers: this.API_HEADERS,
                body
            }).then(result => result.json()).then(result => {
                if (result.errors) {
                    return reject(result.errors);
                }
                return resolve(result.data);
            }).catch(error => {
                return reject(error);
            });
        });
    }
}

export default DotaGraphQLClient;