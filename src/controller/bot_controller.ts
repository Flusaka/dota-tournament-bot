import { IDotaBot } from "../interfaces/bot";
import { DailyMatchesMessage, MatchDetails } from "../interfaces/messages";
import { IMatchesAPI } from "../pandascore/interfaces/matches/api";
import fs from 'fs';
import { ITournamentsAPI } from "../pandascore/interfaces/tournaments/api";

class BotController {
    private bot: IDotaBot;
    private matchesApi: IMatchesAPI;
    private tournamentsApi: ITournamentsAPI;

    constructor(bot: IDotaBot, matchesApi: IMatchesAPI, tournamentsApi: ITournamentsAPI) {
        this.bot = bot;
        this.matchesApi = matchesApi;
        this.tournamentsApi = tournamentsApi;
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
            let upcomingTournaments = await this.tournamentsApi.getUpcomingTournaments({
                sort: '-begin_at'
            });

            // TODO: Split matches within tournament based on stream location
            const tournamentMessages: DailyMatchesMessage[] = upcomingTournaments.map((tournament) => {
                return {
                    tournamentName: tournament.league.name,
                    matches: tournament.matches.map((match) => {
                        let stream = match.streams_list.find(stream => stream.language === "en" || stream.official || stream.main);
                        if (stream === null) {
                            stream = match.streams_list[0];
                        }

                        return {
                            matchId: match.id,
                            matchTitle: match.name,
                            streamLink: stream.raw_url,
                            startTime: match.begin_at
                        }
                    })
                };
            });

            this.bot.postDailyMatches(tournamentMessages);

            // console.log(upcomingTournaments[0].begin_at);
            // console.log(upcomingTournaments[0].end_at);

            // upcomingTournaments[0].matches.forEach((match) => {
            //     console.log(`${match.name} - ${match.begin_at} to ${match.end_at}`);
            // });

            // upcomingTournaments = upcomingTournaments.filter(tournament => (tournament.serie.tier == 'b' || tournament.serie.tier == 'a'));
            // fs.writeFileSync("tournaments.json", JSON.stringify(upcomingTournaments, null, 2));
            // let upcomingMatches = await this.matchesApi.getUpcomingMatches({
            //     sort: '-scheduled_at'
            // });

            // upcomingMatches = upcomingMatches.filter(match => (match.serie.tier == 'c' || match.serie.tier == 'b' || match.serie.tier == 'a'))
            //     .sort((a, b) => a.tournament_id - b.tournament_id);


            // while (upcomingMatches.length > 0) {
            //     const firstMatch = upcomingMatches[0];
            //     const matchesInTournament = upcomingMatches.filter(match => match.tournament_id == firstMatch.tournament_id);
            //     upcomingMatches.splice(0, matchesInTournament.length);

            //     this.bot.postDailyMatches({
            //         tournamentName: `${firstMatch.league.name} - ${firstMatch.tournament.name}`,
            //         matches: matchesInTournament.map(match => {
            //             let streamLink = "";
            //             if (match.streams_list.length > 0) {
            //                 const official = match.streams_list.filter(stream => stream.official);
            //                 if (official.length > 0) {
            //                     streamLink = official[0].raw_url;
            //                 }
            //                 else {
            //                     const main = match.streams_list.filter(stream => stream.main);
            //                     if (main.length > 0) {
            //                         streamLink = main[0].raw_url;
            //                     }
            //                     else {
            //                         streamLink = match.streams_list[0].raw_url;
            //                     }
            //                 }
            //             }
            //             return {
            //                 startTime: new Date(Date.parse(match.begin_at)),
            //                 matchTitle: match.name,
            //                 streamLink: streamLink
            //             };
            //         })
            //     });
            // }
        });
    }
}

export default BotController;