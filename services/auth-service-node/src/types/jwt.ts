import jwt from 'jsonwebtoken';
import { Request } from 'express';

export interface JWTClaims {
  id: number;
  username: string;
  role: string;
}

// Extend Express Request interface
export interface AuthRequest extends Request {
  user?: JWTClaims;
}

export const JWT_SECRET = process.env.JWT_SECRET || 'super-tajni-kljuc-koji-niko-ne-zna-12345';
export const JWT_EXPIRES_IN = '1h';

export const generateToken = (claims: JWTClaims): string => {
  return jwt.sign(claims, JWT_SECRET, { expiresIn: JWT_EXPIRES_IN });
};

export const verifyToken = (token: string): JWTClaims => {
  return jwt.verify(token, JWT_SECRET) as JWTClaims;
};