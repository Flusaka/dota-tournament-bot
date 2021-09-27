import { IDotaBot } from "../interfaces/bot";
import { DailyMatchesMessage } from '../interfaces/messages';
import Discord, { TextChannel } from 'discord.js';
import { ICommandProcessor } from "../interfaces/command_processor";

class DotaBot implements IDotaBot {
    private client: Discord.Client;
    private commandProcessor: ICommandProcessor;

    initialise = (processor: ICommandProcessor, readyCallback: () => void) => {
        this.commandProcessor = processor;
        this.client = new Discord.Client();

        // TODO: Move into env variable
        this.client.login('ODYyMzMyNzY3MDIyNjEyNTIx.YOWz-Q.fvj0mW-pFY3349Qe8A9YRrKZfIw');

        this.client.on('ready', () => {
            readyCallback();
        });

        this.client.on('message', this.messageReceived);
    }

    postDailyMatches = async (messages: DailyMatchesMessage[]) => {
        messages.forEach(message => this.postDailyMatch(message));
    }

    postDailyMatch = async (message: DailyMatchesMessage) => {
        const guild = this.client.guilds.cache.get("328238345098625024");
        const channel = guild.channels.cache.get("555859747547512832") as TextChannel;

        const formatter = new Intl.DateTimeFormat('en-GB', {
            timeZone: "Europe/London",
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hourCycle: 'h12'
        });

        const matchesText: string = message.matches.splice(0, 9).map((match) => {
            const startTime = formatter.format(match.startTime);
            return `[${match.matchId}] ${startTime} - ${match.matchTitle}`;
        }).join('\n');

        channel.send(`${message.tournamentName} matches today!\n` +
            `Games on ${message.matches[0].streamLink}:\n` +
            `${matchesText}`
        );

        console.log(`${message.tournamentName} matches today!\n` +
            `Games on ${message.matches[0].streamLink}:\n` +
            `${matchesText}`
        );
    }

    messageReceived = (message: Discord.Message) => {
        if (message.author.bot) {
            return;
        }

        if (this.commandProcessor.shouldProcess(message.content)) {
            this.commandProcessor.processCommand(message.content);
        }
    }
}

export default DotaBot;