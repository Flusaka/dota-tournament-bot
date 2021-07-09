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

import BotController from './controller/bot_controller';
import TournamentsAPI from './api/tournaments_api';
import DotaBot from './discord/bot';
import MatchesAPI from './api/matches_api';

const botController = new BotController(new DotaBot(), new TournamentsAPI(), new MatchesAPI());
botController.initialise();