import { IDotaBot } from "../interfaces/discord/bot";
import { ITournamentsAPI } from "../interfaces/tournaments/api";

class BotController {
    private bot: IDotaBot;
    private api: ITournamentsAPI;

    constructor(bot: IDotaBot, api: ITournamentsAPI) {
        this.bot = bot;
        this.api = api;
    }

    initialise = async () => {
        // this.bot.initialise(async () => {
        //     const runningTournaments = await this.api.getRunningTournaments({
        //         sort: 'begin_at'
        //     });
        //     runningTournaments.filter(tournament => tournament.serie);
        //     this.bot.postTournaments();
        // });
        const runningTournaments = await this.api.getRunningTournaments({
            sort: 'begin_at'
        });

        const highTierTournaments = runningTournaments.filter(t => t.serie.tier == 'b' || t.serie.tier == 'a')

        console.log(highTierTournaments.map(t => t.teams.map(team => team.name)));
    }
}

export default BotController;