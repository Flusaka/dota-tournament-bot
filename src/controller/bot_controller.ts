import Discord, { Message, TextChannel } from 'discord.js';
import { Command, CommandProcessor } from "./command_processor";
import { DotaTracker } from './dota_tracker';
import MatchesTestAPI from '../test/api/matches_test_api';
import { IMatchesAPI } from '../pandascore/interfaces/matches/api';
import { ITournamentsAPI } from '../pandascore/interfaces/tournaments/api';
import TournamentsTestAPI from '../test/api/tournaments_test_api';
import MatchesAPI from '../pandascore/api/matches_api';
import TournamentsAPI from '../pandascore/api/tournaments_api';


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
        this.commandProcessor.registerCallback(Command.SetDailyTime, this.setDailyTime, 1);
        // this.commandProcessor.registerCallback(Command.Notify, this.notifyUser);

        // Dota trackers map
        this.dotaTrackers = new Map<string, DotaTracker>();

        // API handlers
        this.matchesApi = new MatchesAPI();
        this.tournamentsApi = new TournamentsAPI();

        // TODO: Move into env variable
        this.client.login(process.env.DISCORD_TOKEN);
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
        message.channel.send(":robot: Dota Bot enabled!");

        if (!this.dotaTrackers.has(message.channel.id)) {
            this.dotaTrackers.set(message.channel.id, new DotaTracker(message.channel as TextChannel, this.matchesApi, this.tournamentsApi));
        }
    }

    disableBot = (message: Message, parameters: string[]) => {
        console.log(`disable bot: ${message.channel.id}`);
        message.channel.send(":robot: Dota Bot disabled!");

        this.getDotaTrackerForChannel(message.channel.id, (exists, tracker) => {
            if (exists) {
                tracker.shutdown();
                this.dotaTrackers.delete(message.channel.id);
            }
        });
    }

    setDailyTime = (message: Message, parameters: string[]) => {
        function parseTime(timeString: string): Date {
            if (timeString == '') return null;

            var time = timeString.match(/(\d+)(:(\d\d))?\s*(p?)/i);
            if (time == null) return null;

            var hours = parseInt(time[1], 10);
            if (hours == 12 && !time[4]) {
                hours = 0;
            }
            else {
                hours += (hours < 12 && time[4]) ? 12 : 0;
            }
            var d = new Date();
            d.setHours(hours);
            d.setMinutes(parseInt(time[3], 10) || 0);
            d.setSeconds(0, 0);
            return d;
        }

        this.getDotaTrackerForChannel(message.channel.id, (exists, tracker) => {
            if (!exists) {
                message.channel.send("You need to enable the bot on this channel! Please type \"!dotabot start\" first!")
            }
            else {
                // Parse time
                const dailyTime = parseTime(parameters[0]);
                message.channel.send(`:robot: Daily notifications of games will occur at: ${dailyTime.toTimeString()}`)
                tracker.setDailyNotificationTime(dailyTime);
            }
        });
    }

    notifyUser = (message: Message, parameters: string[]) => {
    }

    getDotaTrackerForChannel = (channelId: string, callback: (exists: boolean, tracker: DotaTracker) => void) => {
        if (this.dotaTrackers.has(channelId)) {
            callback(true, this.dotaTrackers.get(channelId));
            return;
        }

        callback(false, null);
    }
}

export default BotController;