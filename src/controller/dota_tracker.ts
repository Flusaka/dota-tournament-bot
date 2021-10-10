import { IMatchesAPI } from "../pandascore/interfaces/matches/api";
import { ITournamentsAPI } from "../pandascore/interfaces/tournaments/api";
import MessageSender from "./message_sender";
import { TextChannel } from 'discord.js';
import { DailyMatchesMessage } from "./messages";
import { RunningTournamentsResponse } from "../pandascore/interfaces/tournaments/responses";
import IDatabaseConnector from "../database/interfaces/database_connector";

type TimerRef = ReturnType<typeof setTimeout>;

class DotaTracker {
    private channelId: string;

    private matchesApi: IMatchesAPI;

    private tournamentsApi: ITournamentsAPI;

    private databaseConnector: IDatabaseConnector;

    private messageSender: MessageSender;

    private dailyNotificationTime: Date;

    private dailyNotificationRef: TimerRef;

    constructor(channel: TextChannel, matchesApi: IMatchesAPI, tournamentsApi: ITournamentsAPI, databaseConnector: IDatabaseConnector) {
        this.channelId = channel.id;
        this.matchesApi = matchesApi;
        this.tournamentsApi = tournamentsApi;
        this.messageSender = new MessageSender(channel);
        this.databaseConnector = databaseConnector;

        // TODO: Load these from stored ChannelConfig eventually
        this.dailyNotificationTime = null;
        this.dailyNotificationRef = null;
    }

    shutdown = () => {
        // Clear any timeout refs
        if (this.dailyNotificationRef !== null) {
            clearTimeout(this.dailyNotificationRef);
        }
    }

    setDailyNotificationTime = (dateTime: Date) => {
        // TODO: Use the ChannelConfig to base the datetime off the stored timezone

        this.dailyNotificationTime = new Date(dateTime.getTime());

        // If it's at a time before now, add a day to the specified time
        if (dateTime.getTime() < Date.now()) {
            console.log("Time is before current time, adding a day!")
            this.dailyNotificationTime.setDate(dateTime.getDate() + 1);
        }

        console.log(`Setting daily notification time to ${this.dailyNotificationTime.toString()}`)

        // Store the next time to the database
        this.databaseConnector.updateChannelConfiguration(this.channelId, {
            dailyNotificationTime: this.dailyNotificationTime
        });

        if (this.dailyNotificationRef !== null) {
            console.log("Clearing existing notification timeout");
            clearTimeout(this.dailyNotificationRef);
            this.dailyNotificationRef = null;
        }

        const timeout = this.dailyNotificationTime.getTime() - Date.now();
        this.dailyNotificationRef = setTimeout(this.postDailyNotification, timeout);
    }

    postDailyNotification = () => {
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
            const nextNotificationTime = new Date(this.dailyNotificationTime.getTime());
            // nextNotificationTime.setMinutes(nextNotificationTime.getMinutes() + 1);
            nextNotificationTime.setDate(nextNotificationTime.getDate() + 1);
            this.setDailyNotificationTime(nextNotificationTime);
        });
    }
}

export { DotaTracker };