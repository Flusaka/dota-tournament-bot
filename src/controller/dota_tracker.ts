import { TextChannel } from 'discord.js';
import moment, { unitOfTime } from 'moment-timezone';
import { IMatchesAPI } from "../pandascore/interfaces/matches/api";
import { ITournamentsAPI } from "../pandascore/interfaces/tournaments/api";
import MessageSender from "./message_sender";
import { DailyMatchesMessage } from "./messages";
import { RunningTournamentsResponse } from "../pandascore/interfaces/tournaments/responses";
import IDatabaseConnector from "../database/interfaces/database_connector";
import ChannelConfig from "../database/models/channel_models";

type TimerRef = ReturnType<typeof setTimeout>;

class DotaTracker {
    private channelId: string;

    private matchesApi: IMatchesAPI;

    private tournamentsApi: ITournamentsAPI;

    private databaseConnector: IDatabaseConnector;

    private messageSender: MessageSender;

    private dailyNotificationTime: moment.Moment;

    private dailyNotificationRef: TimerRef;

    private config: ChannelConfig;

    constructor(channel: TextChannel, matchesApi: IMatchesAPI, tournamentsApi: ITournamentsAPI, databaseConnector: IDatabaseConnector) {
        this.channelId = channel.id;
        this.matchesApi = matchesApi;
        this.tournamentsApi = tournamentsApi;
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

        this.messageSender.send(`:robot: Daily notifications of games will occur at: ${hour > 12 ? hour - 12 : hour}:${minutes < 10 ? `0${minutes}` : minutes} ${hour < 12 ? 'AM' : 'PM'}`);
        this._setDailyNotificationTime(hour, minutes);
    }

    setTimeZone = (timeZone: string) => {
        // TODO: Re-evaluate notification times, since we've potentially moved some hours ahead or behind...
        this.databaseConnector.updateChannelConfiguration(this.channelId, {
            timeZone: timeZone
        });

        this.messageSender.send(`:robot: Timezone is now set to: ${timeZone}`);
    }

    _setDailyNotificationTime = (hour: number, minutes: number) => {
        const notificationTime = moment.tz(this.config.timeZone);
        notificationTime.set("hour", hour);
        notificationTime.set("minutes", minutes);
        notificationTime.set("seconds", 0);
        notificationTime.set("milliseconds", 0);

        this._setDailyNotificationTimeMoment(notificationTime);
    }

    _setDailyNotificationTimeMoment = (notificationTime: moment.Moment) => {
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

    _setDailyNotificationTimeout = (timeout: number) => {
        if (this.dailyNotificationRef !== null) {
            console.log("Clearing existing notification timeout");
            clearTimeout(this.dailyNotificationRef);
            this.dailyNotificationRef = null;
        }

        console.log(`Timeout callback will fire: ${timeout}`);
        this.dailyNotificationRef = setTimeout(this._postDailyNotification, timeout);
    }

    _postDailyNotification = () => {
        // Get the list of running tournaments
        Promise.all<RunningTournamentsResponse>([
            this.tournamentsApi.getRunningTournaments({
                sort: '-end_at'
            }),
            this.tournamentsApi.getUpcomingTournaments({
                sort: 'begin_at'
            })
        ]).then((upcomingTournaments) => {
            const flattenedTournaments = upcomingTournaments.flat();
            const beginningOfDay = new Date();
            beginningOfDay.setHours(0, 0, 0);

            const endOfDay = new Date();
            endOfDay.setHours(23, 59, 59);

            const filteredTournaments = flattenedTournaments
                .filter(tournament => tournament.serie.tier == 'a' || tournament.serie.tier == 's')
                .filter(tournament => new Date(tournament.end_at) >= beginningOfDay);

            const tournamentMessages: DailyMatchesMessage[] = filteredTournaments.map((tournament) => {
                const filteredMatches = tournament.matches.filter(match => {
                    return new Date(match.begin_at) <= endOfDay && (match.end_at === null || new Date(match.end_at) >= beginningOfDay);
                });

                return {
                    tournamentName: `${tournament.league.name} - ${tournament.name}`,
                    matches: filteredMatches.map((match) => {
                        // Try and get the first official, main and english stream
                        let stream = match.streams_list.find(stream => stream.language === "en" && stream.official && stream.main);
                        if (stream === null) {
                            // If no stream can be found, find the first official && main stream
                            stream = match.streams_list.find(stream => stream.official && stream.main);

                            if (stream === null) {
                                // If no stream can be found still, find the first official stream
                                stream = match.streams_list.find(stream => stream.official);

                                if (stream === null) {
                                    // If _STILL_ no stream can be found, get the first english stream...
                                    stream = match.streams_list.find(stream => stream.language === "en");

                                    if (stream === null) {
                                        // If it's STILL STILL null, just accept the first one...
                                        stream = match.streams_list[0];
                                    }
                                }
                            }
                        }

                        return {
                            matchId: match.id,
                            matchTitle: match.name,
                            streamLink: stream ? stream.raw_url : "Unknown :person_shrugging:",
                            startTime: new Date(match.begin_at)
                        }
                    })
                };
            });

            if (tournamentMessages.length > 0) {
                this.messageSender.postDailyMatches(tournamentMessages);
            }
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