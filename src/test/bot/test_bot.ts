import { IDotaBot } from "../../interfaces/bot";
import { DailyMatchesMessage } from '../../interfaces/messages';

class TestDotaBot implements IDotaBot {
    initialise = (readyCallback: () => void) => {
        readyCallback();
    }

    postDailyMatches = async (messages: DailyMatchesMessage[]) => {
        messages.forEach(message => this.postDailyMatch(message));
    }

    postDailyMatch = async (message: DailyMatchesMessage) => {
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

        console.log(`${message.tournamentName} matches today!\n` +
            `Games on ${message.matches[0].streamLink}:\n` +
            `${matchesText}`
        );
    }
}

export default TestDotaBot;