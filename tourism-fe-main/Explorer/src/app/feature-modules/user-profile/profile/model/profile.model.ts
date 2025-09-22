export interface User {
  id: number;
  first_name: string;
  last_name: string;
  profile_image: string;
  biography: string;
  motto: string;
  role: string;
  is_blocked: boolean;
  created_at: Date;
  updated_at: Date;
}