import IDotaAPIClient from "../interfaces/api_client";
import fetch from 'cross-fetch';
import { ApolloClient, createHttpLink, gql } from '@apollo/client/core';
import { InMemoryCache, NormalizedCacheObject } from '@apollo/client/cache';
import { setContext } from "@apollo/client/link/context"

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

    getUpcomingMatches() {
        this.client.query({
            query: gql`
            query {
                live {
                    matches {
                        matchId
                    }
                }
            }
            `
        }).then(result => {
            console.log("Received result!");
        }).catch(error => {
            console.error(error.networkError.response);
        });
    }
}

export default DotaGraphQLClient;