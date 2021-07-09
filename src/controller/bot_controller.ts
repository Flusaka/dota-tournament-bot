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

        let upcomingMatches = await this.matchesApi.getUpcomingMatches({
            sort: 'begin_at'
        });

        upcomingMatches = upcomingMatches.filter(match => (match.serie.tier == 'b' || match.serie.tier == 'a') && match.opponents.length > 0);

        console.log(upcomingMatches[0].opponents);
    }
}

export default BotController;