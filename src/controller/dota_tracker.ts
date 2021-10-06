import { IMatchesAPI } from "../pandascore/interfaces/matches/api";
import { ITournamentsAPI } from "../pandascore/interfaces/tournaments/api";
import { TextChannel, User } from 'discord.js';

type TimerRef = ReturnType<typeof setTimeout>;

class DotaTracker {
    private channel: TextChannel;

    private dailyNotificationTime: Date;

    private matchesApi: IMatchesAPI;

    private tournamentsApi: ITournamentsAPI;

    private notifications: Map<string, TimerRef>;

    constructor(channel: TextChannel, matchesApi: IMatchesAPI, tournamentsApi: ITournamentsAPI) {
        this.channel = channel;
        this.matchesApi = matchesApi;
        this.tournamentsApi = tournamentsApi;
        this.notifications = new Map<string, TimerRef>();
    }

    setDailyNotificationTime = (dateTime: Date) => {
        if (dateTime.getUTCDate() < Date.now()) {

        }
    }

    registerNotification = (user: User, timeout: number) => {
        setTimeout(() => {
            this.channel.send(`${user}, you have been notified!`);
        }, timeout * 1000);
    }
}

export { DotaTracker };