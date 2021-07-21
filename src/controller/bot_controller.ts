import { IDotaBot } from "../interfaces/discord/bot";
import { DailyMatchesMessage, Match } from "../interfaces/discord/messages";
import { IMatchesAPI } from "../interfaces/matches/api";
import { ITournamentsAPI } from "../interfaces/tournaments/api";

class BotController {
    private bot: IDotaBot;
    private tournamentsApi: ITournamentsAPI;
    private matchesApi: IMatchesAPI;

    constructor(bot: IDotaBot, tournamentsApi: ITournamentsAPI, matchesApi: IMatchesAPI) {
        this.bot = bot;
        this.tournamentsApi = tournamentsApi;
        this.matchesApi = matchesApi;
    }

    initialise = async () => {
        // this.bot.initialise(async () => {
        //     const runningTournaments = await this.api.getRunningTournaments({
        //         sort: 'begin_at'
        //     });
        //     runningTournaments.filter(tournament => tournament.serie);
        //     this.bot.postTournaments();
        // });
        // const runningTournaments = await this.tournamentsApi.getRunningTournaments({
        //     sort: 'begin_at'
        // });

        // const highTierTournaments = runningTournaments.filter(t => t.serie.tier == 'b' || t.serie.tier == 'a');
        // const matches: Match[] = highTierTournaments.map((tournament): Match[] => {
        //     return tournament.matches.map((match): Match => {
        //         return {
        //             gameName: match.name,
        //             startTime: match.begin_at,
        //             streamLink: match.streams_list.length > 0 ? match.streams_list[0].raw_url : ""
        //         };
        //     });
        // }).flat();

        // this.bot.postDailyMatches({
        //     matches
        // });

        this.bot.initialise(async () => {
            let upcomingMatches = await this.matchesApi.getPastMatches({
                sort: '-scheduled_at'
            });
            upcomingMatches = upcomingMatches.filter(match => (match.serie.tier == 'b' || match.serie.tier == 'a'))
                .sort((a, b) => a.tournament_id - b.tournament_id);

            while (upcomingMatches.length > 0) {
                const firstMatch = upcomingMatches[0];
                const matchesInTournament = upcomingMatches.filter(match => match.tournament_id == firstMatch.tournament_id);
                upcomingMatches.splice(0, matchesInTournament.length);

                this.bot.postDailyMatches({
                    tournamentName: `${firstMatch.league.name} - ${firstMatch.tournament.name}`,
                    matches: matchesInTournament.map(match => {
                        let streamLink = "";
                        if (match.streams_list.length > 0) {
                            const official = match.streams_list.filter(stream => stream.official);
                            if (official.length > 0) {
                                streamLink = official[0].raw_url;
                            }
                            else {
                                const main = match.streams_list.filter(stream => stream.main);
                                if (main.length > 0) {
                                    streamLink = main[0].raw_url;
                                }
                                else {
                                    streamLink = match.streams_list[0].raw_url;
                                }
                            }
                        }
                        return {
                            startTime: new Date(Date.parse(match.begin_at)),
                            matchTitle: match.name,
                            streamLink: streamLink
                        };
                    })
                });
                // console.log(`Removing ${matchesInTournament.length} matches. ${upcomingMatches.length} remaining!`);
            }
        });
    }
}

export default BotController;