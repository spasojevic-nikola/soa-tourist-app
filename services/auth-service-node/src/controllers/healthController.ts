import { Request, Response } from 'express';

export class HealthController {
  static healthCheck = (req: Request, res: Response): void => {
    res.json({ status: 'auth-service healthy' });
  };
}