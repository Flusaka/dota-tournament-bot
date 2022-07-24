import ChannelConfig from "../models/channel_models";

interface IDatabaseConnector {
    addChannelConfiguration(channelId: string, config: Partial<ChannelConfig>): void;
    getAllChannelConfigurations(): Promise<Map<string, ChannelConfig>>
    updateChannelConfiguration(channelId: string, config: Partial<ChannelConfig>): void;
    removeChannelConfiguration(channelId: string): void;
}

export default IDatabaseConnector;