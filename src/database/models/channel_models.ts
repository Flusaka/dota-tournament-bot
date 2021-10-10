type ChannelConfig = {
    timeZone?: string;
    dailyNotificationTime?: Date;
}

const defaultChannelConfig: ChannelConfig = {
    timeZone: "Europe/London"
}

export { defaultChannelConfig as DefaultChannelConfig };
export default ChannelConfig;