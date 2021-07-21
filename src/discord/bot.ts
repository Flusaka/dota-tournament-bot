import { IDotaBot } from "../interfaces/discord/bot";
import { DailyMatchesMessage } from '../interfaces/discord/messages';
import Discord, { TextChannel } from 'discord.js';

class DotaBot implements IDotaBot {
    private client: Discord.Client;

    private numberMessages = ['\:one:', '\:two:', '\:three:', '\:four:', '\:five:', '\:six:', '\:seven:', '\:eight:', '\:nine:', '\:keycap_ten:'];
    private numberReactions = ['1️⃣', '2️⃣', '3️⃣', '4️⃣', '5️⃣', '6️⃣', '7️⃣', '8️⃣', '9️⃣', '🔟'];

    initialise = (readyCallback: () => void) => {
        this.client = new Discord.Client();
        // TODO: Move into env variable
        this.client.login('ODYyMzMyNzY3MDIyNjEyNTIx.YOWz-Q.fvj0mW-pFY3349Qe8A9YRrKZfIw');

        this.client.on('ready', () => {
            readyCallback();
        });
    }

    postDailyMatches = async (message: DailyMatchesMessage) => {
        const guild = this.client.guilds.cache.get("328238345098625024");
        const channel = guild.channels.cache.get("555859747547512832") as TextChannel;

        const matchesText: string = message.matches.splice(0, 9).map((match, index) => {
            return `${this.numberMessages[index]} ${match.startTime} - ${match.matchTitle}`;
        }).join('\n');

        const sentMessage = await channel.send(`${message.tournamentName} matches today!\n` +
            `Games on ${message.matches[0].streamLink}:\n` +
            `${matchesText}\n` +
            `React with match number for notification when the game should be starting!`);

        sentMessage.react(this.numberReactions[0]);
    }
}

export default DotaBot;