import { Request, Response } from 'express';
import bcrypt from 'bcryptjs';
import { AppDataSource } from '../config/database';
import { User } from '../models/User';
import { generateToken, JWTClaims } from '../types/jwt';
import { StakeholderUser } from '../types/stakeholder';
import axios from 'axios';

const userRepository = AppDataSource.getRepository(User);

export class AuthController {
  // Register user
  static register = async (req: Request, res: Response): Promise<void> => {
    try {
      const { username, email, password, role = 'tourist' } = req.body;

      // Validate input
      if (!username || !email || !password) {
        res.status(400).json({ error: 'Username, email and password are required' });
        return;
      }

      // Hash password
      const hashedPassword = await bcrypt.hash(password, 12);

      // Create user
      const user = new User();
      user.username = username;
      user.email = email;
      user.password = hashedPassword;
      user.role = role;
      user.blocked = false;
      user.created_at = new Date();
      user.updated_at = new Date();

      // Save to database
      await userRepository.save(user);

      // Send to stakeholders-service
      const stakeholderUser: StakeholderUser = {
        id: user.id,
        first_name: '',
        last_name: '',
        username: user.username,
        role: user.role,
        profile_image: '',
        biography: '',
        motto: ''
      };

      try {
        await axios.post(
          `${process.env.STAKEHOLDERS_SERVICE_URL || 'http://stakeholders-service:8080'}/api/v1/user`,
          stakeholderUser
        );
      } catch (error) {
        // Rollback user creation if stakeholders service fails
        await userRepository.remove(user);
        console.error('Stakeholders-service unreachable, rolling back user:', error);
        res.status(500).json({ 
          error: 'User creation failed in stakeholders-service. Rolled back auth-service entry.' 
        });
        return;
      }

      // Generate token
      const tokenClaims: JWTClaims = {
        id: user.id,
        username: user.username,
        role: user.role
      };

      const token = generateToken(tokenClaims);

      res.status(201).json({ accessToken: token });
    } catch (error: any) {
      console.error('Registration error:', error);
      
      if (error.code === '23505') { // PostgreSQL unique violation
        res.status(409).json({ error: 'User already exists' });
      } else {
        res.status(500).json({ error: 'Internal server error' });
      }
    }
  };

  // Login user
  static login = async (req: Request, res: Response): Promise<void> => {
    try {
      const { username, password } = req.body;

      if (!username || !password) {
        res.status(400).json({ error: 'Username and password are required' });
        return;
      }

      // Find user
      const user = await userRepository.findOne({ where: { username } });

      if (!user) {
        res.status(401).json({ error: 'Invalid credentials' });
        return;
      }

      // Check password
      const isPasswordValid = await bcrypt.compare(password, user.password);

      if (!isPasswordValid) {
        res.status(401).json({ error: 'Invalid credentials' });
        return;
      }

      // Check if user is blocked
      if (user.blocked) {
        res.status(403).json({ error: 'Account is blocked' });
        return;
      }

      // Generate token
      const tokenClaims: JWTClaims = {
        id: user.id,
        username: user.username,
        role: user.role
      };

      const token = generateToken(tokenClaims);

      console.log('🔐 GENERISAN TOKEN:', token);

      res.json({ accessToken: token });
    } catch (error) {
      console.error('Login error:', error);
      res.status(500).json({ error: 'Internal server error' });
    }
  };

  // Get user by ID
  static getUserById = async (req: Request, res: Response): Promise<void> => {
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

      res.json({
        username: user.username,
        email: user.email,
        role: user.role,
        blocked: user.blocked
      });
    } catch (error) {
      console.error('Get user error:', error);
      res.status(500).json({ error: 'Internal server error' });
    }
  };
}