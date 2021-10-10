import ChannelConfig from "../models/channel_models";

interface IDatabaseConnector {
    addChannelConfiguration(channelId: string, config: ChannelConfig): void;
    getAllChannelConfigurations(): Promise<Map<string, ChannelConfig>>
    updateChannelConfiguration(channelId: string, config: ChannelConfig): void;
    removeChannelConfiguration(channelId: string): void;
}

export default IDatabaseConnector;