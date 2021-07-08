import { IDotaBot } from "../interfaces/discord/bot";
import { ITournamentsAPI } from "../interfaces/tournaments/api";

class BotController {
    private bot: IDotaBot;
    private api: ITournamentsAPI;

    constructor(bot: IDotaBot, api: ITournamentsAPI) {
        this.bot = bot;
        this.api = api;
    }

    initialise = () => {
        this.bot.initialise(() => {
            this.bot.postTournaments();
        });
    }
}

export default BotController;