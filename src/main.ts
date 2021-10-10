import express from 'express';
import dotenv from 'dotenv';
import admin from 'firebase-admin';
import BotController from './controller/bot_controller';
import FirebaseDatabaseConnector from './database/firebase/firebase_database_connector';

const app = express();
const port = process.env.PORT || 3000;

app.listen(port, () => {
    return console.log(`Server is listening on ${port}`);
});

dotenv.config();

admin.initializeApp({
    credential: admin.credential.cert({
        clientEmail: process.env.FIREBASE_CLIENT_EMAIL,
        privateKey: process.env.FIREBASE_PRIVATE_KEY.replace('\\n', '\n'),
        projectId: process.env.FIREBASE_PROJECT_ID
    }),
    databaseURL: process.env.FIREBASE_DATABASE_URL
});

const botController = new BotController(new FirebaseDatabaseConnector());
botController.initialise();