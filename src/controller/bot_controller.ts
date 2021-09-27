import { IDotaBot } from "../interfaces/bot";
import { DailyMatchesMessage, MatchDetails } from "../interfaces/messages";
import { IMatchesAPI } from "../pandascore/interfaces/matches/api";
import fs from 'fs';
import { ITournamentsAPI } from "../pandascore/interfaces/tournaments/api";
import { Command, CommandProcessor } from "../discord/command_processor";

class BotController {
    private bot: IDotaBot;
    private matchesApi: IMatchesAPI;
    private tournamentsApi: ITournamentsAPI;
    private commandProcessor: CommandProcessor;

    constructor(bot: IDotaBot, matchesApi: IMatchesAPI, tournamentsApi: ITournamentsAPI) {
        this.bot = bot;
        this.matchesApi = matchesApi;
        this.tournamentsApi = tournamentsApi;

        // Setup command processor
        this.commandProcessor = new CommandProcessor();
        this.commandProcessor.registerCallback(Command.EnableBotInChannel, this.enableBot);
        this.commandProcessor.registerCallback(Command.DisableBotInChannel, this.disableBot);
    }

    initialise = async () => {
        this.bot.initialise(this.commandProcessor, async () => {
            console.log("Bot is ready");
        });
    }

    enableBot = () => {
        console.log("enable bot");
    }

    disableBot = () => {
        console.log("disable bot");
    }
}

export default BotController;