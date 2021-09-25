import { PlayerType } from "../common/player_type";

enum GameStatus {
    NotStarted = "not_started",
    NotPlayed = "not_played",
    Running = "running",
    Finished = "finished"
}

type GameWinner = {
    readonly id: number;

    readonly type: PlayerType | string;
}

type Game = {
    /** When this game is due to start */
    readonly begin_at: Date | null;

    /** Whether game data is complete and won't change */
    readonly complete: boolean;

    /** Whether the game offers full stats */
    readonly detailed_stats: boolean;

    /** When the game is due to end */
    readonly end_at: Date | null;

    /** Whether the game is finished */
    readonly finished: boolean;

    /** Whether the game was forfeit */
    readonly forfeit: boolean;

    /** ID for the game */
    readonly id: number;

    /** Length of the game */
    readonly length: number | null;

    /** ID of the match this game is part of */
    readonly match_id: number;

    /** The position of this game within the series (1, 2, 3 etc.) */
    readonly position: number;

    /** Current status of the game */
    readonly status: GameStatus;

    /** URI to a video of the game */
    readonly video_url: string | null;

    /** The ID and type of the game winner */
    readonly winner: GameWinner;

    /** The type of the game winner (a player or team) */
    readonly winner_type: PlayerType;
};

export {
    Game,
    GameStatus,
    GameWinner
};