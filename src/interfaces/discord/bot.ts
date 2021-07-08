interface IDotaBot {
    initialise: (readyCallback: () => void) => void;
    postTournaments: () => void;
}

export {
    IDotaBot
};