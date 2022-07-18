import { NodeGroup } from "./node";

export enum LeagueRegion {
    Unset = "UNSET",
    NorthAmerica = "NA",
    SouthAmerica = "SA",
    Europe = "EUROPE",
    CIS = "CIS",
    SouthEastAsia = "SEA"
};

export enum LeagueTier {
    Unset = "UNSET",
    Amateur = "AMATEUR",
    Professional = "PROFESSIONAL",
    Minor = "MINOR",
    Major = "MAJOR",
    International = "INTERNATIONAL",
    DPCQualifier = "DPC_QUALIFIER",
    DPCLeagueQualifier = "DPC_LEAGUE_QUALIFIER",
    DPCLeague = "DPC_LEAGUE",
    DPCLeagueFinals = "DPC_LEAGUE_FINALS"
};

export type League = {
    id: number;
    name: string;
    displayName: string;
    region: LeagueRegion;
    tier: LeagueTier;
    description: string;
    nodeGroups: NodeGroup[];
};
