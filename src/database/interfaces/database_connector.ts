import ChannelConfig from "../models/channel_models";

interface IDatabaseConnector {
    addChannelConfiguration(channelId: string, config: ChannelConfig);
    updateChannelConfiguration(channelId: string, config: ChannelConfig);
    removeChannelConfiguration(channelId: string);
}

export default IDatabaseConnector;