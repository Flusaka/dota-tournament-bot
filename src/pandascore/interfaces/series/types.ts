import { PlayerType } from "../common/player_type";
import { BaseLeague } from "../leagues/types";
import { BaseTournament } from "../tournaments/types";

type BaseSerie = {
    readonly begin_at: Date | null;
    readonly description: string | null;
    readonly end_at: Date | null;
    readonly full_name: string;
    readonly id: number;
    readonly league_id: number;
    readonly modified_at: Date | null;
    readonly name: string | null;
    readonly season: string | null;
    readonly slug: string;
    readonly tier: string | null;
    readonly winner_id: number | null;
    readonly winner_type: PlayerType | null;
    readonly year: number;
}

type Serie = BaseSerie & {
    readonly league: BaseLeague;
    readonly tournament: BaseTournament;
}

export {
    BaseSerie,
    Serie
};