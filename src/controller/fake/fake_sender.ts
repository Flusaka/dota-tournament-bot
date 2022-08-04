import IMessageSender from "../interfaces/message_sender";
import { DailyMatchesMessage, MatchDetails } from "../types/messages";

export default class FakeSender implements IMessageSender {
    postDailyMatches(messages: DailyMatchesMessage[]): void {
        messages.forEach(m => this.postDailyMatch(m));
    }

    postDailyMatch(message: DailyMatchesMessage): void {
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

        console.log(matchesText);
    }

    send(message: string): void {
        console.log(message);
    }

}