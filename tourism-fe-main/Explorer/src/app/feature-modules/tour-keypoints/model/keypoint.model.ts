export interface KeyPoint {
  id: number;
  tourId: number;
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  address?: string;
  image: string;
  order: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateKeyPointPayload {
  name: string;
  description: string;
  latitude: number;
  longitude: number;
  image: File | null;
  order: number;
    address?: string
}