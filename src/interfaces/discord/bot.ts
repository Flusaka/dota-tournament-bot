import { DailyMatchesMessage } from './messages';

interface IDotaBot {
    initialise: (readyCallback: () => void) => void;
    postDailyMatches: (message: DailyMatchesMessage) => void;
}

export {
    IDotaBot
};