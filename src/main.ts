import express from 'express';
const app = express();
const port = process.env.PORT || 3000;

app.listen(port, () => {
    return console.log(`Server is listening on ${port}`);
});

import BotController from './controller/bot_controller';

const botController = new BotController();
botController.initialise();