import { TextChannel } from 'discord.js';
import moment, { unitOfTime } from 'moment-timezone';
import MessageSender from "./message_sender";
import { DailyMatchesMessage } from "./messages";
import IDotaAPIClient from '../api/interfaces/api_client';
import IDatabaseConnector from "../database/interfaces/database_connector";
import ChannelConfig from "../database/models/channel_models";

type TimerRef = ReturnType<typeof setTimeout>;

class DotaTracker {
    private channelId: string;

    private dotaApiClient: IDotaAPIClient;

    private databaseConnector: IDatabaseConnector;

    private messageSender: MessageSender;

    private dailyNotificationTime: moment.Moment;

    private dailyNotificationRef: TimerRef;

    private config: ChannelConfig;

    constructor(channel: TextChannel, dotaApiClient: IDotaAPIClient, databaseConnector: IDatabaseConnector) {
        this.channelId = channel.id;
        this.dotaApiClient = dotaApiClient;
        this.messageSender = new MessageSender(channel);
        this.databaseConnector = databaseConnector;

        this.dailyNotificationTime = null;
        this.dailyNotificationRef = null;

        this.dotaApiClient.getMatchesToday().then(leagues => {
            console.log(leagues);
        }).catch(err => {
            console.error(err);
        });
    }

    setup = (config: ChannelConfig) => {
        this.config = { ...config };

        if (!isNaN(this.config.dailyNotificationHour) && !isNaN(this.config.dailyNotificationMinute)) {
            this._setDailyNotificationTime(this.config.dailyNotificationHour, this.config.dailyNotificationMinute);
        }
    }

    shutdown = () => {
        // Clear any timeout refs
        if (this.dailyNotificationRef !== null) {
            clearTimeout(this.dailyNotificationRef);
        }
    }

    setDailyNotificationTime = (hour: number, minutes: number) => {
        // Store the next time to the database
        this.databaseConnector.updateChannelConfiguration(this.channelId, {
            dailyNotificationHour: hour,
            dailyNotificationMinute: minutes
        });

        this.config.dailyNotificationHour = hour;
        this.config.dailyNotificationMinute = minutes;

        this.messageSender.send(`:robot: Daily notifications of games will occur at: ${hour > 12 ? hour - 12 : hour}:${minutes < 10 ? `0${minutes}` : minutes} ${hour < 12 ? 'AM' : 'PM'}`);
        this._setDailyNotificationTime(hour, minutes);
    }

    setTimeZone = (timeZone: string) => {
        // TODO: Re-evaluate notification times, since we've potentially moved some hours ahead or behind...
        this.databaseConnector.updateChannelConfiguration(this.channelId, {
            timeZone: timeZone
        });
        this.config.timeZone = timeZone;

        this.messageSender.send(`:robot: Timezone is now set to: ${timeZone}`);
    }

    private _setDailyNotificationTime = (hour: number, minutes: number) => {
        const notificationTime = moment.tz(this.config.timeZone);
        notificationTime.set("hour", hour);
        notificationTime.set("minutes", minutes);
        notificationTime.set("seconds", 0);
        notificationTime.set("milliseconds", 0);

        this._setDailyNotificationTimeMoment(notificationTime);
    }

    private _setDailyNotificationTimeMoment = (notificationTime: moment.Moment) => {
        const now = moment().tz(this.config.timeZone);

        // If it's at a time before now, add a day to the specified time
        if (notificationTime < now) {
            console.log("Time is before current time, adding a day!")
            notificationTime.add(1, "day");
        }

        console.log(`Setting daily notification time to ${notificationTime.toISOString()}`);
        this.dailyNotificationTime = notificationTime;

        const timeout = notificationTime.valueOf() - now.valueOf();
        this._setDailyNotificationTimeout(timeout);
    }

    private _setDailyNotificationTimeout = (timeout: number) => {
        if (this.dailyNotificationRef !== null) {
            console.log("Clearing existing notification timeout");
            clearTimeout(this.dailyNotificationRef);
            this.dailyNotificationRef = null;
        }

        console.log(`Timeout callback will fire: ${timeout}`);
        this.dailyNotificationRef = setTimeout(this._postDailyNotification, timeout);
    }

    private _postDailyNotification = () => {
        // Get the list of running tournaments
        this.dotaApiClient.getMatchesToday().then((leagues) => {
            const beginningOfDay = moment.tz(this.config.timeZone).startOf("day");
            const endOfDay = moment.tz(this.config.timeZone).endOf("day");
        }).catch((error) => {
            console.log(`Something went wrong when retrieving tournaments... ${error}`);
        }).finally(() => {
            // Setup next notification time to a day in the future
            const nextNotificationTime = moment.tz(this.dailyNotificationTime, this.config.timeZone);
            nextNotificationTime.add(1, "day");
            this._setDailyNotificationTimeMoment(nextNotificationTime);
        });
    }
}

export { DotaTracker };