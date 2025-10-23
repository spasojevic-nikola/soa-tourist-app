import 'dotenv/config';
import express from 'express';
//import cors from 'cors';
import { initializeDatabase } from './config/database';
import { AuthController } from './controllers/authController';
import { AdminController } from './controllers/adminController';
import { HealthController } from './controllers/healthController';
import { authMiddleware } from './middleware/authMiddleware';
import { adminMiddleware } from './middleware/adminMiddleware';

const app = express();
const PORT = process.env.PORT || 8084;

// Middleware
//app.use(cors({
 // origin: 'http://localhost:4200',
  //credentials: true,
  //methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  //allowedHeaders: ['X-Requested-With', 'Content-Type', 'Authorization']
//}));
app.use(express.json());

// Health check route
app.get('/health', HealthController.healthCheck);

// Auth routes
app.post('/api/v1/auth/register', AuthController.register);
app.post('/api/v1/auth/login', AuthController.login);
app.get('/api/v1/auth/user/:id', AuthController.getUserById);

// Admin routes (protected with auth and admin middleware)
app.get('/api/v1/auth/admin/users', authMiddleware, adminMiddleware, AdminController.getAllUsers);
app.post('/api/v1/auth/admin/users/:id/block', authMiddleware, adminMiddleware, AdminController.blockUser);
app.post('/api/v1/auth/admin/users/:id/unblock', authMiddleware, adminMiddleware, AdminController.unblockUser);

// Initialize database and start server
const startServer = async (): Promise<void> => {
  try {
    await initializeDatabase();
    console.log('Database initialized successfully');

    app.listen(PORT, () => {
      console.log(`Auth service running on port ${PORT}`);
    });
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
};

startServer();