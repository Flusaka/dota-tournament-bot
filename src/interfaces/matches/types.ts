import { League } from "../leagues/types";
import { Serie } from "../series/types";
import { Stream } from "../streams/types";
import { Tournament } from "../tournaments/types";

enum OpponentType {
    Player,
    Team
}

type OpponentDetails = {
    readonly acronym: string | null;
    readonly id: number;
    readonly image_url: string | null;
    readonly location: string | null;
    readonly modified_at: string;
    readonly name: string;
    readonly slug: string | null;
}

type Opponent = {
    readonly opponent: OpponentDetails;
    readonly type: OpponentType;
}

type BaseMatch = {
    readonly begin_at: string | null;
    readonly detailed_stats: boolean;
    readonly draw: boolean;
    readonly end_at: string | null;
    readonly forfeit: boolean;
    readonly game_advantage: number | null;
    readonly games: any;
    readonly id: number;
    readonly league: League;
    readonly league_id: number;
    readonly live: any;
    readonly match_type: 'best_of' | 'custom' | 'first_to' | 'ow_best_of';
    readonly modified_at: string;
    readonly name: string;
    readonly number_of_games: number;
    readonly original_scheduled_at: string | null;
    readonly rescheduled: boolean;
    readonly scheduled_at: string | null;
    readonly serie: Serie;
    readonly slug: string;
    readonly status: 'canceled' | 'finished' | 'not_started' | 'postponed' | 'running';
    readonly streams_list: Stream[];
    readonly tournament_id: number;
    readonly winner_id: number | null;
}

type Match = BaseMatch & {
    readonly opponents: Opponent[];
    readonly tournament: Tournament;
}

export {
    BaseMatch,
    Match
};