import { TextChannel } from 'discord.js';
import IMessageSender from './interfaces/message_sender';
import { DailyMatchesMessage, MatchDetails } from './types/messages';

class MessageSender implements IMessageSender {
    private channel: TextChannel;

    constructor(channel: TextChannel) {
        this.channel = channel;
    }

    postDailyMatches = async (messages: DailyMatchesMessage[]) => {
        messages.forEach(message => this.postDailyMatch(message));
    }

    postDailyMatch = async (message: DailyMatchesMessage) => {
        if (message.matches.length === 0) {
            return;
        }

        const organisedMatches: Array<MatchDetails[]> = [];
        while (message.matches.length > 0) {
            const streamLink = message.matches[0].streamLink;
            const matchesOnStream = message.matches.filter(match => match.streamLink === streamLink);
            organisedMatches.push(matchesOnStream);
            matchesOnStream.forEach((match) => {
                message.matches.splice(message.matches.indexOf(match), 1);
            });
        }

        const matchesText: string = organisedMatches.map((matches) => {
            let matchesSetText = `Games on: ${matches[0].streamLink}\n\n`;
            matchesSetText = matchesSetText.concat(matches.map(match => {
                const startTime = match.startTime.format("h:mm A")
                return `${startTime} - ${match.matchTitle}`;
            }).join('\n'));
            return matchesSetText;
        }).join('\n\n');

        const discordMessage = await this.channel.send(`:robot: **${message.leagueName} matches today!**\n` +
            `${matchesText}`
        );
        discordMessage.suppressEmbeds(true);
    }

    send = (message: string) => {
        this.channel.send(message);
    }
}

export default MessageSender;