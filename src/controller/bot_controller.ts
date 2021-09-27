import { IDotaBot } from "../interfaces/bot";
import { DailyMatchesMessage, MatchDetails } from "../interfaces/messages";
import { IMatchesAPI } from "../pandascore/interfaces/matches/api";
import fs from 'fs';
import { ITournamentsAPI } from "../pandascore/interfaces/tournaments/api";
import { Command, CommandProcessor } from "./command_processor";

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
        // this.commandProcessor.registerCallback
    }

    initialise = async () => {
        this.bot.initialise(this.commandProcessor, async () => {
            console.log("Bot is ready");
        });
    }

    enableBot = (userId: string, channelId: string, parameters: string[]) => {
        // TODO: Register channel with 
        console.log(`enable bot: ${userId} ${channelId}`);
    }

    disableBot = (userId: string, channelId: string, parameters: string[]) => {
        console.log(`disable bot: ${userId} ${channelId}`);
    }
}

export default BotController;