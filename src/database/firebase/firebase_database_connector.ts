import IDatabaseConnector from "../interfaces/database_connector";
import ChannelConfig from "../models/channel_models";
import admin from 'firebase-admin';

class FirebaseDatabaseConnector implements IDatabaseConnector {
    private channelsRef: admin.database.Reference;

    constructor() {
        this.channelsRef = admin.database().ref("channels");
    }

    addChannelConfiguration(channelId: string, config: ChannelConfig) {
        this.channelsRef.child(channelId).set(config);
    }

    async getAllChannelConfigurations(): Promise<Map<string, ChannelConfig>> {
        try {
            const channelsSnapshot = await this.channelsRef.get();
            if (!channelsSnapshot.exists()) {
                return Promise.reject("channels property does not exist");
            }

            const channelConfigs = new Map<string, ChannelConfig>();
            channelsSnapshot.forEach(channelSnapshot => {
                channelConfigs.set(channelSnapshot.key, channelSnapshot.val())
            });

            return await Promise.resolve(channelConfigs);
        }
        catch (error) {
            return Promise.reject(error);
        }
    }

    updateChannelConfiguration(channelId: string, config: ChannelConfig) {
        this.channelsRef.child(channelId).update(config);
    }

    removeChannelConfiguration(channelId: string) {
        throw new Error("Method not implemented.");
    }
}

export default FirebaseDatabaseConnector;