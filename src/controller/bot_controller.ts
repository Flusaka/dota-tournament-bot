import Discord, { Message, TextChannel } from 'discord.js';
import { Command, CommandProcessor } from "./command_processor";
import { DotaTracker } from './dota_tracker';
// import MatchesTestAPI from '../test/api/matches_test_api';
import { IMatchesAPI } from '../pandascore/interfaces/matches/api';
import { ITournamentsAPI } from '../pandascore/interfaces/tournaments/api';
// import TournamentsTestAPI from '../test/api/tournaments_test_api';
import MatchesAPI from '../pandascore/api/matches_api';
import TournamentsAPI from '../pandascore/api/tournaments_api';
import IDatabaseConnector from '../database/interfaces/database_connector';
import ChannelConfig, { DefaultChannelConfig } from '../database/models/channel_models';


class BotController {
    private client: Discord.Client;

    private commandProcessor: CommandProcessor;

    private dotaTrackers: Map<string, DotaTracker>;

    private matchesApi: IMatchesAPI;

    private tournamentsApi: ITournamentsAPI;

    private databaseConnector: IDatabaseConnector;

    constructor(databaseConnector: IDatabaseConnector) {
        // Discord client
        this.client = new Discord.Client();

        // Setup command processor
        this.commandProcessor = new CommandProcessor();
        this.commandProcessor.registerCallback(Command.EnableBotInChannel, this._enableBot);
        this.commandProcessor.registerCallback(Command.DisableBotInChannel, this._disableBot);
        this.commandProcessor.registerCallback(Command.SetDailyTime, this._setDailyTime, 1);
        this.commandProcessor.registerCallback(Command.SetTimeZone, this._setTimeZone, 1);
        // this.commandProcessor.registerCallback(Command.Notify, this.notifyUser);

        // Dota trackers map
        this.dotaTrackers = new Map<string, DotaTracker>();

        // API handlers
        this.matchesApi = new MatchesAPI();
        this.tournamentsApi = new TournamentsAPI();

        // Database connector
        this.databaseConnector = databaseConnector;

        this.client.login(process.env.DISCORD_TOKEN);
    }

    initialise = () => {
        this.client.on('ready', async () => {
            console.log("Ready");

            const channelConfigs = await this.databaseConnector.getAllChannelConfigurations();
            this._initialiseTrackers(channelConfigs);
        });

        this.client.on('message', this._messageReceived);
    }

    _initialiseTrackers = (channelConfigs: Map<string, ChannelConfig>) => {
        channelConfigs.forEach((config, channelId) => {
            const discordChannel = this.client.channels.cache.get(channelId) as TextChannel;
            const tracker = new DotaTracker(discordChannel, this.matchesApi, this.tournamentsApi, this.databaseConnector);
            tracker.setup(config);

            this.dotaTrackers.set(channelId, tracker);
        });
    }

    _messageReceived = (message: Message) => {
        if (message.author.bot) {
            return;
        }

        if (this.commandProcessor.shouldProcess(message)) {
            this.commandProcessor.processCommand(message);
        }
    }

    _enableBot = (message: Message, parameters: string[]) => {
        this._getDotaTrackerForChannel(message.channel.id, (exists, tracker) => {
            if (!exists) {
                console.log(`enable bot: ${message.channel.id}`);
                message.channel.send(":robot: Dota Bot enabled!");

                const tracker = new DotaTracker(message.channel as TextChannel, this.matchesApi, this.tournamentsApi, this.databaseConnector);
                tracker.setup(DefaultChannelConfig);
                this.dotaTrackers.set(message.channel.id, tracker);
                this.databaseConnector.addChannelConfiguration(message.channel.id, DefaultChannelConfig);
            }
        })
    }

    _disableBot = (message: Message, parameters: string[]) => {
        console.log(`disable bot: ${message.channel.id}`);
        message.channel.send(":robot: Dota Bot disabled!");

        this._getDotaTrackerForChannel(message.channel.id, (exists, tracker) => {
            if (exists) {
                tracker.shutdown();
                this.dotaTrackers.delete(message.channel.id);
                this.databaseConnector.removeChannelConfiguration(message.channel.id);
            }
        });
    }

    _setDailyTime = (message: Message, parameters: string[]) => {
        function parseTime(timeString: string): [number, number] {
            if (timeString == '') return null;

            var time = timeString.match(/(\d+)(:(\d\d))?\s*(p?)/i);
            if (time == null) return [-1, -1];

            var hours = parseInt(time[1], 10);
            if (hours == 12 && !time[4]) {
                hours = 0;
            }
            else {
                hours += (hours < 12 && time[4]) ? 12 : 0;
            }

            return [hours, parseInt(time[3], 10) || 0]
        }

        this._getDotaTrackerForChannel(message.channel.id, (exists, tracker) => {
            if (!exists) {
                message.channel.send("You need to enable the bot on this channel! Please type \"!dotabot start\" first!")
            }
            else {
                // Parse time
                const [hours, minutes] = parseTime(parameters[0]);
                if (hours < 0 || minutes < 0) {
                    message.channel.send("Please enter a time in the correct format! e.g. 5PM, 5:00PM, 17:00 etc.")
                }
                else {
                    tracker.setDailyNotificationTime(hours, minutes);
                }
            }
        });
    }

    _setTimeZone = (message: Message, parameters: string[]) => {
        this._getDotaTrackerForChannel(message.channel.id, (exists, tracker) => {
            if (!exists) {
                message.channel.send("You need to enable the bot on this channel! Please type \"!dotabot start\" first!")
            }
            else {
                // Pass to tracker
                tracker.setTimeZone(parameters[0]);
            }
        })
    }

    notifyUser = (message: Message, parameters: string[]) => {
    }

    _getDotaTrackerForChannel = (channelId: string, callback: (exists: boolean, tracker: DotaTracker) => void) => {
        if (this.dotaTrackers.has(channelId)) {
            callback(true, this.dotaTrackers.get(channelId));
            return;
        }

        callback(false, null);
    }
}

export default BotController;