import { BaseSerie } from "../series/types";
import { VideoGame } from "../videogames/types";

type BaseLeague = {
    readonly id: number;
    readonly image_url: string | null;
    readonly modified_at: string | null;
    readonly name: string;
    readonly slug: string;
    readonly url: string | null;
}

type League = BaseLeague & {
    readonly series: BaseSerie[];
    readonly videogame: VideoGame;
}

export {
    BaseLeague,
    League
};