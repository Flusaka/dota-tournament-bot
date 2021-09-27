import { ICommandProcessor } from './command_processor';
import { DailyMatchesMessage } from './messages';

interface IDotaBot {
    initialise: (commandProcessor: ICommandProcessor, readyCallback: () => void) => void;
    postDailyMatches: (messages: DailyMatchesMessage[]) => void;
    postDailyMatch: (message: DailyMatchesMessage) => void;
}

export {
    IDotaBot
};