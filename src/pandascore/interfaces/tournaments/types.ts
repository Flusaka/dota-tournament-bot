import { BaseLeague } from '../leagues/types';
import { BaseMatch } from '../matches/types';
import { BaseSerie } from '../series/types';
import { BaseTeam } from '../teams/types';
import { PlayerType } from '../common/player_type';
import { VideoGame, VideoGameVersion } from '../videogames/types';

type BaseTournament = {
    readonly begin_at: string | null;
    readonly end_at: string | null;
    readonly id: number;
    readonly league_id: number;
    readonly live_supported: boolean;
    readonly modified_at: string | null;
    readonly name: string;
    readonly prizepool: string | null;
    readonly serie_id: number;
    readonly slug: string;
    readonly winner_id: number | null;
    readonly winner_type: PlayerType;
}

type Tournament = BaseTournament & {
    readonly league: BaseLeague;
    readonly matches: BaseMatch[];
    readonly serie: BaseSerie;
    readonly teams: BaseTeam[];
    readonly videogame: VideoGame;
    readonly videogame_version: VideoGameVersion;
}

export {
    BaseTournament,
    Tournament
};