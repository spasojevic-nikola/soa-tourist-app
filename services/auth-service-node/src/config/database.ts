import { DataSource } from 'typeorm';
import { User } from '../models/User';

// Odredite da li smo u Docker-u ili lokalno
const isDocker = process.env.DOCKER_ENV === 'true';

const dbHost = isDocker ? 'postgres-auth' : 'localhost';
const dbPort = parseInt(process.env.EXT_AUTH_DB_PORT || '5432');

export const AppDataSource = new DataSource({
  type: 'postgres',
  host: process.env.AUTH_DB_HOST || 'localhost',
  port: parseInt(process.env.EXT_AUTH_DB_PORT || '5436'),
  username: process.env.AUTH_DB_USER || 'postgres',
  password: process.env.AUTH_DB_PASSWORD || 'password',
  database: process.env.AUTH_DB_NAME || 'authdb',
  entities: [User],
  synchronize: false,
  logging: false,
});

// Funkcija za inicijalizaciju baze sa retry logikom
export const initializeDatabase = async (): Promise<void> => {
  const maxRetries = 5;
  const retryDelay = 5000;

  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      console.log(`Attempting to connect to database (attempt ${attempt}/${maxRetries})...`);
      console.log(`Connecting to: ${dbHost}:${dbPort}`);
      await AppDataSource.initialize();
      console.log('Successfully connected to the database!');
      return;
    } catch (error) {
      console.error(`Failed to connect to database. Attempt ${attempt}/${maxRetries}. Error:`, error);
      
      if (attempt === maxRetries) {
        throw new Error('Fatal: Could not connect to the database after multiple retries.');
      }
      
      console.log(`Retrying in ${retryDelay / 1000} seconds...`);
      await new Promise(resolve => setTimeout(resolve, retryDelay));
    }
  }
};