import { IDotaBot } from "../interfaces/discord/bot";
import { ITournamentsAPI } from "../interfaces/tournaments/api";

class BotController {
    _bot: IDotaBot;
    _api: ITournamentsAPI;

    constructor(bot: IDotaBot, api: ITournamentsAPI) {
        this._bot = bot;
        this._api = api;
    }

    initialise = () => {
        
    }
}

export default BotController;