type VideoGameVersion = {
    readonly current: boolean;
    readonly name: string;
}

type VideoGame = {
    readonly id: number;
    readonly name: string;
    readonly slug: string;
}

export {
    VideoGame,
    VideoGameVersion
};