interface ICommandProcessor {
    shouldProcess: (message: string) => boolean;
    processCommand: (message: string) => void;
}

export { ICommandProcessor };