import { Message } from 'discord.js';

type ProcessorCallback = (message: Message, parameters: string[]) => void;

enum Command {
    Invalid = "invalid",
    EnableBotInChannel = "start",
    DisableBotInChannel = "stop",
    Notify = "notify"
}

type CommandOptions = {
    callback: ProcessorCallback;
    numParameters: number;
}

class CommandProcessor {
    private commandCallbacks: Map<Command, CommandOptions>;


    constructor() {
        this.commandCallbacks = new Map<Command, CommandOptions>();
    }

    registerCallback = (command: Command, callback: ProcessorCallback, expectedNumParameters: number = 0) => {
        this.registerCallbackWithOptions(command, {
            numParameters: expectedNumParameters,
            callback: callback
        })
    }

    registerCallbackWithOptions = (command: Command, options: CommandOptions) => {
        this.commandCallbacks.set(command, options);
    }

    shouldProcess = (message: Message): boolean => {
        if (!message.content.startsWith("!dotabot")) {
            return false;
        }

        const parts = message.content.split(" ");
        if (parts.length <= 1) {
            return false;
        }

        return true;
    }

    processCommand = (message: Message): void => {
        const splitMessage = message.content.split(" ");

        // [1] is the actual command
        const commandString = splitMessage[1];

        // Find the command's enum
        const command = this.getCommandEnum(commandString);
        if (command == Command.Invalid) {
            return;
        }

        // Check if we have a valid command callback setup
        if (!this.commandCallbacks.has(command)) {
            return;
        }

        // Get the definition/options for the command
        const commandDef = this.commandCallbacks.get(command);

        // Deduct 2 for !dotabot and the command string
        const numParameters = splitMessage.length - 2;

        // Check that num of parameters is valid
        if (numParameters != commandDef.numParameters) {
            console.log("invalid number of parameters");
            return;
        }

        // [2-length] are the parameters (if there are any)
        // TODO: Validate that there are the correct amount of parameters for the command (if any)
        let parameters: string[] = [];
        if (splitMessage.length > 2) {
            parameters = splitMessage.slice(2);
        }

        commandDef.callback(message, parameters);
    }

    getCommandEnum = (commandString: string): Command => {
        const keys = Object.keys(Command).filter(k => Command[k] == commandString);
        return keys.length > 0 ? Command[keys[0]] : Command.Invalid;
    }
}

export { CommandProcessor, Command };