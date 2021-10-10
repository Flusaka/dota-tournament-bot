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

    updateChannelConfiguration(channelId: string, config: ChannelConfig) {
        this.channelsRef.child(channelId).update(config);
    }

    removeChannelConfiguration(channelId: string) {
        throw new Error("Method not implemented.");
    }
}

export default FirebaseDatabaseConnector;