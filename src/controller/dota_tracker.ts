import { TextChannel } from 'discord.js';
import moment from 'moment-timezone';
import MessageSender from "./message_sender";
import { DailyMatchesMessage, MatchDetails } from "./messages";
import IDotaAPIClient from '../api/interfaces/api_client';
import IDatabaseConnector from "../database/interfaces/database_connector";
import ChannelConfig from "../database/models/channel_models";
import { LeagueTier } from '../api/models/league';

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
        this.dotaApiClient.getMatchesToday([LeagueTier.DPCLeague, LeagueTier.International, LeagueTier.Major]).then((leagues) => {
            // Beginning of the "day" is at the time of the notifications to go out
            const beginningOfDay = moment.tz(this.dailyNotificationTime, this.config.timeZone);
            // End of the "day" is at the time the notifications will go out tomorrow
            const endOfDay = moment(beginningOfDay);
            endOfDay.add(1, 'day');

            // Find all leagues with matches today
            const leaguesWithMatchesToday = leagues.map(league => {
                const nodeGroups = league.nodeGroups.map(group => {
                    return {
                        ...group,
                        nodes: group.nodes.filter(node => {
                            if (node.hasStarted || node.isCompleted) {
                                return false;
                            }

                            const scheduledTime = moment.tz(node.scheduledTime, this.config.timeZone);
                            return scheduledTime >= beginningOfDay && scheduledTime < endOfDay;
                        })
                    };
                }).filter(group => group.nodes.length > 0);

                return {
                    ...league,
                    nodeGroups
                };
            }).filter(league => league.nodeGroups.length > 0);

            // Now go through each league and post the daily messages, split based on streams (i.e. node groups)
            const messages = leaguesWithMatchesToday.map(league => {
                const matches = league.nodeGroups.map(group => {
                    const matchDetails: MatchDetails[] = group.nodes.map(node => {
                        return {
                            matchId: node.id,
                            matchTitle: `${node.teamOne.name} vs ${node.teamTwo.name}`,
                            startTime: moment.tz(node.scheduledTime, this.config.timeZone),
                            streamLink: node.streams[0].streamUrl
                        }
                    });
                    return matchDetails;
                }).flat().sort((a, b) => a.startTime.diff(b.startTime));

                const dailyMessage: DailyMatchesMessage = {
                    leagueName: league.displayName,
                    matches
                };

                return dailyMessage;
            });

            this.messageSender.postDailyMatches(messages);
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