import IDotaAPIClient from "../interfaces/api_client";
import fetch from 'cross-fetch';
import { League } from "../models/league";
import { LEAGUES_QUERY } from "./queries";

class DotaGraphQLClient implements IDotaAPIClient {
    API_URL = 'https://api.stratz.com/graphql';

    getMatchesToday(): Promise<League[]> {
        return this._query(LEAGUES_QUERY);
    }

    private _query<DataT>(query: string, variables?: object): Promise<DataT> {
        return new Promise((resolve, reject) => {
            const body = JSON.stringify({
                query,
                variables
            });

            fetch(this.API_URL, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': `Bearer ${process.env.STRATZ_TOKEN}`
                },
                body
            }).then(result => result.json()).then(result => {
                return resolve(result.data);
            }).catch(error => {
                return reject(error);
            });
        });
    }
}

export default DotaGraphQLClient;