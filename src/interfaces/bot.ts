import { DailyMatchesMessage } from './messages';

interface IDotaBot {
    initialise: (readyCallback: () => void) => void;
    postDailyMatches: (messages: DailyMatchesMessage[]) => void;
    postDailyMatch: (message: DailyMatchesMessage) => void;
}

export {
    IDotaBot
};