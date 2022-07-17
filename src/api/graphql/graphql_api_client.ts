import IDotaAPIClient from "../interfaces/api_client";
import fetch from 'cross-fetch';
import { ApolloClient, createHttpLink, gql } from '@apollo/client/core';
import { InMemoryCache, NormalizedCacheObject } from '@apollo/client/cache';
import { setContext } from "@apollo/client/link/context"
import { League } from "../models/league";

class DotaGraphQLClient implements IDotaAPIClient {
    private client: ApolloClient<NormalizedCacheObject>;

    constructor() {
        console.log(process.env.STRATZ_TOKEN);
        const authLink = setContext((_, { headers }) => {
            return {
                headers: {
                    ...headers,
                    authorization: `Bearer ${process.env.STRATZ_TOKEN}`
                }
            }
        });

        const httpLink = createHttpLink({
            uri: "https://api.stratz.com/graphql",
            fetch
        });

        this.client = new ApolloClient<NormalizedCacheObject>({
            cache: new InMemoryCache(),
            link: authLink.concat(httpLink)
        });
    }

    getActiveLeagues(): Promise<League[]> {
        return new Promise(async (resolve, reject) => {
            try {
                const leagues = await this.client.query<League[]>({
                    query: gql`
                    query {
                        leagues(request: {tiers: DPC_LEAGUE, leagueEnded: false}) {
                            id
                            displayName
                            region
                            startDateTime
                            endDateTime
                            description
                        }
                    }
                    `
                });

                return resolve(leagues.data);
            }
            catch (error) {
                return reject(error);
            }
        });
    }
}

export default DotaGraphQLClient;