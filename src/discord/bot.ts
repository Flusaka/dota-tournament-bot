import { IDotaBot } from "../interfaces/discord/bot";
import { DailyMatchesMessage } from '../interfaces/discord/messages';
import Discord, { TextChannel } from 'discord.js';

class DotaBot implements IDotaBot {
    private client: Discord.Client;

    initialise = (readyCallback: () => void) => {
        this.client = new Discord.Client();
        // TODO: Move into env variable
        this.client.login('ODYyMzMyNzY3MDIyNjEyNTIx.YOWz-Q.fvj0mW-pFY3349Qe8A9YRrKZfIw');

        this.client.on('ready', () => {
            readyCallback();
        });
    }

    postDailyMatches = (message: DailyMatchesMessage) => {
        // const guild = this.client.guilds.cache.get("328238345098625024");
        // const channel = guild.channels.cache.get("555859747547512832") as TextChannel;
        for (let i = 0; i < message.matches.length; ++i) {
            console.log(message.matches[i]);
        }
    }
}

export default DotaBot;