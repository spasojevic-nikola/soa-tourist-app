import { Request, Response } from 'express';
import { AppDataSource } from '../config/database';
import { User } from '../models/User';

const userRepository = AppDataSource.getRepository(User);

export class AdminController {
  // Get all users
  static getAllUsers = async (req: Request, res: Response): Promise<void> => {
    try {
      const users = await userRepository.find();
      
      // Remove passwords from response
      const usersWithoutPasswords = users.map(user => ({
        id: user.id,
        username: user.username,
        email: user.email,
        role: user.role,
        blocked: user.blocked,
        created_at: user.created_at,
        updated_at: user.updated_at
      }));

      res.json(usersWithoutPasswords);
    } catch (error) {
      console.error('Get all users error:', error);
      res.status(500).json({ error: 'Failed to fetch users' });
    }
  };

  // Block user
  static blockUser = async (req: Request, res: Response): Promise<void> => {
    try {
      const userId = parseInt(req.params.id);

      if (isNaN(userId)) {
        res.status(400).json({ error: 'Invalid user ID' });
        return;
      }

      const user = await userRepository.findOne({ where: { id: userId } });

      if (!user) {
        res.status(404).json({ error: 'User not found' });
        return;
      }

      if (user.role === 'administrator') {
        res.status(400).json({ error: 'Cannot block administrator' });
        return;
      }

      user.blocked = true;
      user.updated_at = new Date();
      
      await userRepository.save(user);

      // Return user without password
      const { password, ...userWithoutPassword } = user;
      res.json(userWithoutPassword);
    } catch (error) {
      console.error('Block user error:', error);
      res.status(500).json({ error: 'Failed to block user' });
    }
  };

  // Unblock user
  static unblockUser = async (req: Request, res: Response): Promise<void> => {
    try {
      const userId = parseInt(req.params.id);

      if (isNaN(userId)) {
        res.status(400).json({ error: 'Invalid user ID' });
        return;
      }

      const user = await userRepository.findOne({ where: { id: userId } });

      if (!user) {
        res.status(404).json({ error: 'User not found' });
        return;
      }

      user.blocked = false;
      user.updated_at = new Date();
      
      await userRepository.save(user);

      // Return user without password
      const { password, ...userWithoutPassword } = user;
      res.json(userWithoutPassword);
    } catch (error) {
      console.error('Unblock user error:', error);
      res.status(500).json({ error: 'Failed to unblock user' });
    }
  };
}