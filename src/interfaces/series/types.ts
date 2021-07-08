type Serie = {
    readonly begin_at: Date | null;
    readonly description: string | null;
    readonly end_at: Date | null;
    readonly full_name: string;
    readonly id: number;
    readonly league: any;
    readonly league_id: number;
    readonly modified_at: Date;
    readonly name: string | null;
    readonly season: string | null;
    readonly slug: string;
    readonly tier: string | null;
    readonly tournaments: any;
    readonly videogame: any;
    readonly videogame_title: any | null;
    readonly winner_id: number | null;
    readonly winner_type: string | null;
    readonly year: number;
}

export {
    Serie
};