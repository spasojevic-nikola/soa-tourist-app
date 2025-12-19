import { Response, NextFunction } from 'express';
import { AuthRequest } from '../types/jwt';

export const adminMiddleware = (req: AuthRequest, res: Response, next: NextFunction) => {
  if (!req.user) {
    return res.status(401).json({ error: 'Authentication required' });
  }

  if (req.user.role !== 'administrator') {
    return res.status(403).json({ error: 'Admin access required' });
  }

  next();
};