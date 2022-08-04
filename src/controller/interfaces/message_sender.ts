import { DailyMatchesMessage } from "../types/messages";

export default interface IMessageSender {
    postDailyMatches(messages: DailyMatchesMessage[]): void;
    postDailyMatch(message: DailyMatchesMessage): void;
    send(message: string): void;
}