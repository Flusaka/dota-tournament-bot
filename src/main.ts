import express from 'express';
import dotenv from 'dotenv';
import BotController from './controller/bot_controller';

const app = express();
const port = process.env.PORT || 3000;

app.listen(port, () => {
    return console.log(`Server is listening on ${port}`);
});

dotenv.config();

const botController = new BotController();
botController.initialise();