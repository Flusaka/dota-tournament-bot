import Discord, { Message, TextChannel } from 'discord.js';
import { DailyMatchesMessage, MatchDetails } from "../interfaces/messages";
import { Command, CommandProcessor } from "./command_processor";
import { DotaTracker } from './dota_tracker';
import MatchesAPI from '../pandascore/api/matches_api';
import TournamentsAPI from '../pandascore/api/tournaments_api';
import MatchesTestAPI from '../test/api/matches_test_api';
import { IMatchesAPI } from '../pandascore/interfaces/matches/api';
import { ITournamentsAPI } from '../pandascore/interfaces/tournaments/api';
import TournamentsTestAPI from '../test/api/tournaments_test_api';


class BotController {
    private client: Discord.Client;

    private commandProcessor: CommandProcessor;

    private dotaTrackers: Map<string, DotaTracker>;

    private matchesApi: IMatchesAPI;

    private tournamentsApi: ITournamentsAPI;

    constructor() {
        // Discord client
        this.client = new Discord.Client();

        // Setup command processor
        this.commandProcessor = new CommandProcessor();
        this.commandProcessor.registerCallback(Command.EnableBotInChannel, this.enableBot);
        this.commandProcessor.registerCallback(Command.DisableBotInChannel, this.disableBot);
        this.commandProcessor.registerCallback(Command.Notify, this.notifyUser);

        // Dota trackers map
        this.dotaTrackers = new Map<string, DotaTracker>();

        // API handlers
        this.matchesApi = new MatchesTestAPI();
        this.tournamentsApi = new TournamentsTestAPI();

        // TODO: Move into env variable
        this.client.login('ODYyMzMyNzY3MDIyNjEyNTIx.YOWz-Q.fvj0mW-pFY3349Qe8A9YRrKZfIw');
    }

    initialise = async () => {
        this.client.on('ready', () => {
            console.log("Ready");
        });

        this.client.on('message', this.messageReceived);
    }

    messageReceived = (message: Message) => {
        if (message.author.bot) {
            return;
        }

        if (this.commandProcessor.shouldProcess(message)) {
            this.commandProcessor.processCommand(message);
        }
    }

    enableBot = (message: Message, parameters: string[]) => {
        console.log(`enable bot: ${message.channel.id}`);

        if (!this.dotaTrackers.has(message.channel.id)) {
            this.dotaTrackers.set(message.channel.id, new DotaTracker(message.channel as TextChannel, this.matchesApi, this.tournamentsApi));
        }
    }

    disableBot = (message: Message, parameters: string[]) => {
        console.log(`disable bot: ${message.channel.id}`);

        // TODO: Clear out any existing timeouts etc. properly shutdown the tracker
        this.dotaTrackers.delete(message.channel.id);
    }

    notifyUser = (message: Message, parameters: string[]) => {
        if (this.dotaTrackers.has(message.channel.id)) {
            this.dotaTrackers.get(message.channel.id).registerNotification(message.author, 5);
        }
    }
}

export default BotController;