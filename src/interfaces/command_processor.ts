import { Message } from 'discord.js';

interface ICommandProcessor {
    shouldProcess: (message: Message) => boolean;
    processCommand: (message: Message) => void;
}

export { ICommandProcessor };