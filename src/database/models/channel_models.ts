enum Tournament {
    International = "international",
    Major = "major",
    Minor = "minor",
    DotaProCircuit = "dpc",
    DotaProCircuitQualifiers = "dpc_qualifiers",
    OtherPro = "other_pro"
}

type ChannelConfig = {
    timeZone?: string;
    dailyNotificationHour?: number;
    dailyNotificationMinute?: number;
    tournaments: Tournament[];
}

const defaultChannelConfig: ChannelConfig = {
    timeZone: "Europe/London",
    tournaments: [Tournament.DotaProCircuit, Tournament.Major, Tournament.International]
}

export { defaultChannelConfig as DefaultChannelConfig };
export default ChannelConfig;