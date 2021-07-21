import { League } from '../leagues/types';
import { BaseMatch } from '../matches/types';
import { Serie } from '../series/types';
import { Team } from '../teams/types';
import { VideoGame } from '../videogames/types';

type Tournament = {
    readonly begin_at: string | null;
    readonly end_at: string | null;
    readonly id: number;
    readonly league: League;
    readonly league_id: number;
    readonly live_supported: boolean;
    readonly matches: BaseMatch[];
    readonly modified_at: string | null;
    readonly name: string;
    readonly prizepool: string | null;
    readonly serie: Serie;
    readonly serie_id: number;
    readonly slug: string;
    readonly teams: Team[];
    readonly videogame: VideoGame;
    readonly winner_id: number | null;
    readonly winner_type: string | null;
}

export {
    Tournament
};