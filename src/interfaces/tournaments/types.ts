import { BaseMatch } from "../matches/types";

type BaseTournament = {
    readonly begin_at: string | null;
    readonly end_at: string | null;
    readonly id: number;
    readonly league: any;
    readonly league_id: number;
    readonly live_supported: boolean;
    readonly matches: BaseMatch[];
    readonly modified_at: string | null;
    readonly name: string;
    readonly prizepool: string | null;
    readonly serie: any;
    readonly serie_id: number;
    readonly slug: string;
    readonly teams: any[];
    readonly videogame: any;
    readonly winner_id: number | null;
    readonly winner_type: string | null;
}

export {
    BaseTournament
};