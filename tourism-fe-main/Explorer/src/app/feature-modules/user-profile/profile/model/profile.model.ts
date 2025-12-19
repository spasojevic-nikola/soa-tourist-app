export interface User {
  id: number;
  username: string;
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

// podaci koje saljemo prilikom izmene
export interface UpdateUserProfilePayload {
  first_name?: string;
  last_name?: string;
  profile_image?: string;
  biography?: string;
  motto?: string;
}