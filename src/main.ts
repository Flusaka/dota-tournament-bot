// import express from 'express';
// import Discord from 'discord.js';

// const client = new Discord.Client();

// client.on('ready', () => {
//     console.log(`Logged in as ${client.user.tag}`);
// });

// client.login('ODYyMzMyNzY3MDIyNjEyNTIx.YOWz-Q.fvj0mW-pFY3349Qe8A9YRrKZfIw');

// const app = express();
// const port = 3000;

// app.listen(port, () => {
//     return console.log(`Server is listening on ${port}`);
// });

import TournamentsAPI from "./api/tournaments_api";

const tournamentsApi = new TournamentsAPI();
tournamentsApi.getRunningTournaments({
    sort: 'begin_at'
});
