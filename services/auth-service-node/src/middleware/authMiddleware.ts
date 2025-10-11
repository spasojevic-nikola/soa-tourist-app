import { Response, NextFunction } from 'express';
import { verifyToken, AuthRequest, JWTClaims } from '../types/jwt';

export const authMiddleware = (req: AuthRequest, res: Response, next: NextFunction) => {
  // Koristite req.headers umesto req.header()
  const tokenHeader = req.headers.authorization;
  
  if (!tokenHeader) {
    return res.status(401).json({ error: 'Authorization header required' });
  }

  // Remove "Bearer " prefix if present
  const token = tokenHeader.startsWith('Bearer ') ? tokenHeader.slice(7) : tokenHeader;

  try {
    const decoded = verifyToken(token) as JWTClaims;
    req.user = decoded;
    next();
  } catch (error) {
    return res.status(401).json({ error: 'Invalid token' });
  }
};