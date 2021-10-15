type ChannelConfig = {
    timeZone?: string;
    dailyNotificationHour?: number;
    dailyNotificationMinute?: number;
}

const defaultChannelConfig: ChannelConfig = {
    timeZone: "Europe/London"
}

export { defaultChannelConfig as DefaultChannelConfig };
export default ChannelConfig;