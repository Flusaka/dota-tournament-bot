import { IMatchesAPI } from "../pandascore/interfaces/matches/api";
import { ITournamentsAPI } from "../pandascore/interfaces/tournaments/api";
import MessageSender from "./message_sender";
import { TextChannel } from 'discord.js';
import { DailyMatchesMessage } from "./messages";

type TimerRef = ReturnType<typeof setTimeout>;

class DotaTracker {
    private messageSender: MessageSender;

    private dailyNotificationTime: Date;

    private matchesApi: IMatchesAPI;

    private tournamentsApi: ITournamentsAPI;

    private dailyNotificationRef: TimerRef;

    private notifications: Map<string, TimerRef>;

    constructor(channel: TextChannel, matchesApi: IMatchesAPI, tournamentsApi: ITournamentsAPI) {
        this.messageSender = new MessageSender(channel);
        this.matchesApi = matchesApi;
        this.tournamentsApi = tournamentsApi;
        this.notifications = new Map<string, TimerRef>();
        this.dailyNotificationTime = null;
        this.dailyNotificationRef = null;
    }

    setDailyNotificationTime = (dateTime: Date) => {
        // TODO: Cancel existing notification if there is one...
        this.dailyNotificationTime = new Date(dateTime.getTime());

        // If it's at a time before now, add a day to the specified time
        if (dateTime.getTime() < Date.now()) {
            console.log("Time is before current time, adding a day!")
            this.dailyNotificationTime.setDate(dateTime.getDate() + 1);
        }

        console.log(`Setting daily notification time to ${this.dailyNotificationTime.toString()}`)

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
        this.tournamentsApi.getRunningTournaments({
            sort: '-end_at'
        }).then((upcomingTournaments) => {
            const beginningOfDay = new Date();
            beginningOfDay.setHours(0, 0, 0);

            const endOfDay = new Date();
            endOfDay.setHours(23, 59, 59);

            const filteredTournaments = upcomingTournaments
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
                        console.log(match.streams_list);
                        if (stream === null) {
                            // If no stream can be found, find the first official && main stream
                            stream = match.streams_list.find(stream => stream.official && stream.main);

                            if(stream === null) {
                                // If no stream can be found still, find the first official stream
                                stream = match.streams_list.find(stream => stream.official);

                                if(stream === null) {
                                    // If _STILL_ no stream can be found, get the first english stream...
                                    stream = match.streams_list.find(stream => stream.language === "en");

                                    if(stream === null) {
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
            console.log(`Something went wrong when retrieving upcoming matches... ${error}`);
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