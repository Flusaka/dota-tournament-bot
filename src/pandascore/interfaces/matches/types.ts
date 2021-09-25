import { PlayerType } from "../common/player_type";
import { Game } from "../games/types";
import { BaseLeague } from "../leagues/types";
import { BaseSerie } from "../series/types";
import { Stream } from "../streams/types";
import { BaseTournament } from "../tournaments/types";
import { VideoGame, VideoGameVersion } from "../videogames/types";

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
    readonly type: PlayerType;
}

enum MatchType {
    BestOf = "best_of",
    Custom = "custom",
    FirstTo = "first_to"
    // TODO: OwBestOf may be needed
}

enum MatchStatus {
    NotStarted = "not_started",
    Running = "running",
    Finished = "finished",
    Postponed = "postponed",
    Canceled = "canceled",
}

type BaseMatch = {
    /** When the entire match is due to start */
    readonly begin_at: string | null;

    /** Whether the match offers full stats */
    readonly detailed_stats: boolean;
    readonly draw: boolean;
    readonly end_at: Date | null;
    readonly forfeit: boolean;
    readonly game_advantage: number | null;
    readonly id: number;
    // TODO: Could add "live" property here...
    readonly live_embed_url: string | null;
    readonly match_type: MatchType;
    readonly modified_at: Date | null;
    readonly name: string;
    readonly number_of_games: number;
    readonly official_stream_url: string | null;
    readonly original_scheduled_at: Date | null;
    readonly rescheduled: boolean;
    readonly scheduled_at: Date | null;
    readonly slug: string;
    readonly status: MatchStatus;
    // TODO: Add "streams" object?
    readonly streams_list: Stream[];
    readonly tournament_id: number;
    readonly winner_id: number | null;
}

type Match = BaseMatch & {
    readonly games: Game[];
    readonly league: BaseLeague;
    readonly league_id: number;
    readonly opponents: Opponent[];
    // TODO: Figure out what this looks like...
    readonly results: any;
    readonly serie_id: number;
    readonly serie: BaseSerie;
    readonly tournament: BaseTournament;
    readonly videogame: VideoGame;
    readonly videogame_version: VideoGameVersion;
}

export {
    BaseMatch,
    Match
};