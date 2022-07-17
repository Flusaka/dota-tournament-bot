export type League = {
    id: number;
    displayName: string;
    region: "UNSET" | "NA" | "SA" | "EUROPE" | "CIS" | "CHINA" | "SEA";
    startDateTime: number;
    endDateTime: number;
    description: string;
};