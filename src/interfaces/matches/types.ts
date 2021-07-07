type BaseMatch = {
    readonly begin_at: Date | null;
    readonly detailed_stats: boolean;
    readonly draw: boolean;
    readonly end_at: Date | null;
    readonly forfeit: boolean;
    readonly game_advantage: number | null;
    readonly games: any;
    readonly id: number;
    readonly league: any;
    readonly league_id: number;
    readonly live: any;
    readonly match_type: 'best_of' | 'custom' | 'first_to' | 'ow_best_of';
    readonly modified_at: Date;
    readonly name: string;
    readonly number_of_games: number;
    readonly original_scheduled_at: Date | null;
    readonly rescheduled: boolean;
    readonly scheduled_at: Date | null;
    readonly slug: string;
    readonly status: 'canceled' | 'finished' | 'not_started' | 'postponed' | 'running';
    readonly streams_list: any;
    readonly tournament_id: number;
    readonly winner_id: number | null;
}

export {
    BaseMatch
};